package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AppTokenClaims struct {
	UserID          string   `json:"user_id"`
	Email           string   `json:"email"`
	FullName        string   `json:"full_name"`
	OrgID           string   `json:"org_id"`
	Role            string   `json:"role"`
	Apps            []string `json:"apps"`
	Plan            string   `json:"plan"`
	IsPlatformAdmin bool     `json:"is_platform_admin"`
}

// handleSSO exchanges a one-time platform app token for a PRODUCT_NAME session cookie.
// Platform redirects here with ?token=<t> after the user clicks this product in the launcher.
func handleSSO() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "token required"})
			return
		}

		url := fmt.Sprintf("%s/api/v1/internal/verify-app-token?token=%s", cfg.authServiceURL, token)
		req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, url, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		req.Header.Set("X-Internal-Key", cfg.internalKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		var claims AppTokenClaims
		if err := json.Unmarshal(body, &claims); err != nil || claims.UserID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "bad token claims"})
			return
		}

		now := time.Now()
		sessionJWT, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":               claims.UserID,
			"user_id":           claims.UserID,
			"email":             claims.Email,
			"full_name":         claims.FullName,
			"org_id":            claims.OrgID,
			"role":              claims.Role,
			"apps":              claims.Apps,
			"plan":              claims.Plan,
			"is_platform_admin": claims.IsPlatformAdmin,
			"iat":               now.Unix(),
			"exp":               now.Add(8 * time.Hour).Unix(),
		}).SignedString(sessionSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "session error"})
			return
		}

		secure := cfg.cookieDomain != "localhost"
		c.SetCookie(sessionCookie, sessionJWT, 8*3600, "/", cfg.cookieDomain, secure, true)

		redirect := c.Query("redirect")
		if redirect == "" {
			redirect = "/"
		}
		c.Redirect(http.StatusFound, redirect)
	}
}

// handleBootstrap returns session context for the frontend on load.
func handleBootstrap() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		mc, _ := claims.(jwt.MapClaims)

		var apps []string
		if raw, ok := mc["apps"]; ok {
			if list, ok := raw.([]any); ok {
				for _, v := range list {
					if s, ok := v.(string); ok {
						apps = append(apps, s)
					}
				}
			}
		}
		plan, _ := mc["plan"].(string)
		isAdmin, _ := mc["is_platform_admin"].(bool)

		c.JSON(http.StatusOK, gin.H{
			"user": gin.H{
				"id":                c.GetString("user_id"),
				"email":             c.GetString("email"),
				"display_name":      c.GetString("display_name"),
				"role":              c.GetString("role"),
				"is_platform_admin": isAdmin,
			},
			"org": gin.H{
				"id": c.GetString("org_id"),
			},
			"capabilities": gin.H{
				"apps": apps,
				"plan": plan,
			},
		})
	}
}

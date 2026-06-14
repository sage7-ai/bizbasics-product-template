package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const sessionCookie = "PRODUCT_NAME_session"

// sessionSecret signs THIS product's own session cookie. It is product-local
// and independent of the platform — never the platform JWT secret.
var sessionSecret = []byte(cfg.sessionSecret)

// CORSMiddleware allows requests from known bizbasics origins.
func CORSMiddleware() gin.HandlerFunc {
	allowed := map[string]bool{
		"https://app.bizbasics.ai": true,
	}
	for _, o := range cfg.corsOrigins {
		allowed[strings.TrimSpace(o)] = true
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if allowed[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// AuthMiddleware validates the product session cookie and sets gin context keys.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, err := c.Cookie(sessionCookie)
		if err != nil || raw == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing session"})
			return
		}

		tok, err := jwt.Parse(raw, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return sessionSecret, nil
		})
		if err != nil || !tok.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
			return
		}

		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "bad claims"})
			return
		}

		userID, _ := claims["user_id"].(string)
		orgID, _ := claims["org_id"].(string)
		email, _ := claims["email"].(string)
		role, _ := claims["role"].(string)
		name, _ := claims["full_name"].(string)
		if name == "" {
			name = email
		}

		c.Set("user_id", userID)
		c.Set("org_id", orgID)
		c.Set("role", role)
		c.Set("email", email)
		c.Set("display_name", name)
		c.Set("claims", claims)
		c.Next()
	}
}

// requireAppAccess rejects requests from orgs not entitled to this product.
func requireAppAccess(product string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		mc, _ := claims.(jwt.MapClaims)
		apps, _ := mc["apps"].([]any)
		for _, a := range apps {
			if s, ok := a.(string); ok && s == product {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "product not entitled"})
	}
}

package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var gitSHA = "dev"

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service", "PRODUCT_NAME"))

	db, err := sql.Open("postgres", cfg.dbURL)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("db ping: %v", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(10 * time.Minute)

	runMigrations(db)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(CORSMiddleware(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "PRODUCT_NAME", "version": gitSHA})
	})
	r.GET("/ready", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "db unavailable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	r.GET("/auth/sso", handleSSO())
	r.GET("/bootstrap", AuthMiddleware(), handleBootstrap())

	api := r.Group("/api/v1")
	api.Use(AuthMiddleware(), requireAppAccess("PRODUCT_NAME"))
	{
		// Register your routes here:
		// api.GET("/items", handleListItems(db))
		// api.POST("/items", handleCreateItem(db))
		// api.PATCH("/items/:id", handleUpdateItem(db))
		// api.DELETE("/items/:id", handleDeleteItem(db))
	}

	slog.Info("starting", "port", cfg.port, "version", gitSHA)
	if err := r.Run(":" + cfg.port); err != nil {
		log.Fatalf("server: %v", err)
	}
}

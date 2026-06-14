package main

import (
	"log"
	"os"
	"strings"
)

type config struct {
	port           string
	dbURL          string
	redisURL       string
	authServiceURL string
	platformAPIURL string
	internalKey    string
	cookieDomain   string
	corsOrigins    []string
	sessionSecret  string
}

var cfg = func() config {
	c := config{
		port:           getenv("PORT", "PRODUCT_PORT"),
		dbURL:          mustenv("DATABASE_URL"),
		redisURL:       mustenv("REDIS_URL"),
		authServiceURL: getenv("AUTH_SERVICE_URL", "http://auth-service:8091"),
		platformAPIURL: getenv("PLATFORM_API_URL", "http://platform-api:8092"),
		internalKey:    os.Getenv("PLATFORM_INTERNAL_KEY"),
		// Host-only cookie by default — your product's session belongs to YOUR
		// domain, not the shared parent. Do not set this to bizbasics.ai.
		cookieDomain: getenv("COOKIE_DOMAIN", ""),
		// Your product's OWN session-signing secret. Generate a random value
		// (>=32 bytes); it is NOT the platform JWT secret and must never be.
		sessionSecret: mustenv("APP_SESSION_SECRET"),
	}
	origins := os.Getenv("CORS_EXTRA_ORIGINS")
	if origins != "" {
		c.corsOrigins = strings.Split(origins, ",")
	}
	return c
}()

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required env var %s is not set", key)
	}
	return v
}

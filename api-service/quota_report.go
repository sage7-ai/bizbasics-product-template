package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

var quotaHTTPClient = &http.Client{Timeout: 5 * time.Second}

// reportQuota fires a best-effort quota consume call to platform-api.
// Call this after performing a quota-gated operation. Runs in a goroutine —
// never blocks the request path.
//
// resource: one of "ai_calls", "api_calls", "storage_bytes", "seats"
// amount:   number of units consumed (use storage bytes for file ops)
func reportQuota(orgID, resource string, amount int64) {
	if amount <= 0 || cfg.platformAPIURL == "" || cfg.internalKey == "" {
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		payload, _ := json.Marshal(map[string]any{
			"org_id":     orgID,
			"resource":   resource,
			"amount":     amount,
			"product_id": "PRODUCT_NAME", // replace with your product_id
		})
		req, err := http.NewRequestWithContext(ctx, http.MethodPost,
			fmt.Sprintf("%s/v1/quota/consume", cfg.platformAPIURL),
			bytes.NewReader(payload),
		)
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Internal-Key", cfg.internalKey)

		resp, err := quotaHTTPClient.Do(req)
		if err != nil {
			slog.Warn("quota report failed", "org_id", orgID, "resource", resource, "err", err)
			return
		}
		resp.Body.Close()
	}()
}

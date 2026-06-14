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

var wsPublishClient = &http.Client{Timeout: 10 * time.Second}

// publishWorkspaceRecord pushes a record to the platform workspace records API.
// Call this after creating/updating a significant business object so monk and
// other platform products can surface it as cross-product context.
//
// Runs in a goroutine — never blocks the request path. Failures are logged
// as warnings only; they never affect the caller's response.
//
// recordType: a short noun describing the object (e.g. "task", "invoice", "ticket")
// sourceRef:  the object's primary key in your product's DB
// title:      display title shown in search results
// body:       plain-text summary (max ~2000 chars; longer values are truncated)
func publishWorkspaceRecord(orgID, recordType, sourceRef, title, body string) {
	if cfg.platformAPIURL == "" || cfg.internalKey == "" {
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if len(body) > 2000 {
			body = body[:2000]
		}

		payload, _ := json.Marshal(map[string]any{
			"org_id":         orgID,
			"source_product": "PRODUCT_NAME", // replace with your product_id
			"record_type":    recordType,
			"source_ref":     sourceRef,
			"title":          title,
			"body":           body,
		})
		req, err := http.NewRequestWithContext(ctx, http.MethodPost,
			fmt.Sprintf("%s/v1/workspace-records", cfg.platformAPIURL),
			bytes.NewReader(payload),
		)
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Internal-Key", cfg.internalKey)

		resp, err := wsPublishClient.Do(req)
		if err != nil {
			slog.Warn("workspace record publish failed", "org_id", orgID, "record_type", recordType, "err", err)
			return
		}
		resp.Body.Close()
	}()
}

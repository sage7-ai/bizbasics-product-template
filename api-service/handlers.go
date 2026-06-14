package main

// Add your domain handlers here.
//
// Rule: every DB query MUST filter by org_id from the gin context.
//
// Example pattern:
//
//   func handleListItems(db *sql.DB) gin.HandlerFunc {
//       return func(c *gin.Context) {
//           orgID := c.GetString("org_id")
//           userID := c.GetString("user_id")
//
//           rows, err := db.QueryContext(c.Request.Context(), `
//               SELECT id, title, status, created_at
//               FROM "PRODUCT_NAME".example_items
//               WHERE org_id = $1
//               ORDER BY created_at DESC
//               LIMIT 100
//           `, orgID)
//           if err != nil {
//               c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
//               return
//           }
//           defer rows.Close()
//
//           var items []gin.H
//           for rows.Next() {
//               var id, title, status string
//               var createdAt time.Time
//               rows.Scan(&id, &title, &status, &createdAt)
//               items = append(items, gin.H{
//                   "id": id, "title": title,
//                   "status": status, "created_at": createdAt,
//               })
//           }
//           c.JSON(http.StatusOK, gin.H{"items": items, "org_id": orgID})
//
//           // Publish a workspace record when something significant happens:
//           // publisher.NewFromEnv().Publish(c.Request.Context(), publisher.Record{
//           //     OrgID: orgID, SourceProduct: "PRODUCT_NAME",
//           //     RecordType: "item", SourceRef: id, Title: title,
//           // })
//       }
//   }

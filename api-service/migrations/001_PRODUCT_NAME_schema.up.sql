-- Replace PRODUCT_NAME with your actual product slug (e.g. outreach, esign).
-- Every table MUST have org_id and filter by it in all queries.
CREATE SCHEMA IF NOT EXISTS "PRODUCT_NAME";

CREATE TABLE IF NOT EXISTS "PRODUCT_NAME".example_items (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id      UUID NOT NULL,
    created_by  UUID NOT NULL,
    title       TEXT NOT NULL,
    status      TEXT NOT NULL DEFAULT 'active',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_example_items_org_id
    ON "PRODUCT_NAME".example_items (org_id);

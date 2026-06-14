# bizbasics Product Scaffold

Starter template for building a product on the bizbasics platform.

## Quickstart

```bash
# Copy scaffold into your new product repo
cp -r . my-product
cd my-product

# Fill in your product details
sed -i 's/PRODUCT_NAME/my-product/g' **/*
sed -i 's/PRODUCT_PORT/8100/g' **/*

# Start dev environment
cp .env.example .env
docker compose up
```

## Structure

```
my-product/
├── api-service/         Go/Gin or Python/FastAPI backend
│   ├── main.go          — HTTP server, routes, SSO endpoint
│   ├── auth.go          — Platform token validation
│   ├── db.go            — Database setup (always filter by org_id)
│   ├── handlers/
│   └── Dockerfile
├── frontend/            React/Next.js frontend
│   ├── src/
│   └── Dockerfile
├── infra/k8s/           Kubernetes manifests (submit to platform team)
│   ├── api-deployment.yaml
│   ├── frontend-deployment.yaml
│   └── ingress.yaml
├── .github/workflows/
│   └── build.yml
└── docker-compose.yml
```

## Minimum Endpoints

See the [developer docs](https://bizbasics.io/docs).

Required:
- `POST /auth/sso`
- `GET  /bootstrap`
- `GET  /health`
- `GET  /ready`

## Tenant Isolation

**Every** database query must include `org_id` from the validated JWT.

## Design System

Install `@bizbasics/ui` for all UI components:

```bash
npm install @bizbasics/ui
```

Never build auth forms, nav, or loading states yourself — use the shared components.

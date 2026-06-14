FROM golang:1.22-alpine AS builder

WORKDIR /src
RUN apk add --no-cache git ca-certificates

COPY go.mod ./
RUN GONOSUMDB='*' GOFLAGS='-mod=mod' go mod download

COPY . .
ARG GIT_SHA=dev
RUN CGO_ENABLED=0 GOOS=linux GOFLAGS='-mod=mod' \
    go build -ldflags="-s -w -X main.gitSHA=${GIT_SHA}" -o /PRODUCT_NAME ./...

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /PRODUCT_NAME /app/PRODUCT_NAME
COPY --from=builder /src/migrations /app/migrations
EXPOSE PRODUCT_PORT
ENTRYPOINT ["/app/PRODUCT_NAME"]

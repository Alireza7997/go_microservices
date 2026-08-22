# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM golang:1.25-alpine AS build

WORKDIR /src

# Copy the workspace first so module resolution works
COPY go.work go.work.sum ./

# Copy the modules required to build the auth service
COPY pkg ./pkg
COPY config ./config
COPY general ./general
COPY auth_service ./auth_service

WORKDIR /src/auth_service

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/auth_service .

# ---- Runtime stage ----
FROM alpine:3.20

RUN adduser -D -u 10001 appuser

# Keep the same layout as local runs:
#   binary + migrations live in /app/auth_service, config is read from /app/env.yaml
WORKDIR /app/auth_service

COPY --from=build /out/auth_service /app/auth_service/auth_service
COPY --from=build /src/auth_service/migrations /app/auth_service/migrations

USER appuser

EXPOSE 6001

CMD ["./auth_service"]

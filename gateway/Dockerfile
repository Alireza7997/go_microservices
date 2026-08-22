# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM golang:1.25-alpine AS build

WORKDIR /src

# Copy the workspace first so module resolution works
COPY go.work go.work.sum ./

# Copy the modules required to build the gateway
COPY pkg ./pkg
COPY config ./config
COPY general ./general
COPY gateway ./gateway

WORKDIR /src/gateway

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/gateway .

# ---- Runtime stage ----
FROM alpine:3.20

RUN adduser -D -u 10001 appuser

WORKDIR /app/gateway

COPY --from=build /out/gateway /app/gateway/gateway

USER appuser

EXPOSE 6000

CMD ["./gateway"]

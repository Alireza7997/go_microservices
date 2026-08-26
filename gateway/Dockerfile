# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM golang:1.25-alpine AS build

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /src

COPY go.work.sum ./

# Copy the modules required to build the gateway
COPY pkg ./pkg
COPY config ./config
COPY general ./general
COPY auth_service ./auth_service
COPY chat_service ./chat_service
COPY greet_service ./greet_service
COPY gateway ./gateway

# Create a workspace scoped to only the modules present in this build context
RUN go work init ./gateway ./pkg ./config ./general ./auth_service ./chat_service ./greet_service

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

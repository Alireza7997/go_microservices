# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM golang:1.25-alpine AS build

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /src

COPY go.work.sum ./

COPY pkg ./pkg
COPY config ./config
COPY greet_service ./greet_service

# Create a workspace scoped to only the modules present in this build context
RUN go work init ./greet_service ./pkg ./config

WORKDIR /src/greet_service

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/greet_service .

# ---- Runtime stage ----
FROM alpine:3.20

RUN adduser -D -u 10001 appuser

WORKDIR /app/greet_service

COPY --from=build /out/greet_service /app/greet_service/greet_service

USER appuser

EXPOSE 6003

CMD ["./greet_service"]

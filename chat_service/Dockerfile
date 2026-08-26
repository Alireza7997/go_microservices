# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM golang:1.25-alpine AS build

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /src

COPY go.work.sum ./

COPY pkg ./pkg
COPY config ./config
COPY chat_service ./chat_service

# Create a workspace scoped to only the modules present in this build context
RUN go work init ./chat_service ./pkg ./config

WORKDIR /src/chat_service

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/chat_service .

# ---- Runtime stage ----
FROM alpine:3.20

RUN adduser -D -u 10001 appuser

WORKDIR /app/chat_service

COPY --from=build /out/chat_service /app/chat_service/chat_service

USER appuser

EXPOSE 6002

CMD ["./chat_service"]

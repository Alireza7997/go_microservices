.PHONY: help build test tidy lint proto run-auth run-gateway docker-build docker-up docker-down clean

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build all services
	cd gateway && go build -o bin/gateway .
	cd auth_service && go build -o bin/auth_service .

test: ## Run all tests
	cd pkg && go test ./...
	cd auth_service && go test ./...
	cd general && go test ./...
	cd config && go test ./...

tidy: ## Tidy all module dependencies
	for m in pkg config general auth_service gateway; do (cd $$m && go mod tidy); done

lint: ## Run go vet on all modules
	for m in pkg config general auth_service gateway; do (cd $$m && go vet ./...); done

proto: ## Regenerate protobuf code
	protoc --proto_path=. --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		general/general.proto auth_service/auth_pb/auth.proto

run-gateway: ## Run the gateway locally (expects ../env.yaml)
	cd gateway && go run .

run-auth: ## Run the auth service locally (expects ../env.yaml)
	cd auth_service && go run .

docker-build: ## Build all docker images
	docker compose build

docker-up: ## Start postgres + auth + gateway
	docker compose up --build -d

docker-down: ## Stop the stack
	docker compose down

clean: ## Remove built binaries and test caches
	rm -rf gateway/bin auth_service/bin
	go clean -cache -testcache

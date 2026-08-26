.PHONY: help build test tidy lint proto run-auth run-chat run-greet run-gateway docker-build docker-up docker-down clean

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build all services
	for m in gateway auth_service chat_service greet_service; do (cd $$m && go build -o bin/$(notdir $$m) .); done

test: ## Run all tests
	for m in pkg config general auth_service chat_service greet_service gateway; do (cd $$m && go test ./...); done

tidy: ## Tidy all module dependencies
	for m in pkg config general auth_service chat_service greet_service gateway; do (cd $$m && go mod tidy); done

lint: ## Run go vet on all modules
	for m in pkg config general auth_service chat_service greet_service gateway; do (cd $$m && go vet ./...); done

proto: ## Regenerate protobuf code
	PATH=$$PATH:$$(go env GOPATH)/bin protoc --proto_path=. --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		general/general.proto auth_service/auth_pb/auth.proto \
		chat_service/chat_pb/chat.proto greet_service/greet_pb/greet.proto

run-gateway: ## Run the gateway locally (expects ../env.yaml)
	cd gateway && go run .

run-auth: ## Run the auth service locally (expects ../env.yaml)
	cd auth_service && go run .

run-chat: ## Run the chat service locally (expects ../env.yaml)
	cd chat_service && go run .

run-greet: ## Run the greet service locally (expects ../env.yaml)
	cd greet_service && go run .

docker-build: ## Build all docker images
	docker compose build

docker-up: ## Start postgres + all services + gateway
	docker compose up --build -d

docker-down: ## Stop the stack
	docker compose down

clean: ## Remove built binaries and test caches
	rm -rf */bin
	go clean -cache -testcache

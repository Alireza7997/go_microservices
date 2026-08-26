# go_microservices

A Go microservice starter built as a [Go workspace](go.work) with an HTTP gateway that talks to backend services over gRPC.

## Architecture

```
                        HTTP (JSON)
    Client ────────────────────────────►  Gateway :6000
                                             │
                     gRPC                    │ gRPC                gRPC
             ┌───────────────────────────────┼──────────────────────────┐
             ▼                               ▼                          ▼
      Auth Service :6001              Chat Service :6002        Greet Service :6003
      (users, JWT, Postgres)          (in-memory room messages) (stateless ping)
```

- **Gateway** (`gateway/`) — single public HTTP entrypoint. Validates and parses requests, forwards calls to services over gRPC.
- **Auth Service** (`auth_service/`) — user registration, bcrypt hashing, JWT issuing. Runs SQL migrations on startup.
- **Chat Service** (`chat_service/`) — in-memory per-room message store. Demonstrates a stateful service without external dependencies.
- **Greet Service** (`greet_service/`) — stateless ping RPC. Demonstrates how trivially new services plug into the architecture.
- Shared libraries in `pkg/` (router, errors, database, config loader), shared protobuf messages in `general/`, shared config structs in `config/`.

## Tech stack

| Area       | Choice                                        |
|------------|-----------------------------------------------|
| Language   | Go ≥ 1.25                                     |
| RPC        | gRPC + Protocol Buffers                       |
| HTTP       | `net/http` with a custom router (`pkg/router`)|
| Database   | PostgreSQL (via `lib/pq` + `sql-migrate`)     |
| Auth       | bcrypt password hashing, HS256 JWT (`golang-jwt/jwt/v5`) |
| Deploy     | Docker / docker compose                       |

## Repository layout

```
├── gateway/          # HTTP gateway (public API)
├── auth_service/     # Auth microservice (gRPC) + migrations
├── chat_service/     # Chat microservice (gRPC, in-memory)
├── greet_service/    # Greet microservice (gRPC, stateless)
├── pkg/              # Shared libraries: router, errors, database, yaml loader
├── config/           # Shared configuration structs
├── general/          # Shared protobuf messages
├── env.yaml.example  # Config template — copy to env.yaml
└── docker-compose.yml
```

## Getting started

### Prerequisites

- Go ≥ 1.25
- Docker & docker compose *(optional)*
- PostgreSQL *(only if running without docker)*

### 1. Configuration

```sh
cp env.yaml.example env.yaml   # then fill in real values
```

`env.yaml` is gitignored — never commit it.

### 2a. Run with docker compose

```sh
POSTGRES_PASSWORD=secret docker compose up --build -d
```

> For the full stack in docker, set `microservices.auth.ip` to `0.0.0.0`
> in `env.yaml`, and point the gateway at the service by hostname if you
> use container DNS names.

### 2b. Run locally

Start PostgreSQL, create a database named `authDB` (matching `env.yaml`),
then:

```sh
make run-auth      # terminal 1 — gRPC :6001
make run-chat      # terminal 2 — gRPC :6002
make run-greet     # terminal 3 — gRPC :6003
make run-gateway   # terminal 4 — HTTP :6000
```

All services read `../env.yaml` relative to their directory.

## API endpoints

| Method | Path                        | Body / Query                                                                        | Success |
|--------|-----------------------------|-------------------------------------------------------------------------------------|---------|
| POST   | `/api/auth/register/`       | `{"username": "...", "password": "...", "password_confirm": "...", "email": "..."}` | 201     |
| POST   | `/api/chat/messages/`       | `{"room": "...", "username": "...", "body": "..."}`                                 | 201     |
| GET    | `/api/chat/messages/{room}/`| —                                                                                   | 200     |
| GET    | `/api/greet/ping/?name=...` | — (`name` optional)                                                                 | 200     |

Examples:

```sh
curl -s -X POST localhost:6000/api/greet/ping/ # see below
```

```sh
# greet
curl 'localhost:6000/api/greet/ping/?name=dev'
# {"greeting":"hello, dev!","server_time":1755000000}

# chat: post and list messages
curl -X POST localhost:6000/api/chat/messages/ \
  -d '{"room":"general","username":"dev","body":"hi"}'
curl localhost:6000/api/chat/messages/general/
```

Errors are returned as JSON:

```json
{"code": 400, "message": "username already taken"}
```

## Development

```sh
make build        # compile all services
make test         # run unit tests
make lint         # go vet across all modules
make tidy         # tidy all go.mod files
make proto        # regenerate protobuf/gRPC code (needs protoc + plugins)
make help         # list all targets
```

### Adding a new service

1. `mkdir my_service/{my_pb,global,load,service}` and add `module github.com/Alireza7997/go_microservices/my_service` in its `go.mod`.
2. Define the `.proto`, run `make proto`.
3. Implement the service, listen on an address from `env.yaml` (see `greet_service/main.go` — the smallest example).
4. Register it under `microservices:` in `env.yaml`, add a call helper + handler in the gateway, and a route in `gateway/routes/http.go`.

Tests in `auth_service/service` use [`go-sqlmock`](https://github.com/DATA-DOG/go-sqlmock),
so they run without a database. See `CONTRIBUTING.md` for conventions.

## Security notes

- Secrets live only in `env.yaml` (gitignored). Use `env.yaml.example` as template.
- JWTs are signed HS256 with `secret_key` and expire after 24h.
- Passwords are hashed with bcrypt; they are never stored or logged in plaintext.

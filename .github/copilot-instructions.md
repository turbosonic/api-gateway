# GitHub Copilot Instructions for turbosonic/api-gateway

## Project Overview

This is a lightweight, sub-millisecond API Gateway for microservices, written in **Go**. It acts as the reverse-proxy conduit between the outside world and an internal microservice ecosystem. Routing is defined in a YAML configuration file. The gateway supports authentication, CORS, request ID tracking, role/scope-based access control, and pluggable logging backends. The final Docker image is built on `scratch` and is approximately 7 MB.

## Tech Stack

- **Language**: Go (module path `github.com/turbosonic/api-gateway`)
- **HTTP router**: [Goji](https://goji.io) (`goji.io`)
- **TLS**: `github.com/zenazn/goji/graceful`
- **Configuration**: YAML via `gopkg.in/yaml.v2`
- **Env vars**: `github.com/joho/godotenv`
- **Authentication**: Auth0 JWT validation (`github.com/auth0-community/go-auth0`)
- **Logging backends**: stdout, Elasticsearch, InfluxDB, Azure Application Insights
- **UUID**: `github.com/satori/go.uuid` (request IDs)
- **CI**: Travis CI (`.travis.yml`)
- **Containerisation**: Multi-stage Docker build → final `scratch` image

## Repository Layout

```
api-gateway/
├── app.go                          # main() — bootstraps mux, middleware, and server
├── go.mod
├── Makefile                        # make shell | run | test | build
├── Dockerfile                      # multi-stage: golang:1.13 build → scratch final
├── data/config.yaml                # example YAML routing config (mounted at runtime)
├── .env.example                    # environment variable template
├── authentication/
│   ├── authentication.go           # middleware: selects auth provider from env
│   └── clients/auth0/auth0.go      # Auth0 JWT validation, extracts roles/scopes
├── configurations/
│   └── configurations.go           # loads & parses data/config.yaml
├── factories/
│   └── logging.go                  # factory: returns the configured LogClient
├── initializer/
│   └── initializer.go              # registers all routes on the Goji mux
├── logging/
│   ├── logging.go                  # LogClient interface + middleware + log structs
│   └── clients/
│       ├── stdout/                 # prints to stdout
│       ├── elasticsearch/          # ships logs to Elasticsearch
│       ├── influxdb/               # batch-writes to InfluxDB
│       └── applicationinsights/    # sends telemetry to Azure App Insights
├── parammap/
│   └── parammap.go                 # maps :param URL segments → values
├── relay/
│   └── relay.go                    # HTTP client that forwards requests to backends
├── responseMarshal/
│   └── responseMarshal.go          # CORS middleware + request-ID header injection
└── tests/
    └── factories_test.go           # unit test: verifies LogClient factory returns non-nil
```

## Key Concepts and Patterns

### Middleware Stack (app.go)
Middleware is applied in this order on every request:
1. `responseMarshal.CorsHandler` — CORS headers, short-circuits OPTIONS
2. `responseMarshal.AddHeaders` — injects `X-Request-Id` (UUID v4)
3. `logging.LogHandlerFunc` — records external request + response asynchronously
4. `authentication.Handler` — validates token, attaches `roles` and `scopes` to context

### Authentication (authentication/authentication.go)
Controlled by the `AUTHENTICATION_PROVIDER` environment variable:
- `auth0` — validates JWT against Auth0 JWKS (`AUTH0_DOMAIN`, `AUTH0_AUDIENCE`)
- `none` — allows all requests (development / open endpoints)
- anything else (or unset) — **denies all requests** (fail-safe default)

Roles and scopes from a valid token are stored in the request context and checked per-endpoint in `initializer`.

### Routing & Relay (initializer/initializer.go + relay/relay.go)
`RegisterEndpoints` iterates the parsed configuration and registers one Goji pattern per `(url, method)` pair. For each request:
1. Role and scope checks are performed against the token claims in context.
2. URL path parameters (`:id` style) are substituted using `parammap`.
3. The query string is forwarded as-is.
4. The request body is read and sent to the backend via `relay.RelayRequest`.
5. Response headers and body are copied back to the client.

The relay HTTP client reuses connections (`MaxIdleConns: 10`, `IdleConnTimeout: 30s`) and filters out hop-by-hop / CORS headers before forwarding.

### Logging (logging/)
`logging.LogClient` is an interface with two methods:
```go
LogRequest(*RequestLog, string, string)  // external request log
LogRelay(*RelayLog, string, string)      // backend relay log
```
The factory in `factories/logging.go` reads `LOGGING_PROVIDER` and returns the appropriate concrete client. All logging calls are made asynchronously (`go logClient.Log...`).

### Configuration Format (data/config.yaml)
```yaml
name: "/web/v1"          # base path prefix
endpoints:
  - url: /resource/:id   # supports :param segments
    methods:
      - method: GET
        roles: ["admin"]   # required roles (use "*" for wildcard)
        scopes: ["read"]   # required scopes (use "*" for wildcard)
        destination:
          name: my-service
          host: http://my-service  # internal hostname (Docker network)
          url: /resource/:id       # backend URL (params are substituted)
```

## Environment Variables

| Variable | Purpose | Default |
|---|---|---|
| `PORT` | Listening port | `8080` |
| `CORS_ORIGIN` | Allowed CORS origin | *(none)* |
| `CORS_MAX_AGE` | CORS preflight cache seconds | `600` |
| `AUTHENTICATION_PROVIDER` | `auth0` / `none` | deny-all |
| `AUTH0_DOMAIN` | Auth0 tenant domain | — |
| `AUTH0_AUDIENCE` | Auth0 API audience | — |
| `LOGGING_PROVIDER` | `elasticsearch` / `influxdb` / `applicationinsights` / `stdout` | `stdout` |
| `LOGGING_ELASTIC_URL` | Elasticsearch base URL | — |
| `LOGGING_EXTERNAL_REQUEST_INDEX_NAME` | Elasticsearch index for gateway logs | — |
| `LOGGING_INTERNAL_REQUEST_INDEX_NAME` | Elasticsearch index for relay logs | — |
| `LOGGING_INFLUXDB_URL` | InfluxDB HTTP URL | — |
| `LOGGING_INFLUX_DB_NAME` | InfluxDB database name | — |
| `APPLICATIONINSIGHTS_INTRUMENTATION_KEY` | Azure App Insights key (note: "INTRUMENTATION" is a typo in the source — use this exact spelling) | — |

## Development Workflow

```bash
# Run tests
make test          # go test -v ./tests

# Build Docker image
make build         # docker build -t turbosonic/api-gateway .

# Run locally in Docker (port 8080)
make run

# Drop into an interactive Go shell
make shell
```

TLS is automatically enabled when `./data/certs/cert.pem` and `./data/certs/key.pem` are present.

## Coding Conventions

- Package names match their directory name (e.g. `package authentication`, `package relay`).
- Interfaces are defined in `logging/logging.go`; concrete implementations live under `logging/clients/<name>/`.
- The factory pattern (`factories/`) decouples provider selection from implementation.
- Error handling: non-recoverable startup errors use `panic()`; runtime errors use `log.Println()`.
- All log calls are asynchronous (`go logClient.LogRequest(...)`) to avoid blocking the request path.
- URL path parameters use the `:name` convention (not `{name}`).
- YAML config structs use exported fields; environment variable reads use `os.Getenv`.

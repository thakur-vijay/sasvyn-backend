# Sasvyn Backend

<p align="center">
  <img alt="Go version" src="https://img.shields.io/badge/Go-1.27-00ADD8?logo=go" />
  <img alt="License" src="https://img.shields.io/badge/License-MIT-green.svg" />
  <img alt="Docker" src="https://img.shields.io/badge/Docker-Postgres-2496ED?logo=docker" />
</p>

Sasvyn Backend is the Go API for the Sasvyn professional identity and portfolio application. The current codebase is a modular HTTP service built with Go's standard library, PostgreSQL, and migration-driven schema updates.

## Overview

The repository currently implements a small but functional backend foundation for:

- Apple ID based social login bootstrap
- session creation and token refresh
- authenticated access checks
- user profile retrieval and updates
- skills and language CRUD for authenticated users
- upload URL generation for object storage-backed profile assets
- PostgreSQL-backed persistence with automatic migrations
- request rate limiting for auth and social-login flows

This is not a full production authentication provider yet. The current social-login flow accepts an Apple ID from the client and creates or looks up the user locally; it does not yet verify Apple identity tokens on the server side.

## Stack

- Go 1.27
- PostgreSQL 17 via Docker Compose for local development
- pgx for PostgreSQL access
- golang-migrate for schema migrations
- standard library `net/http` router and middleware
- AWS S3 SDK for Cloudflare R2 presigned upload URLs
- `swaggo/swag` generated OpenAPI artifacts under `docs/`

## Architecture

The app follows a small modular monolith pattern rooted in `internal/` packages:

- `cmd/server` starts the application
- `internal/app` wires the HTTP router, DB, and rate limiters
- `internal/config` loads runtime configuration from environment variables
- `internal/database` connects to PostgreSQL and runs migrations
- `internal/modules/*` contains domain-oriented handlers, services, repositories, and routes
- `internal/ratelimit` applies request throttling per route policy
- `internal/storage` initializes the Cloudflare R2/S3 client used for upload URLs
- `internal/response` provides consistent JSON responses
- `internal/validation` validates request DTOs

Request flow is intentionally simple:

```text
HTTP request
  ↓
Router
  ↓
Auth middleware / route handler
  ↓
Repository / service logic
  ↓
PostgreSQL
```

## Project structure

```text
.
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── router.go
│   │   └── router_test.go
│   ├── app/
│   │   ├── app.go
│   │   ├── dependencies.go
│   │   ├── ratelimit.go
│   │   └── server.go
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   ├── database.go
│   │   ├── migrate.go
│   │   └── migrations/
│   ├── modules/
│   │   ├── auth/
│   │   ├── languages/
│   │   ├── sessions/
│   │   ├── skills/
│   │   ├── upload/
│   │   └── users/
│   ├── ratelimit/
│   ├── response/
│   ├── storage/
│   └── validation/
├── docs/
├── docker-compose.yml
├── Dockerfile
├── dev.sh
├── restart.sh
├── go.mod
├── go.sum
├── LICENSE
├── README.md
└── .gitignore
```

## Authentication

The repository currently implements a session-based auth flow using bearer tokens:

1. Client sends `POST /api/v1/auth/socialLogin` with an Apple user identifier and user details.
2. The server looks up the user by `apple_id`. If not found, it creates a new user row.
3. It generates a new access token and refresh token pair, stores only hashed token values in PostgreSQL, and returns both to the client.
4. Protected routes use the `Authorization: Bearer <token>` header.
5. The auth middleware validates the access token and stores `user_id` and `session_id` in request context.
6. `POST /api/v1/auth/refresh` accepts a refresh token and rotates access/refresh credentials.
7. `POST /api/v1/auth/logout` revokes the active session.

Important note: the implementation currently treats the request payload as user-provided identity data from the client and does not verify Apple JWTs or authorization codes server-side. The code is therefore best described as a working session foundation rather than production-ready Apple Sign In verification.

## API overview

All routes are mounted under `/api/v1` except the health check at `/health`.

| Method | Route | Auth | Purpose |
| --- | --- | --- | --- |
| GET | `/health` | No | Basic service health check |
| POST | `/api/v1/auth/socialLogin` | No | Create or look up a user and return tokens |
| POST | `/api/v1/auth/refresh` | No | Rotate session tokens |
| POST | `/api/v1/auth/logout` | Yes | Revoke the current session |
| GET | `/api/v1/users/{id}` | Yes | Return a user record |
| PUT | `/api/v1/users/{id}` | Yes | Update profile details such as full name, date of birth, and image key |
| GET | `/api/v1/skills` | Yes | List skills for the authenticated user |
| POST | `/api/v1/skills` | Yes | Create a skill |
| PUT | `/api/v1/skills/{id}` | Yes | Update a skill |
| DELETE | `/api/v1/skills/{id}` | Yes | Delete a skill |
| GET | `/api/v1/languages` | Yes | List languages for the authenticated user |
| POST | `/api/v1/languages` | Yes | Create a language entry |
| PUT | `/api/v1/languages/{id}` | Yes | Update a language entry |
| DELETE | `/api/v1/languages/{id}` | Yes | Delete a language entry |
| POST | `/api/v1/upload-url` | Yes | Generate a presigned upload URL for an image |

The API response format is a JSON envelope with `status_code`, `message`, and `data` fields. The code also applies per-route rate limiting for social login and token refresh requests.

## Database

The application uses PostgreSQL as the system of record and applies migrations automatically at startup.

### Current migrations

```text
000001_create_users.up.sql
000002_create_sessions.up.sql
000003_add_date_of_birth_to_users.up.sql
000004_create_skills.up.sql
000005_create_languages.up.sql
000006_add_img_url_to_users.up.sql
```

### Key tables

- `users`
  - `id` UUID primary key
  - `apple_id` unique text identifier
  - `full_name`, `email`, `date_of_birth`, `img_key`
  - `created_at`, `updated_at`
- `sessions`
  - `id` UUID primary key
  - `user_id` foreign key to `users.id`
  - hashed `access_token_hash` and `refresh_token_hash`
  - `expires_at`, `refresh_expires_at`, `revoked_at`
  - indexed by `user_id` and expiry timestamp
- `skills`
  - `id`, `user_id`, `skill`, `category`, `created_at`, `updated_at`
  - unique index on `(user_id, LOWER(skill))`
- `languages`
  - `id`, `user_id`, `language_code`, `language`, `proficiency`, timestamps
  - unique index on `(user_id, LOWER(language))`

`internal/database/migrate.go` runs `migrate.Up()` automatically when the app starts. This is the source of truth for schema evolution.

## Configuration and environment variables

The app currently reads configuration from environment variables with defaults only for the port. The practical configuration is:

| Variable | Required | Purpose |
| --- | --- | --- |
| `APP_PORT` | No | HTTP port for the server; defaults to `8080` |
| `DATABASE_URL` | Yes | PostgreSQL connection string |
| `R2_ACCOUNT_ID` | Yes for uploads | Cloudflare R2 account identifier |
| `R2_ACCESS_KEY_ID` | Yes for uploads | R2 access key ID |
| `R2_SECRET_ACCESS_KEY` | Yes for uploads | R2 secret access key |
| `R2_BUCKET_NAME` | Yes for uploads | Bucket name used for presigned upload URLs |
| `R2_PUBLIC_URL` | Yes for user profile image URLs | Public base URL used when returning image URL metadata |

The project does not currently include a checked-in `.env.example` file. Local development typically sets these variables in a local `.env` file, which is excluded by `.gitignore`.

## Local development

### Prerequisites

- Go 1.27+
- Docker and Docker Compose

### Start PostgreSQL

```bash
docker compose up -d postgres
```

This starts PostgreSQL on `localhost:5433` with the following local development database values:

- Host: `localhost`
- Port: `5433`
- Database: `sasvyn`
- User: `sasvyn`
- Password: `sasvyn`

### Start the server

```bash
export DATABASE_URL='postgresql://sasvyn:sasvyn@localhost:5433/sasvyn?sslmode=disable'
export APP_PORT='8080'

go run ./cmd/server
```

The server listens on `http://localhost:8080` and runs migrations automatically on startup.

### Helper scripts

```bash
./dev.sh
./restart.sh
```

`dev.sh` runs the server directly, while `restart.sh` tears down the local database volume and re-creates the container before launching the app.

## Testing

The repository includes Go tests for the router and ratelimit behavior.

Run the suite with:

```bash
go test ./...
```

Run a build and static validation with:

```bash
go vet ./...
go build ./...
```

## Docker and container notes

The repository includes:

- `docker-compose.yml` for a local PostgreSQL instance
- `Dockerfile` for building the Go binary into a container image

The local development setup currently uses Docker for PostgreSQL only; the application itself is typically run with `go run ./cmd/server` rather than via the Docker image during local development.

## Swagger / API docs

Generated Swagger artifacts are present in `docs/`:

- `docs/swagger.json`
- `docs/swagger.yaml`
- `docs/docs.go`

The server does not currently register a Swagger UI route, and the development script does not automatically regenerate the docs. The generated files are present for reference and can be used to inspect the API contract from the repository itself.

## Security and operational notes

- Tokens are stored as hashes, not plain-text secrets.
- Access tokens are validated and expired by timestamp.
- Refresh tokens are rotated when refreshed.
- Session revocation is supported with a `revoked_at` field.
- The app uses rate limiting on auth and social-login endpoints.
- Sensitive environment variables are expected to be loaded from local environment state or a local `.env` file, never committed to source control.

The current implementation does not perform Apple JWT validation, key verification, or authorization code exchange against Apple's servers. It should be treated as a local session model with a client-supplied Apple identifier, not as a fully verified Apple Sign In implementation.

## License

This project is licensed under the MIT License. See `LICENSE` for details.

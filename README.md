# Sasvyn Backend

Sasvyn Backend is the Go HTTP API for the Sasvyn professional identity and portfolio application. It is a modular monolith backed by PostgreSQL, with domain-oriented packages, versioned REST routes, and database migrations applied at server startup.

> **Implementation status:** This README describes the code currently in this repository. In particular, the Apple sign-in endpoint is a local account/session bootstrap, not verified Sign in with Apple authentication. See [Current limitations](#current-limitations) before exposing the API to production traffic.

## Contents

- [Implemented capabilities](#implemented-capabilities)
- [Technology](#technology)
- [Architecture](#architecture)
- [Getting started](#getting-started)
- [Configuration](#configuration)
- [HTTP API](#http-api)
- [Request and response conventions](#request-and-response-conventions)
- [Database and migrations](#database-and-migrations)
- [Rate limiting and idempotency](#rate-limiting-and-idempotency)
- [Tests and quality checks](#tests-and-quality-checks)
- [Docker](#docker)
- [API specification files](#api-specification-files)
- [Current limitations](#current-limitations)
- [License](#license)

## Implemented capabilities

- Apple-ID-keyed local user lookup/creation and database-backed access/refresh sessions.
- Authenticated user lookup and profile updates.
- User-owned skills, languages, and social-link endpoints.
- Authenticated Cloudflare R2-compatible presigned image-upload URLs.
- PostgreSQL persistence using SQL through Go's `database/sql` package and the pgx stdlib driver.
- Twelve versioned SQL migrations, automatically applied at server startup.
- In-memory, per-client-IP HTTP rate limiting.
- Database-backed idempotency for selected write operations.
- A lightweight `GET /health` endpoint.

This is not a complete portfolio API: projects, experience, education, documents, and other future domains are not implemented here.

## Technology

| Area | Implementation |
| --- | --- |
| Language | Go 1.27 (module and Docker build image) |
| HTTP | Go standard library `net/http`, including method-aware `ServeMux` patterns |
| Database | PostgreSQL 17 for the supplied local Compose service |
| SQL driver | pgx v5 through `database/sql` |
| Migrations | golang-migrate; SQL files under `internal/database/migrations/` |
| Object storage | S3-compatible AWS SDK v2 client configured for Cloudflare R2 |
| Request validation | `go-playground/validator` for DTOs that call the validation helper |
| API prefix | `/api/v1`; health is served separately at `/health` |

## Architecture

The application is composed in `cmd/server` and wired in `internal/app`. The `internal/modules/` packages own feature routes, handlers, persistence models, and repositories. Auth, sessions, and idempotency also have service logic; not every module uses a separate service layer.

```text
cmd/server
    └── config → app composition → HTTP server
                              ├── API router and middleware
                              ├── feature handlers/repositories
                              ├── PostgreSQL and migrations
                              └── R2/S3 client for upload URLs
```

### Repository map

```text
cmd/
├── server/                 # API server entry point
└── cleanup/                # Removes expired idempotency records
internal/
├── api/                    # Health endpoint and route composition
├── app/                    # Dependencies, rate limits, server lifecycle
├── config/                 # APP_PORT and DATABASE_URL
├── database/               # PostgreSQL connection and migration runner
│   └── migrations/         # Up/down SQL migration pairs
├── domain/model/           # Shared base fields for selected models
├── modules/
│   ├── auth/               # Social-login bootstrap, refresh, logout, auth middleware
│   ├── idempotency/        # Idempotency-key persistence and execution
│   ├── languages/
│   ├── sessions/
│   ├── skills/
│   ├── socialLinks/
│   ├── upload/
│   └── users/
├── ratelimit/              # In-memory limiter and HTTP policy middleware
├── response/               # JSON decoding and response envelopes
├── storage/                # R2 client and object URL helpers
└── validation/             # DTO validation helpers and messages
docs/                       # Generated Swagger artifacts (not currently served)
docker-compose.yml          # Local PostgreSQL service only
Dockerfile                  # Multi-stage API and cleanup binary image
```

## Getting started

### Prerequisites

- Go 1.27 or a compatible later Go toolchain.
- Docker with the Compose plugin (for the local PostgreSQL service).

### 1. Start PostgreSQL

From the repository root:

```bash
docker compose up -d postgres
```

The Compose service publishes PostgreSQL on `localhost:5433` and creates a **local-development-only** database with:

| Setting | Value |
| --- | --- |
| Database | `sasvyn` |
| User | `sasvyn` |
| Password | `sasvyn` |
| Host port | `5433` |

Do not reuse these local credentials in a shared or production environment.

### 2. Configure and run the API

The server loads `.env` from the current working directory when present; otherwise it uses the process environment. `.env` is ignored by Git. There is currently no checked-in `.env.example`.

```bash
export DATABASE_URL='postgres://sasvyn:sasvyn@localhost:5433/sasvyn?sslmode=disable'
export APP_PORT='8080'

go run ./cmd/server
```

On startup, the server connects to PostgreSQL, applies pending migrations, and listens on `http://localhost:8080`. Run it from the repository root so the migration source path resolves.

Verify the HTTP liveness endpoint:

```bash
curl -i http://localhost:8080/health
```

Expected response:

```json
{"status":"ok"}
```

`/health` is a simple liveness response; it does not check PostgreSQL or R2 availability.

## Configuration

| Variable | Required | Default / purpose |
| --- | --- | --- |
| `APP_PORT` | No | HTTP listen address; defaults to `8080`. A leading `:` is accepted. |
| `DATABASE_URL` | Yes | PostgreSQL connection URL. Required by both server and cleanup command. |
| `R2_ACCOUNT_ID` | For R2 operations | Cloudflare account identifier used to construct the R2 endpoint. |
| `R2_ACCESS_KEY_ID` | For R2 operations | R2 access key ID. |
| `R2_SECRET_ACCESS_KEY` | For R2 operations | R2 secret access key. |
| `R2_BUCKET_NAME` | For upload URL generation | Bucket used when presigning PUT requests. |
| `R2_PUBLIC_URL` | For profile image URLs | Base URL prepended to a saved user image key when returning a user. |

`APP_PORT` and `DATABASE_URL` are loaded in `internal/config`; R2 values are read directly from the environment by storage and upload code. The application does not validate R2 configuration at startup, so configure it before relying on upload or profile-image behavior. Never commit real credentials.

## HTTP API

All feature routes are mounted beneath `/api/v1`. Routes marked **Yes** require `Authorization: Bearer <access_token>`. The exact camel case in `/auth/socialLogin` and `/socialLinks` is intentional.

### Health and authentication

| Method | Path | Auth | Purpose |
| --- | --- | :---: | --- |
| `GET` | `/health` | No | Basic liveness response (`{"status":"ok"}`); not a dependency/readiness check. |
| `POST` | `/api/v1/auth/socialLogin` | No | Find or create a local user by the submitted `apple_id`, then create a session. |
| `POST` | `/api/v1/auth/refresh` | No | Validate and rotate the submitted refresh token. |
| `POST` | `/api/v1/auth/logout` | Yes | Revoke the current authenticated session. |

Social-login request:

```json
{
  "apple_id": "apple-subject-from-client",
  "full_name": "Example User",
  "email": "user@example.com"
}
```

The server currently accepts these identity fields as client-provided data; it does **not** validate an Apple identity token or exchange an authorization code. The first local user lookup uses `apple_id`; later logins create a new session for that user.

Refresh request:

```json
{
  "refresh_token": "<refresh-token>"
}
```

Successful login returns `data.user`, `data.access_token`, and `data.refresh_token`. A newly created access token is valid for 15 minutes and its refresh token for 30 days. Refreshing rotates both token values.

### Users

| Method | Path | Auth | Idempotency key | Purpose |
| --- | --- | :---: | :---: | --- |
| `GET` | `/api/v1/users/{id}` | Yes | — | Retrieve the user at the supplied ID. |
| `PUT` | `/api/v1/users/{id}` | Yes | Required | Update supplied profile fields. |

`PUT /users/{id}` accepts any non-empty subset of:

```json
{
  "full_name": "Example User",
  "date_of_birth": "1990-01-31",
  "img_key": "users/<user-id>/<object-id>"
}
```

`img_key` is the storage object key, not the uploaded file itself. The returned user representation may include `img_url` built from `R2_PUBLIC_URL` and that key. `ImageKey` is not included in user JSON responses.

### Skills

| Method | Path | Auth | Idempotency key | Purpose |
| --- | --- | :---: | :---: | --- |
| `GET` | `/api/v1/skills` | Yes | — | List the authenticated user's skills. |
| `POST` | `/api/v1/skills` | Yes | — | Create a skill; the server generates its ID. |
| `PUT` | `/api/v1/skills/{id}` | Yes | — | Update a skill's name and category. |
| `DELETE` | `/api/v1/skills/{id}` | Yes | — | Delete one of the authenticated user's skills. |

Create/update fields are `skill` and `category`:

```json
{
  "skill": "Go",
  "category": "Backend"
}
```

### Languages

| Method | Path | Auth | Idempotency key | Purpose |
| --- | --- | :---: | :---: | --- |
| `GET` | `/api/v1/languages` | Yes | — | List the authenticated user's languages. |
| `GET` | `/api/v1/languages/{id}` | Yes | — | Retrieve one of the authenticated user's languages. |
| `POST` | `/api/v1/languages` | Yes | Required | Create a language entry. |
| `PUT` | `/api/v1/languages/{id}` | Yes | Required | Partially update a language entry. |
| `DELETE` | `/api/v1/languages/{id}` | Yes | Required | Delete a language entry. |

Create requires a client-generated UUID v4 `id`, `language_code`, `language`, and `proficiency` from 1 to 5:

```json
{
  "id": "d9428888-122b-4c7e-9a6a-45f69c9a7b0d",
  "language_code": "en",
  "language": "English",
  "proficiency": 5
}
```

Update accepts any non-empty subset of `language_code`, `language`, and `proficiency`.

### Social links

| Method | Path | Auth | Idempotency key | Purpose |
| --- | --- | :---: | :---: | --- |
| `GET` | `/api/v1/socialLinks` | Yes | — | List the authenticated user's social links. |
| `GET` | `/api/v1/socialLinks/{id}` | Yes | — | Retrieve one of the authenticated user's social links. |
| `POST` | `/api/v1/socialLinks` | Yes | Required | Create a social link. |
| `PUT` | `/api/v1/socialLinks/{id}` | Yes | Required | Partially update a social link. |
| `DELETE` | `/api/v1/socialLinks/{id}` | Yes | Required | Delete a social link. |

Create requires a client-generated UUID v4 `id`, a `type`, and a valid `url`. Supported types are `github`, `linkedin`, `x`, `instagram`, `youtube`, `dribbble`, `behance`, `medium`, and `website`. Update accepts `type` and/or `url`.

### Upload URLs

| Method | Path | Auth | Idempotency key | Purpose |
| --- | --- | :---: | :---: | --- |
| `POST` | `/api/v1/upload-url` | Yes | — | Generate a presigned R2-compatible object upload URL. |

Request:

```json
{
  "type": "profile_image",
  "content_type": "image/png"
}
```

Accepted `type` values: `profile_image`, `project_screenshot`, `app_icon`. Accepted `content_type` values: `image/jpeg`, `image/png`, `image/webp`. The response contains `data.upload_url` and `data.img_key`; the presigned PUT URL expires after 10 minutes. This endpoint creates an upload URL only—it does not upload the file or persist a project record.

## Request and response conventions

- Feature endpoints use JSON request bodies. Authenticated requests send `Authorization: Bearer <access_token>`.
- A successful item/list response uses an envelope with `status_code`, `message`, and `data`. List responses return an empty array when there are no records.
- Error responses generally use the same JSON envelope with a message and no `data`; status codes vary by handler (commonly `400`, `401`, `404`, `409`, and `500`).
- Rate-limit rejections are emitted by `http.Error` as plain text, not the normal JSON envelope.
- Mutations marked above require a non-empty `Idempotency-Key` header. Reusing a completed key for the same authenticated user replays its stored status and response. A key already being processed returns `409 Conflict`.
- List endpoints currently do not implement pagination or client-selectable limits.
- Authenticated child-resource queries use the authenticated user ID from request context. See the ownership limitation below for user-by-ID endpoints.

Example authenticated request:

```bash
curl -i \
  -H 'Authorization: Bearer <access-token>' \
  -H 'Content-Type: application/json' \
  -H 'Idempotency-Key: <unique-request-key>' \
  -d '{"id":"d9428888-122b-4c7e-9a6a-45f69c9a7b0d","language_code":"en","language":"English","proficiency":5}' \
  http://localhost:8080/api/v1/languages
```

## Database and migrations

PostgreSQL is the persistent system of record. Migration files are stored in `internal/database/migrations/` and `internal/database.Migrate` applies pending migrations automatically during server startup.

The current migration sequence is:

```text
000001_create_users
000002_create_sessions
000003_add_date_of_birth_to_users
000004_create_skills
000005_create_languages
000006_add_img_url_to_users
000007_add_sync_version_to_users
000008_create_social_links
000009_add_sync_version_to_social_links
000010_create_idempotency_keys
000011_add_status_to_idempotency_keys
000012_add_sync_version_to_languages
```

The schema currently contains:

| Table | Purpose and notable constraints |
| --- | --- |
| `users` | UUID primary key; unique Apple ID; name, email, optional date of birth and image key; sync version and timestamps. |
| `sessions` | User-owned sessions with unique access/refresh token hashes, expirations, and optional revocation time. |
| `skills` | User-owned skill and category; case-insensitive skill uniqueness per user. |
| `languages` | User-owned language code, name, proficiency, and sync version; case-insensitive name uniqueness per user. |
| `social_links` | User-owned link type and URL with sync version; case-insensitive URL uniqueness per user. |
| `idempotency_keys` | Stored response/status associated with a user and idempotency key; the pair is unique. |

Skills, languages, social links, sessions, and idempotency records reference users with cascading deletion. Skills currently do not have a `sync_version`; user, language, and social-link records do.

### Expired idempotency-record cleanup

`cmd/cleanup` removes idempotency records older than 24 hours. It is a one-shot command, not a scheduled service:

```bash
go run ./cmd/cleanup
```

Set `DATABASE_URL` (or provide it through `.env`) before running it. No scheduler is configured in this repository.

## Rate limiting and idempotency

The in-memory rate limiter keys requests by client IP and applies these per-minute policies:

| Requests | Limit |
| --- | ---: |
| Default policy (all other requests) | 20 |
| `POST /api/v1/auth/refresh` | 10 |
| `POST /api/v1/auth/socialLogin` | 5 |

Responses include `X-RateLimit-Limit` and `X-RateLimit-Remaining`. Rejected requests return `429 Too Many Requests` with `Retry-After`. Because these counters are in memory, they are not shared across server instances and reset when the process restarts.

Idempotency currently applies to user updates, language create/update/delete, and social-link create/update/delete. The key is stored per user in PostgreSQL and completed response bodies are replayed. The cleanup command removes records after 24 hours; cleanup is not automatically scheduled.

## Tests and quality checks

Run the Go test suite, static analysis, and build from the repository root:

```bash
go test ./...
go vet ./...
go build ./...
```

The current tests cover router composition and rate-limiter behavior. There is no database integration-test setup in the repository.

## Docker

`docker-compose.yml` starts PostgreSQL only. The `Dockerfile` builds both the server and cleanup binaries, copies migration files into the runtime image, and starts the API server. The Compose file does not currently define or launch the API container.

## API specification files

Files under `docs/` are generated Swagger artifacts, but the application does not register a Swagger UI or serve an OpenAPI route. The checked-in generated specification currently describes a `/vijay` test endpoint rather than the API routes documented above; treat it as stale and not as the API contract.

## Current limitations

- **Apple identity is not verified.** `/auth/socialLogin` trusts the submitted `apple_id`, `full_name`, and `email`. It does not verify Apple JWT signatures/claims or exchange authorization codes. Do not treat this endpoint as production-ready Sign in with Apple.
- **User ID ownership is not enforced on user-by-ID routes.** `GET` and `PUT /users/{id}` require authentication, but the current handlers use the path ID without comparing it to the authenticated user's ID. Other user-owned collection/item queries generally scope database access by the authenticated user ID. Do not assume the user endpoints enforce ownership.
- **No pagination.** Collection routes return all matching records.
- **Uploads are presigning only.** The client uploads directly to storage using the returned URL; this API does not verify that upload completed or persist a project/document entity.
- **Rate limiting is process-local.** Limits are not coordinated across multiple instances.
- **Health is liveness only.** `/health` does not check database or object-storage dependencies.
- **Configuration is partly distributed.** R2 environment variables are read directly by the storage/upload code rather than loaded into the central config object.
- **No automated idempotency cleanup.** Expired keys are removed only when `cmd/cleanup` is run.
- **No implemented project or broader portfolio modules.** Upload type names do not imply project CRUD endpoints exist.

## License

This project is licensed under the MIT License. See [LICENSE](./LICENSE).

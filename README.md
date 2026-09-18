# Sasvyn Backend

Sasvyn Backend is the Go HTTP API for the Sasvyn professional identity and portfolio application. It currently provides a small authentication foundation backed by PostgreSQL:

- Social login user creation and lookup
- Short-lived access tokens and refresh tokens
- Authenticated access to the current user
- Automatic database migrations on startup

The service uses Go's standard `net/http` package, PostgreSQL through `pgx`, and `golang-migrate` for schema migrations.

## Current Status

This repository is under active development. The current social login endpoint accepts an `apple_id` in the request body; server-side verification of Apple identity tokens has not been implemented yet. Do not treat the current endpoint as production-ready Sign in with Apple authentication.

There is currently no health-check endpoint. The server exposes the routes listed in the API section below.

## Prerequisites

- Go 1.27 or newer
- Docker and Docker Compose, for the local PostgreSQL instance

## Quick Start

1. Start PostgreSQL:

	 ```sh
	 docker compose up -d postgres
	 ```

	 The container creates the following local database:

	 | Setting | Value |
	 | --- | --- |
	 | Host | `localhost` |
	 | Port | `5433` |
	 | Database | `sasvyn` |
	 | User | `sasvyn` |
	 | Password | `sasvyn` |

2. Start the API:

	 ```sh
	 go run ./cmd/server
	 ```

	 On startup, the application connects to PostgreSQL, applies pending migrations, and listens on `http://localhost:8080`.

3. To stop the local database:

	 ```sh
	 docker compose down
	 ```

	 To remove the database volume as well:

	 ```sh
	 docker compose down -v
	 ```

The database connection string and HTTP port are currently defined in the source code. `DATABASE_URL` and `PORT` are not supported yet.

## Development Commands

Run the server through the development script:

```sh
./dev.sh
```

Reset the local PostgreSQL container and start the server:

```sh
./restart.sh
```

Run the test suite:

```sh
go test ./...
```

Format the Go source:

```sh
go fmt ./...
```

Run static checks and build the application:

```sh
go vet ./...
go build ./...
```

## API

### `POST /socialLogin`

Creates a user when the supplied Apple ID is new, or creates a new session for an existing user.

Request:

```json
{
	"apple_id": "apple-user-id",
	"full_name": "Jane Appleseed",
	"email": "jane@example.com"
}
```

The response contains the user, an access token, and a refresh token.

### `POST /auth/refresh`

Rotates the tokens for a valid refresh token.

Request:

```json
{
	"refresh_token": "refresh-token"
}
```

### `GET /me`

Returns the authenticated user. Send the access token using the Bearer scheme:

```sh
curl http://localhost:8080/me \
	-H 'Authorization: Bearer ACCESS_TOKEN'
```

Access tokens expire after 15 minutes. Refresh tokens expire after 30 days.

## Database Migrations

Migration files are stored in [`internal/database/migrations`](internal/database/migrations). The application applies pending migrations automatically when it starts.

The current schema contains:

- `users`, identified by a unique `apple_id`
- `sessions`, containing hashed access and refresh tokens linked to users

Schema changes must be made through versioned migration files rather than manual database edits.

## Project Structure

```text
cmd/server/              Application entry point
internal/auth/           Authentication middleware and token refresh handler
internal/database/       PostgreSQL connection and migrations
internal/sessions/       Session persistence and token lifecycle
internal/users/          User persistence and user endpoints
docs/                    Generated Swagger artifacts
docker-compose.yml       Local PostgreSQL service
```

## Documentation

Generated Swagger artifacts are available in [`docs`](docs). The current server does not register a Swagger UI route. `dev.sh` starts the server; Swagger regeneration is not currently enabled in that script.

## License

See [`LICENSE`](LICENSE).
Sasvyn Backend — AGENTS.md

1. Project Identity

Project name: Sasvyn Backend

Repository: sasvyn-backend

Purpose:

Sasvyn Backend is the backend/API service for the Sasvyn professional identity and portfolio application.

The backend will provide:

* Authentication
* User identity
* Profile synchronization
* Portfolio data synchronization
* Projects
* Skills
* Experience
* Education
* Languages
* Social links
* Documents
* Mockups
* GitHub integration
* Future cloud synchronization
* Future AI-related backend capabilities where appropriate

The backend is written in Go and uses PostgreSQL as its primary relational database.

⸻

2. Current Project State

The repository starts intentionally minimal.

Do NOT assume that existing architecture, folders, packages, models, APIs, or database tables exist.

Build the system incrementally.

Do not generate the entire backend at once.

Each feature/module must be implemented only when explicitly requested.

Example:

“Build the Apple authentication module.”

The agent should implement only the necessary foundation and Apple authentication functionality required for that module.

Do not proactively implement unrelated modules.

⸻

3. Core Technology Stack

Use the following stack unless explicitly instructed otherwise.

Language

Go.

Use the current stable Go version available at development time.

Database

PostgreSQL.

PostgreSQL is the primary persistent database.

HTTP

Prefer Go’s standard library:

net/http

Do not introduce Gin, Echo, Fiber, or another web framework unless explicitly requested.

PostgreSQL driver

Use:

pgx

Prefer the official/current pgx ecosystem APIs appropriate for the selected version.

Database migrations

Use:

golang-migrate

Database schema changes must be represented as version-controlled migrations.

API

REST/HTTP APIs.

Use JSON request and response bodies.

Configuration

Environment variables.

Never hard-code:

* passwords
* API keys
* private keys
* JWT secrets
* database credentials
* Apple credentials
* OAuth secrets

⸻

4. Architecture Philosophy

Sasvyn Backend should initially be a:

Modular Monolith

Do NOT create microservices.

The backend should be structured by business/domain modules rather than by generic technical layers spread across the entire application.

Preferred conceptual structure:

internal/
├── auth/
├── user/
├── profile/
├── project/
├── skill/
├── experience/
├── education/
├── language/
├── sociallink/
├── document/
├── mockup/
└── github/

A module may contain files such as:

module/
├── handler.go
├── service.go
├── repository.go
├── model.go
└── ...

Only create files that are actually necessary.

Do not create empty or speculative abstractions.

⸻

5. Dependency Philosophy

Keep the dependency graph small.

Before introducing a third-party library, ask:

1. Does Go’s standard library already solve this?
2. Does the project genuinely need the dependency?
3. Does the dependency simplify the architecture?
4. Does it introduce unnecessary maintenance or security risk?

Do not add dependencies simply because they are popular.

Avoid dependency-driven architecture.

⸻

6. ORM Policy

Do NOT introduce GORM or another ORM unless explicitly requested.

Prefer:

SQL
↓
pgx
↓
PostgreSQL

The project should maintain clear control over SQL queries and database behavior.

If query volume/repetition becomes significant, sqlc may be considered.

Do not introduce sqlc automatically without evaluating whether it is useful for the current module.

⸻

7. Database Rules

PostgreSQL is the source of truth for server-side persistent data.

Every schema change must use migrations.

Example:

migrations/
├── 000001_create_users.up.sql
├── 000001_create_users.down.sql
├── 000002_create_profiles.up.sql
└── 000002_create_profiles.down.sql

Never make undocumented schema changes manually in the development database.

Never rely on TablePlus changes as the source of truth.

TablePlus is only a database client.

The migration files are the source of truth for schema evolution.

⸻

8. Database Design Principles

Use:

* UUIDs where appropriate
* foreign keys
* unique constraints
* indexes where justified
* NOT NULL where appropriate
* database-level constraints for important invariants
* timestamps for lifecycle tracking

Avoid:

* unnecessary nullable fields
* duplicated data
* arbitrary JSON blobs when relational modeling is more appropriate
* storing derived data unnecessarily
* database designs optimized only for the current UI

Database design must represent domain concepts, not SwiftUI screens.

⸻

9. Naming Conventions

Go:

Use idiomatic Go naming.

Examples:

User
UserRepository
CreateUser
GetUserByID

Avoid:

UserModel
UserServiceManager
UserHelper
UserUtil

unless the suffix genuinely communicates responsibility.

Database:

Use:

snake_case

Examples:

user_id
created_at
updated_at
apple_user_id

Go files:

lowercase.go

Package names:

auth
user
project
github

Avoid package names such as:

helpers
utils
common
misc
stuff

unless there is a strong architectural reason.

⸻

10. API Design

All APIs must have predictable HTTP semantics.

Use appropriate HTTP methods:

GET
POST
PUT
PATCH
DELETE

Use appropriate HTTP status codes.

Examples:

200 OK
201 Created
204 No Content
400 Bad Request
401 Unauthorized
403 Forbidden
404 Not Found
409 Conflict
422 Unprocessable Entity
500 Internal Server Error

Do not return 200 OK for every situation.

⸻

11. API Response Design

Responses should be consistent.

Success responses should contain only information useful to the client.

Errors should be structured JSON.

Example:

{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "The request is invalid."
  }
}

Do not expose:

* PostgreSQL errors
* SQL queries
* stack traces
* internal filesystem paths
* secrets
* private tokens
* implementation details

to clients.

Internal logs may contain more debugging information, but must still avoid secrets and sensitive credentials.

⸻

12. Error Handling

Follow idiomatic Go error handling.

Prefer:

if err != nil {
    return err
}

Do not silently ignore errors.

Do not use panic for ordinary application errors.

Use panic only for genuinely unrecoverable initialization/programming conditions where appropriate.

Errors should preserve useful context.

Prefer:

fmt.Errorf("create user: %w", err)

rather than:

return errors.New("error")

when wrapping an underlying error.

⸻

13. Context

Use context.Context for operations that may involve:

* HTTP requests
* database queries
* external APIs
* network operations

Pass request context through the relevant layers.

Do not create unnecessary background contexts inside business logic.

Avoid:

context.Background()

when an existing request context should be propagated.

⸻

14. Authentication

Sasvyn initially supports:

Sign in with Apple only.

Do not implement:

* Google authentication
* GitHub authentication
* email/password authentication
* Facebook authentication
* arbitrary OAuth providers

unless explicitly requested later.

Apple authentication must be designed around server-side verification of Apple’s identity/token information.

Never blindly trust identity information supplied by the client.

The backend must verify the appropriate Apple credential/token before treating the user as authenticated.

⸻

15. Apple Identity

Apple’s stable user identifier should be treated as an identity identifier.

Do not use the user’s display name as identity.

Do not use email alone as the primary identity key.

The backend must account for Apple’s private relay email addresses.

Apple may provide the user’s name only under specific circumstances, so the implementation must not assume the name will always be available.

The first successful authentication should persist appropriate user information.

Subsequent logins should correctly identify the existing user.

⸻

16. Authentication Security

Never log:

* Apple identity tokens
* refresh tokens
* access tokens
* private keys
* client secrets
* authorization codes
* database passwords

Do not store sensitive credentials in plain text unless there is a specific and justified requirement.

Authentication secrets must come from environment variables or an appropriate secret-management mechanism.

⸻

17. Authorization

Authentication and authorization are separate concerns.

Authentication answers:

Who is this user?

Authorization answers:

Is this user allowed to access this resource?

Every protected resource must verify ownership/authorization where applicable.

Example:

A user must not be able to request another user’s project merely by changing:

/project/{id}

in the request.

⸻

18. User Data Ownership

Sasvyn is a personal professional identity application.

User-owned resources should be associated with a user identity.

Conceptually:

users
  |
  ├── profile
  ├── projects
  ├── skills
  ├── experiences
  ├── education
  ├── languages
  ├── social_links
  ├── documents
  └── mockups

Queries involving user-owned resources must enforce ownership.

Do not trust a client-provided user_id when the authenticated user identity is already available.

Prefer deriving the user identity from the authenticated session/token.

⸻

19. Local-First Architecture

Sasvyn’s iOS/macOS application is designed around local-first data storage.

The backend should therefore NOT assume that every user interaction requires a network request.

The client may operate locally and synchronize with the backend.

Backend architecture should remain compatible with future synchronization.

Do not build APIs that unnecessarily force a network-dependent application architecture.

⸻

20. Sync Philosophy

Future synchronization may involve:

Local database
      ↕
Sasvyn API
      ↕
PostgreSQL

The backend should therefore avoid designs that make synchronization unnecessarily difficult.

Where appropriate, domain records should have:

id
created_at
updated_at

and potentially additional synchronization metadata when the actual sync design requires it.

Do NOT add complicated sync/versioning systems before the synchronization requirements are defined.

Do not invent conflict-resolution logic prematurely.

⸻

21. File Storage

Documents, mockups, images, and other potentially large binary assets should not automatically be stored directly inside PostgreSQL.

When cloud file storage is introduced, use object storage and store metadata/reference information in PostgreSQL.

Possible future architecture:

iOS
 ↓
Sasvyn API
 ↓
Object Storage

and:

PostgreSQL
 ↓
metadata
file IDs
URLs/references
ownership
timestamps

Do not introduce an object-storage provider until the actual requirement is defined.

⸻

22. GitHub Integration

GitHub integration is separate from authentication.

Do not use GitHub login as the primary Sasvyn authentication mechanism.

GitHub integration will eventually allow users to:

* connect GitHub
* disconnect GitHub
* persist the connection
* fetch repositories
* select repositories/projects for Sasvyn

OAuth tokens/secrets must never be exposed to the iOS application unnecessarily.

The backend should own sensitive GitHub integration credentials where appropriate.

⸻

23. Logging

Use structured, useful logging.

Logs should help diagnose:

* request failures
* database failures
* external API failures
* authentication failures
* application startup problems

Do not log secrets.

Do not log complete authentication tokens.

Do not log sensitive user information unnecessarily.

Avoid excessive debug logging in production.

⸻

24. Configuration

Configuration should be centralized.

Example:

APP_ENV
APP_PORT
DATABASE_URL
APPLE_CLIENT_ID
APPLE_TEAM_ID
APPLE_KEY_ID
APPLE_PRIVATE_KEY

Do not scatter environment-variable reads throughout the application.

Prefer loading configuration once and passing configuration explicitly where required.

⸻

25. Environment Files

Local development may use:

.env

The repository should contain:

.env.example

but never commit actual secrets.

.gitignore must include:

.env

⸻

26. Docker

Docker should be used where it improves reproducibility.

PostgreSQL may run through Docker during development.

The application itself does not need to be containerized immediately unless required.

Do not add Kubernetes or other orchestration infrastructure.

⸻

27. Testing

Important business logic must have tests.

Use Go’s built-in testing framework:

testing

Prefer unit tests for:

* services
* validation
* authentication logic
* authorization
* business rules
* transformations

Database-dependent behavior should have appropriate integration tests where useful.

Do not write meaningless tests that only increase coverage percentages.

Tests must verify actual behavior.

⸻

28. Test Naming

Prefer descriptive test names.

Example:

func TestCreateUser_WithExistingAppleID_ReturnsExistingUser(t *testing.T)

rather than:

func TestUser1(t *testing.T)

Table-driven tests should be used when testing multiple related cases.

⸻

29. Input Validation

Never assume client input is valid.

Validate:

* required fields
* formats
* lengths
* IDs
* enums
* pagination parameters
* URLs where applicable
* data relationships

Validation should happen before business logic and database operations where appropriate.

⸻

30. SQL Safety

Never construct SQL using untrusted string concatenation.

Bad:

query := "SELECT * FROM users WHERE email = '" + email + "'"

Use parameterized queries.

Example:

query := `SELECT id, email FROM users WHERE email = $1`

⸻

31. Transactions

Use database transactions when multiple database operations must succeed or fail together.

Example:

Create user
+
Create profile
+
Create initial data

If these operations represent one atomic business operation, use a transaction.

Do not use transactions unnecessarily for independent operations.

⸻

32. Concurrency

Go is concurrent by design.

Do not introduce goroutines without understanding ownership and lifecycle.

Every goroutine should have a clear reason to exist.

Avoid goroutine leaks.

Use context cancellation where appropriate.

Do not create unbounded goroutines from HTTP requests.

⸻

33. HTTP Server

The HTTP server must have:

* sensible timeouts
* graceful shutdown
* proper context handling
* clear route registration
* centralized middleware where appropriate

Do not use an unconfigured HTTP server in production.

Avoid:

http.ListenAndServe(...)

as the final production implementation without server configuration.

⸻

34. Health Check

The backend should provide a simple health endpoint.

Example:

GET /health

The health endpoint should provide enough information for local development and deployment health checks.

Do not expose internal system information.

⸻

35. API Versioning

Do not create unnecessary API versions.

If versioning is introduced, prefer a clear structure such as:

/api/v1/...

Do not create:

/v1
/v2
/v3

without an actual compatibility requirement.

⸻

36. Pagination

List endpoints should eventually support pagination where datasets can grow.

Do not return an unlimited number of records from endpoints such as:

GET /projects
GET /github/repositories

For small fixed datasets, pagination may not be necessary.

Use sensible limits.

Never allow a client to request an arbitrarily huge result set.

⸻

37. API Documentation

Document public APIs clearly.

When a module is implemented, document:

* endpoint
* HTTP method
* request body
* authentication requirement
* response
* errors
* important behavior

Do not create massive documentation for nonexistent APIs.

Documentation should reflect the actual implementation.

⸻

38. Git Rules

Commit meaningful changes.

Prefer commits such as:

feat: add apple authentication
feat: add project repository
fix: handle duplicate apple identity
test: add user service tests
refactor: simplify auth middleware
chore: add database migrations

Avoid:

changes
update
stuff
final
final2
working

Do not commit:

* .env
* secrets
* private keys
* generated credentials
* build artifacts

⸻

39. Code Formatting

Always format Go code with:

gofmt

or:

go fmt ./...

Code must be formatted before completion.

Use idiomatic Go.

Avoid unnecessary comments.

Comments should explain why, not obvious implementation details.

⸻

40. Linting

Use golangci-lint when configured.

Before considering a module complete, run:

go test ./...
go vet ./...
golangci-lint run

Fix meaningful issues rather than disabling linters blindly.

Never add a lint exclusion merely to make CI pass without understanding why the warning exists.

⸻

41. Build Verification

Before declaring a module complete:

go test ./...
go vet ./...
go build ./...

If linting is configured:

golangci-lint run

The agent must not claim a feature is complete if the project does not compile.

⸻

42. Security Principles

Treat all client input as untrusted.

Protect:

* authentication
* authorization
* database access
* external API credentials
* file access
* user-owned resources

Follow least privilege.

Never expose internal credentials.

Never trust IDs from the client without verifying ownership.

Do not implement custom cryptography.

Use established, maintained cryptographic libraries and protocols.

⸻

43. Performance Philosophy

Do not optimize prematurely.

First prioritize:

1. Correctness
2. Security
3. Maintainability
4. Testability
5. Observability
6. Performance

Optimize after identifying an actual bottleneck.

Do not introduce:

* Redis
* Kafka
* RabbitMQ
* Elasticsearch
* Kubernetes
* caching layers
* background worker systems

unless an actual requirement justifies them.

⸻

44. Scalability Philosophy

The initial architecture should be capable of scaling without prematurely implementing distributed systems.

Start with:

One Go application
        ↓
PostgreSQL

Scale vertically first.

When necessary, the Go API should be designed to allow multiple application instances behind a load balancer.

Do not introduce distributed architecture before it is necessary.

⸻

45. Generic Utility Rule

Avoid creating generic utility packages.

Bad:

internal/utils/
internal/helpers/
internal/common/

unless the functionality genuinely belongs there.

Prefer domain ownership.

For example:

internal/auth/

should own authentication-specific logic.

⸻

46. Dependency Injection

Use simple dependency injection.

Prefer explicit constructors:

func NewService(repo Repository) *Service

Avoid large dependency-injection frameworks.

The application’s composition/root wiring should happen near the application entry point.

⸻

47. Application Startup

The application should have a clear startup flow:

Load configuration
        ↓
Initialize logger
        ↓
Connect to PostgreSQL
        ↓
Initialize repositories
        ↓
Initialize services
        ↓
Initialize HTTP handlers/routes
        ↓
Start HTTP server

Shutdown should reverse the relevant resources gracefully.

⸻

48. Domain Layer Rules

Business rules should not live inside HTTP handlers.

Bad:

HTTP Handler
 ├── validation
 ├── business logic
 ├── SQL
 └── response formatting

Prefer:

HTTP Handler
     ↓
Service / Use Case
     ↓
Repository
     ↓
PostgreSQL

The exact implementation may vary by module, but responsibilities must remain clear.

⸻

49. Repository Rules

Repositories are responsible for persistence concerns.

They should not:

* construct HTTP responses
* know about HTTP status codes
* contain UI logic
* depend on Swift/iOS concepts

Repositories deal with domain data and database operations.

⸻

50. Service Rules

Services/use cases contain business logic.

They should not directly depend on HTTP-specific concerns unless there is a compelling reason.

They should operate on meaningful domain operations.

Example:

CreateProject
UpdateProject
DeleteProject
GetProject
ListProjects

rather than generic:

DoSomething
ProcessData
HandleRequest

⸻

51. Handler Rules

HTTP handlers should be thin.

Responsibilities:

1. Read request
2. Validate transport-level input
3. Authenticate/obtain user context
4. Call service/use case
5. Translate result into HTTP response

Do not put complex business logic into handlers.

⸻

52. No Fake Functionality

Never implement fake functionality just to make an endpoint appear complete.

Do not return:

{
  "success": true
}

when the underlying operation has not actually happened.

Do not create mock repositories for production code unless explicitly requested.

If functionality is incomplete, clearly state what remains.

⸻

53. No Speculative Features

Do not implement future features without being asked.

Examples:

Do NOT automatically implement:

* Google login
* GitHub login
* subscriptions
* payments
* notifications
* AI agents
* Redis caching
* analytics
* admin dashboards
* recommendation systems
* microservices

Sasvyn will be built incrementally.

⸻

54. AI Coding Agent Rules

The coding agent must:

1. Inspect the existing repository before modifying anything.
2. Understand existing architecture before introducing new patterns.
3. Reuse existing abstractions where appropriate.
4. Avoid unnecessary dependencies.
5. Avoid unrelated changes.
6. Never overwrite working code without a reason.
7. Never delete functionality without explicit permission.
8. Keep changes focused on the requested module.
9. Run tests/build/lint after meaningful changes.
10. Report assumptions clearly.

If requirements are ambiguous and the ambiguity affects architecture, ask before making a major architectural decision.

Do not ask unnecessary questions when the intended implementation is obvious.

⸻

55. Module Development Workflow

Every new module should follow this general workflow:

1. Understand requirements
        ↓
2. Identify domain model
        ↓
3. Design database schema
        ↓
4. Create migration
        ↓
5. Implement repository
        ↓
6. Implement service/use case
        ↓
7. Implement HTTP handler
        ↓
8. Register routes
        ↓
9. Add validation
        ↓
10. Add tests
        ↓
11. Run formatting
        ↓
12. Run tests
        ↓
13. Run vet/lint
        ↓
14. Verify API behavior

Not every module needs every layer.

Do not create abstractions solely because the workflow says they exist.

⸻

56. First Development Milestone

The first milestone should establish the backend foundation.

Expected initial structure:

sasvyn-backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── database/
│   └── http/
├── migrations/
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── docker-compose.yml
└── README.md

Initial functionality:

Go application starts
        ↓
PostgreSQL connection
        ↓
HTTP server starts
        ↓
GET /health
        ↓
successful response

Do not implement authentication or business modules as part of this milestone unless explicitly requested.

⸻

57. Future Module Order

The exact order may change, but the initial logical direction is:

Foundation
    ↓
Apple Authentication
    ↓
User
    ↓
Profile
    ↓
Projects
    ↓
Skills
    ↓
Experience
    ↓
Education
    ↓
Languages
    ↓
Social Links
    ↓
Documents
    ↓
Mockups
    ↓
GitHub Integration
    ↓
Synchronization

Implement only the module explicitly requested.

⸻

58. Important Product Architecture Constraint

Sasvyn is not simply a CRUD application.

The backend exists to support:

Sasvyn iOS/macOS
        ↓
Local-first experience
        ↓
Cloud synchronization
        ↓
Sasvyn Backend
        ↓
PostgreSQL

Therefore, backend decisions must consider:

* offline usage
* synchronization
* user ownership
* data consistency
* conflict handling
* eventual multi-device usage

However, do not implement synchronization complexity until the synchronization specification is explicitly defined.

⸻

59. Definition of Done

A module is considered complete only when:

* The implementation matches the requested requirements.
* The code compiles.
* Database migrations work.
* Database constraints are correct.
* Authentication/authorization is correctly enforced where required.
* Errors are handled.
* Important business logic has tests.
* No secrets are committed.
* Code is formatted.
* go test ./... passes.
* go vet ./... passes.
* Linting passes when configured.
* API behavior has been manually or automatically verified.
* No unrelated features were added.

⸻

60. Final Rule

Do not over-engineer Sasvyn Backend.

Prefer:

simple
explicit
idiomatic Go
secure
testable
maintainable

over:

clever
abstract
framework-heavy
distributed
prematurely optimized

When there are multiple technically valid solutions, prefer the solution with the smallest complexity that satisfies the actual Sasvyn requirement.

Build the backend incrementally.

Do not guess future requirements.

Do not invent architecture.

Wait for the requested module, understand its requirements, then implement it properly.
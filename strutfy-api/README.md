# strutfy-api

REST API backend for Strutfy — a multi-tenant project management platform.

## Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.25 |
| HTTP Router | [httprouter](https://github.com/julienschmidt/httprouter) |
| Database | PostgreSQL 17 |
| ORM | [Bun](https://bun.uptrace.dev/) |
| Auth | JWT (HS256) + bcrypt |
| Cache / Queue | Redis 7 |
| Containers | Docker Compose |

## Getting Started

### Prerequisites

- Go 1.25+
- Docker and Docker Compose

### 1. Start infrastructure

```bash
docker compose up -d
```

This brings up PostgreSQL 17 on `:5432` and Redis 7 on `:6379`.

### 2. Configure environment

```bash
cp .env .env.local
# edit .env.local as needed
```

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_PORT` | `8080` | API listen port |
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/strutfy?sslmode=disable` | PostgreSQL DSN |
| `REDIS_URL` | `localhost:6379` | Redis address |
| `JWT_SECRET` | `dev-secret` | **Change in production** — HMAC-SHA256 signing key |
| `JWT_ISSUER` | `strutfy-api` | JWT `iss` claim |
| `STRIPE_SECRET_KEY` | — | Stripe (not yet implemented) |
| `POSTMARK_SERVER_TOKEN` | — | Email service (not yet implemented) |
| `OPENAI_API_KEY` | — | OpenAI (not yet implemented) |

### 3. Run the API

```bash
go run ./cmd/api
```

Migrations run automatically on startup. The API will be available at `http://localhost:8080`.

## API

### Authentication

All protected endpoints require a Bearer token in the `Authorization` header:

```
Authorization: Bearer <access_token>
```

Tokens expire after **15 minutes**.

### Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/health` | — | Health check |
| `POST` | `/v1/auth/register` | — | Register user + create organization |
| `POST` | `/v1/auth/login` | — | Authenticate and receive token |
| `POST` | `/v1/projects` | Required | Create project |
| `DELETE` | `/v1/projects/:id` | Required | Delete project |

#### POST /v1/auth/register

```json
{
  "company": "Acme Inc",
  "document": "12.345.678/0001-99",
  "name": "Jane Doe",
  "email": "jane@acme.com",
  "password": "s3cr3t"
}
```

Creates an organization and registers the user as `admin`. Returns a JWT access token.

#### POST /v1/auth/login

```json
{
  "email": "jane@acme.com",
  "password": "s3cr3t"
}
```

#### POST /v1/projects

Requires `X-Org-ID` header with the organization UUID.

```json
{
  "name": "My Project",
  "description": "Optional description"
}
```

## Architecture

```
cmd/api/            Application entry point
internal/
  app/              Dependency wiring and application lifecycle
  config/           Environment-based configuration
  database/         PostgreSQL connection, migration runner, embedded SQL
  domain/auth/      JWT claims and permission definitions
  http/
    handlers/       HTTP request/response layer
    middleware/     JWT validation, role and permission guards
    router.go       Route registration
  auth/             Registration and login business logic
  service/          JWT token service
  user/             User model and repository
  organization/     Organization model and repository
  project/          Project model, repository, and service
  workflow/         In-memory job queue with worker pool
  shared/response/  Common JSON response helpers
  realtime/         Placeholder — WebSocket / SSE
  integrations/     Placeholder — external service clients
```

### Multi-tenancy

All resources are scoped to an `Organization`. The organization UUID must be passed via the `X-Org-ID` request header for protected endpoints. JWTs also embed `org_id` to prevent cross-tenant access.

### Permissions

| Role | Permissions |
|------|------------|
| `admin` | `projects:read`, `projects:create`, `projects:update`, `projects:delete` |
| `member` | `projects:read` |

### Workflow Engine

A simple in-memory worker pool (4 workers, 100-job buffer) processes background jobs. Planned integrations include Stripe webhooks, transactional email, bank reconciliation, AI pipelines, and report generation.

### Database Migrations

SQL files are embedded into the binary via `//go:embed`. Applied migrations are tracked in the `schema_migrations` table and run automatically at startup.

## Project Status

Phase 1 complete: core authentication and project CRUD.

Planned:
- Token refresh (`POST /v1/auth/refresh`)
- Project list and detail endpoints
- Stripe billing integration
- Email notifications via Postmark
- AI-powered features
- Real-time events (WebSocket / SSE)

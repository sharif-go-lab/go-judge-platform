# Go Judge Platform

Go Judge Platform is a teaching-oriented online judge implemented in Go. It delivers end-to-end functionality for managing programming problems, running contests, and judging Go submissions inside isolated Docker containers. The repository includes both the HTTP application that students interact with and a lightweight code runner service for executing submissions.

## Table of contents

- [Features](#features)
- [Architecture](#architecture)
- [Tech stack](#tech-stack)
- [Repository layout](#repository-layout)
- [Configuration](#configuration)
- [Running with Docker](#running-with-docker)
- [Running locally for development](#running-locally-for-development)
- [Database schema and migrations](#database-schema-and-migrations)
- [Useful Go commands](#useful-go-commands)
- [Contributing](#contributing)

## Features

- **Authentication and sessions** – Cookie-based sessions keep users signed in. Passwords are hashed using bcrypt before being stored, and a default administrator account is seeded for first-time setups.【F:internal/db/db.go†L16-L67】
- **Role-based access control** – Administrators can publish problems and manage user roles, while regular users can create drafts and submit solutions.【F:internal/handler/admin.go†L15-L116】【F:internal/handler/question.go†L15-L195】
- **Problem lifecycle** – Problems can be drafted, edited, and published with time and memory limits, sample I/O, and ownership metadata.【F:internal/model/models.go†L16-L30】
- **Submission pipeline** – Users submit Go code which is recorded as `Pending`, then evaluated by the code runner and updated with the final verdict (Accepted, Wrong Answer, Runtime Error, etc.).【F:internal/model/models.go†L32-L67】【F:internal/handler/submission.go†L20-L153】
- **Self-service portals** – Users receive profile pages with per-problem statistics, submission history, and authoring dashboards. Administrators can promote/demote users directly from the UI.【F:internal/handler/profile.go†L14-L63】【F:internal/handler/admin.go†L15-L72】
- **Isolated code execution** – Submissions are executed inside ephemeral Docker containers with enforced CPU, memory, and networking limits to protect the host environment.【F:cmd/code-runner/main.go†L18-L118】
- **Configurable services** – All components read from `config.yaml` and can be driven by environment variables thanks to Viper, making deployments flexible.【F:internal/config/config.go†L9-L32】【F:config.yaml†L1-L12】

## Architecture

```
+--------------------+        +------------------+
|  HTTP server       | <----> | PostgreSQL       |
|  (cmd/server)      |        | (users, problems |
|  Gin templates     |        |  submissions)    |
+---------+----------+        +---------+--------+
          |                             ^
          v                             |
+---------+----------+        +---------+--------+
| Internal REST APIs | <----> | Code runner      |
| for judges         |        | (cmd/code-runner |
|                    |        |  Docker sandbox) |
+--------------------+        +------------------+
```

- The **web server** (Gin + templates) handles HTML rendering, authentication, and business logic for problems, submissions, and profiles.【F:cmd/server/main.go†L1-L55】
- The **code runner** exposes a `/run` endpoint that accepts source code and test limits, then spins up a constrained Docker container to compile and execute the program.【F:cmd/code-runner/main.go†L18-L118】
- Both components share the same configuration loader and talk to PostgreSQL via GORM models defined under `internal/model`.【F:internal/config/config.go†L9-L32】【F:internal/model/models.go†L5-L67】

## Tech stack

- Go 1.22+
- Gin web framework with HTML templates
- GORM ORM for PostgreSQL
- Docker for secure execution
- Viper for configuration management
- bcrypt for password hashing

## Repository layout

```
cmd/
  server/         # Web server entry point
  code-runner/    # Submission execution service
config.yaml       # Default configuration (used by docker-compose)
internal/
  config/         # Configuration bootstrap
  db/             # Database initialization and seeding
  handler/        # HTTP handlers grouped by feature
  middleware/     # Shared Gin middleware
  model/          # GORM models
migrations/       # SQL migrations (PostgreSQL)
static/, templates/ # Front-end assets and HTML templates
```

## Configuration

Configuration is read in this order: environment variables, `config.yaml`, then baked-in defaults. Nested keys use dot notation which maps to environment variables via `.` → `_` replacement (for example `SERVER_LISTEN`).【F:internal/config/config.go†L9-L32】

Key settings include:

| Key | Description | Default |
| --- | ----------- | ------- |
| `server.listen` | Address Gin listens on | `:8080` |
| `code_runner.listen` | Address code runner listens on | `:9000` |
| `code_runner.judge_dockerfile` | Path to the Dockerfile injected into the runner container | `./judge.Dockerfile` |
| `database.dsn` | PostgreSQL DSN string | `postgres://user:pass@localhost:5432/go_judge?sslmode=disable` |
| `session.secret` | Cookie signing key | `super-secret-key` |

When booting via Docker Compose the bundled `config.yaml` sets the DSN, ports, and secrets to match the containers.【F:config.yaml†L1-L12】

## Running with Docker

1. Ensure Docker (and optionally Docker Compose v2) is installed and that your user can access the Docker daemon (required by the code runner).
2. Copy `config.yaml` if you need to adjust secrets or ports before starting.
3. Start the stack:

   ```bash
   docker compose up --build
   ```

   This brings up PostgreSQL, the web server (`cmd/server`), and the code runner (`cmd/code-runner`).【F:docker-compose.yaml†L1-L38】

4. Visit http://localhost:8080 and log in using the seeded administrator (`admin` / `admin123`). Change the password immediately from the profile page.

The code runner container mounts `/var/run/docker.sock` read-only so it can launch sandboxed containers for each submission.【F:docker-compose.yaml†L23-L36】

## Running locally for development

1. **Start PostgreSQL** locally or in Docker. The DSN must match `database.dsn`.
2. **Set environment variables** or edit `config.yaml` to point at your database and set a `session.secret`.
3. **Run migrations** (see next section) if you prefer SQL files; otherwise the server will auto-migrate tables on startup.【F:internal/db/db.go†L24-L49】
4. **Run the web server**:

   ```bash
   go run ./cmd/server
   ```

   The server loads configuration, establishes DB connections, registers templates, and listens on `server.listen`.【F:cmd/server/main.go†L1-L55】

5. **Run the code runner** in a separate terminal:

   ```bash
   go run ./cmd/code-runner
   ```

   Keep Docker running in the background so submissions can be executed.【F:cmd/code-runner/main.go†L18-L118】

6. Access the UI at `http://localhost:8080`. The default admin user is seeded automatically if none exists.【F:internal/db/db.go†L34-L67】

## Database schema and migrations

- GORM auto-migrates the schema for users, problems, submissions, sessions, and test cases whenever the server boots.【F:internal/db/db.go†L28-L49】
- SQL migrations are provided under `migrations/` for environments where declarative migrations are preferred. These scripts cover table creation and later schema adjustments.【F:migrations/0001_create_users.up.sql†L1-L33】【F:migrations/0007_add_sample_fields_to_problems.up.sql†L1-L17】
- To run the SQL migrations manually you can use the [`golang-migrate`](https://github.com/golang-migrate/migrate) CLI:

  ```bash
  migrate -path migrations -database "$DATABASE_DSN" up
  ```

  Replace `$DATABASE_DSN` with your PostgreSQL connection string.

## Useful Go commands

The project does not ship with a `Makefile`, but you can use native Go tooling:

- `go run ./cmd/server` – start the web server with HTML templates.
- `go run ./cmd/code-runner` – start the sandboxed execution service.
- `go test ./...` – execute unit tests (add your own as the project evolves).

## Contributing

1. Fork the repository and clone it locally.
2. Create a feature branch and make your changes.
3. Run `go fmt ./...` and `go test ./...` to ensure code quality.
4. Submit a pull request describing your changes and screenshots when touching the UI templates.

Contributions that improve the judging pipeline, add language support, or tighten security are especially welcome!

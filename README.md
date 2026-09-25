# Go Judge Platform

[![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go&logoColor=white)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

An online judge for Go, written in Go. Users sign up, read published problems and submit Go solutions from the browser. A separate code-runner service runs each submission in a throwaway Docker container with no network access and a memory cap, then compares the program's output with the problem's expected output.

Any signed-in user can write a problem. It stays a draft until an admin publishes it. Admins also manage user roles.

## Features

- Registration and sign-in with bcrypt-hashed passwords and signed cookie sessions. You can sign in with either your username or your email.
- Problems with a statement, time and memory limits, and a sample input with its expected output. Authors can edit their own problems, and admins can edit any of them.
- Problems stay hidden until an admin publishes them. Regular users see published problems, 10 per page; admins see every problem and can publish or unpublish it.
- Go submissions that are judged asynchronously. A submission is saved as `Pending`, and its page shows the verdict once the runner answers.
- Sandboxed execution: each run gets its own `golang:1.24.2` container with `--network none`, the problem's memory limit, one CPU and at most 64 processes.
- Submission history for each user. Admins can open anyone's submissions.
- Public profile pages and an admin page for promoting and demoting users.
- A server-rendered UI built with Gin and Go HTML templates. There is no JavaScript framework.

## How it works

```
 browser ──HTTP──► server (cmd/server, :8080)  ──GORM──►  PostgreSQL
                     Gin + HTML templates                  users, problems,
                     cookie sessions                       submissions
                          │
                          │ POST /run  (code, sample input/output, limits)
                          ▼
                   code-runner (cmd/code-runner, :9000)
                          │
                          │ docker run --rm --network none --memory … golang:1.24.2
                          ▼
                   one container per submission
```

When a user submits code, the server stores a `Submission` with status `Pending` and redirects to its page. In the background it sends the code, the problem's sample input and expected output, and the problem's limits to the runner at `http://code-runner:9000/run`. It then stores whatever verdict comes back. A second background timer marks the submission `Failed` if it is still pending after `server.submission_time_out` seconds (100 in the bundled `config.yaml`).

The runner writes the code to `main.go` inside a fresh container and runs it with `go run`, piping the sample input to stdin. It picks the verdict from the exit code and the output:

| Outcome | Verdict |
|---|---|
| Exit code 0, nothing on stderr, and stdout equals the expected output (surrounding whitespace ignored) | `Accepted` |
| Exit code 0 and nothing on stderr, but the output differs | `Wrong Answer` |
| Exit code 0 with something written to stderr | `Runtime Error` |
| Exit code 137 (the container was killed, usually for exceeding the memory limit) | `Memory Limit Exceeded` |
| Exit code 139 (segmentation fault) | `Runtime Error` |
| Exit code 124 | `Time Limit Exceeded` |
| Any other non-zero exit code, including a failed compile | `Compilation Error` |
| The runner can't be reached, returns something unreadable, or doesn't answer in time | `Failed` (set by the server) |

The runner needs the Docker CLI and access to the host's Docker socket. In Docker Compose the socket is mounted into the runner container at `/var/run/docker.sock`.

## Quick start with Docker Compose

`docker-compose.yaml` starts PostgreSQL 15, the server and the code runner:

```bash
docker compose up --build
```

Then open http://localhost:8080. On first start the server creates an admin account:

| Username | Password |
|---|---|
| `admin` | `admin123` |

There's no page for changing a password yet, so don't expose the app to other people with this account in place.

The compose file builds the server from `cmd/server/Dockerfile` and the runner from `cmd/code-runner/Dockerfile`. Neither file is in the repository at the moment, because `.gitignore` ignores every file named `Dockerfile`. Add both files, and force-add them with `git add -f`, before running the command above from a fresh clone. The runner's image needs the Docker CLI, and it needs `judge.Dockerfile` at the path set in `code_runner.judge_dockerfile` (the runner rejects every run if it can't read that file).

## Running without Docker Compose

You still need Docker, both for PostgreSQL and for the containers that run submissions.

1. Start PostgreSQL:

   ```bash
   docker run -d --name go-judge-db -p 5432:5432 \
     -e POSTGRES_USER=myuser -e POSTGRES_PASSWORD=mysecurepassword -e POSTGRES_DB=go_judge \
     postgres:15
   ```

2. Start the server from the repository root, pointing it at that database:

   ```bash
   export DATABASE_DSN="host=localhost user=myuser password=mysecurepassword dbname=go_judge sslmode=disable"
   go run ./cmd/server
   ```

   The server creates or updates its tables with GORM on every start.

3. Start the runner in another terminal:

   ```bash
   CODE_RUNNER_JUDGE_DOCKERFILE=cmd/code-runner/judge.Dockerfile go run ./cmd/code-runner
   ```

   Pull the image once with `docker pull golang:1.24.2` so the first submission doesn't wait for the download.

4. The server always calls the runner at the hostname `code-runner`. Outside Compose, make that name resolve to your machine:

   ```bash
   echo "127.0.0.1 code-runner" | sudo tee -a /etc/hosts
   ```

## Configuration

Both programs read `config.yaml` from the working directory (or `./config/`), and environment variables override it. An environment variable's name is the key in upper case with dots replaced by underscores, so `database.dsn` becomes `DATABASE_DSN`.

| Key | Used by | Default | Value in `config.yaml` |
|---|---|---|---|
| `server.listen` | server | `:8080` | `:8080` |
| `server.submission_time_out` | server | none | `100` (seconds before a pending submission becomes `Failed`) |
| `database.dsn` | server | `postgres://user:pass@localhost:5432/go_judge?sslmode=disable` | points at the `db` service from Compose |
| `session.secret` | server | `super-secret-key` | a fixed random string (replace it for any real deployment) |
| `code_runner.listen` | runner and server | `:9000` | `:9000` |
| `code_runner.judge_dockerfile` | runner | none | `./judge.Dockerfile` |

## Pages

| Path | Access | What it does |
|---|---|---|
| `/` | public | Landing page. Signed-in users are redirected to their profile. |
| `/auth/login`, `/auth/register`, `/auth/logout` | public | Sign in, sign up (alphanumeric username, valid email, password of 6+ characters), sign out. |
| `/questions/` | public | Published problems, 10 per page. Admins also see drafts. |
| `/questions/:id` | public | Problem statement, limits, sample input and output, and the submit form. Drafts are visible only to their author and to admins. |
| `/profile/:username` | public | A user's profile. |
| `/profile` | signed in | Your own profile. |
| `/questions/create`, `/questions/my`, `/questions/edit/:id` | signed in | Write a problem, list the problems you wrote, edit one of them. |
| `/submissions/`, `/submissions/:id` | signed in | Your submissions and a single submission with its code and verdict. |
| `POST /submissions/submit/:question_id` | signed in | Submit code for a problem. |
| `/admin/users` | admin | List users; promote or demote them. |
| `POST /admin/questions/publish/:id`, `POST /admin/questions/unpublish/:id` | admin | Publish or unpublish a problem. |

The runner has one endpoint, `POST /run`, which takes JSON with `code`, `sample_input`, `sample_output`, `time_limit` (ms) and `memory_limit` (MB), and returns `{"result": "<verdict>"}`, plus `stdout` or `stderr` for some verdicts.

## Database

The server creates and updates its tables (`users`, `problems`, `submissions`, `sessions`, `test_cases`) with GORM's auto-migration on startup. The `migrations/` folder has the same schema as plain SQL files, in case you prefer to manage it with [golang-migrate](https://github.com/golang-migrate/migrate):

```bash
migrate -path migrations -database "$DATABASE_URL" up
```

## Project layout

```
cmd/
  server/                 web application (Gin, templates, sessions)
  code-runner/            judging service and judge.Dockerfile
internal/
  config/                 Viper setup: defaults, config.yaml, environment
  db/                     connection, auto-migration, default admin
  handler/                pages and form handlers, grouped by feature
  middleware/             sign-in and admin checks
  model/                  GORM models and the runner's request type
migrations/               SQL schema for golang-migrate
templates/                HTML templates (layout, auth, questions, submissions, profile, admin)
static/                   CSS and a small script
config.yaml               configuration used with Docker Compose
docker-compose.yaml       PostgreSQL, server and runner
```

## Limitations

- Only Go is supported, and each submission is checked against the problem's single sample input and output. The `test_cases` table exists but isn't used for judging yet.
- The time limit is sent to the runner but not enforced inside the container, so a submission that never finishes ends up as `Failed` after `server.submission_time_out` seconds instead of `Time Limit Exceeded`.
- A panic (exit code 2) is reported as `Compilation Error`, because every non-zero exit code other than 124, 137 and 139 maps to that verdict.
- The profile's solved count looks for the status `accepted`, while the runner stores `Accepted`, so it currently shows 0 solved problems.
- The runner has no authentication, and whoever can reach port 9000 can start containers through the host's Docker socket. Keep that port private.
- Sessions live in a signed cookie. The `sessions` table is created but not used.

## Authors

- [Kasra Siavashpour](https://github.com/kasra-sia)
- [Ardalan Siavashpour](https://github.com/Ardalan-Sia)

## License

[MIT](LICENSE)

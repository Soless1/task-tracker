# Task Tracker

Task management application written in **Go** with CLI and Web interfaces.

The project demonstrates **Clean Architecture**, REST API development, PostgreSQL, code generation and multiple storage implementations.

## Features

* Create, update and delete tasks
* Change task status
* Filter tasks by status
* CLI interface
* Web interface
* REST API
* JSON and PostgreSQL storage
* Swagger documentation
* Database migrations
* Docker support
* Unit tests with mocks

## Tech Stack

* Go
* PostgreSQL
* pgx / pgxpool
* `net/http`
* sqlc
* Goose
* Swagger / Swaggo
* Testify
* GoMock
* HTML / CSS / JavaScript
* Docker Compose
* Task

## Architecture

The project follows **Clean Architecture** principles.

```text
CLI ──────┐
          │
Web ──────┤
          ▼
       Usecases
          │
          ▼
        Domain
          │
          ▼
   TaskRepository
      ▲       ▲
      │       │
    JSON   PostgreSQL
```

The application depends on the `TaskRepository` interface, so the storage implementation can be changed without modifying the usecases.

## Project Structure

```text
task-tracker/
├── cmd/
│   └── task-cli/
├── internal/
│   ├── domain/
│   ├── usecase/
│   ├── repository/
│   │   ├── json/
│   │   └── postgres/
│   └── transport/
│       └── http/
├── migrations/
├── docs/
├── web/
├── Taskfile.yml
├── docker-compose.yml
├── sqlc.yaml
├── go.mod
└── go.sum
```

## Requirements

Install the following tools before running the project:

* Go
* Docker and Docker Compose
* Task
* sqlc
* Swag

The Go dependencies, including **Testify** and **GoMock**, are managed through `go.mod` and `go.sum`.

No separate dependency installation command is required.

## Setup and Run

After installing the required tools, run:

```bash
task setup
```

This starts PostgreSQL, generates sqlc and Swagger code, and starts the Web application with PostgreSQL storage.

The application will be available at:

```text
http://localhost:8080
```

Swagger UI:

```text
http://localhost:8080/swagger/index.html
```

## Task Commands

Start the project:

```bash
task setup
```

Start PostgreSQL:

```bash
task postgres
```

Generate sqlc and Swagger:

```bash
task generate
```

Run tests:

```bash
task test
```

## Testing

Tests use **Testify** for assertions and **GoMock** for creating and configuring mocks.

Run all tests:

```bash
task test
```

## CLI

The application can also be run in CLI mode.

### JSON storage

```bash
go run ./cmd/task-cli -mode=cli -storage=json add "Learn Go"
go run ./cmd/task-cli -mode=cli -storage=json list
go run ./cmd/task-cli -mode=cli -storage=json mark-in-progress 1
go run ./cmd/task-cli -mode=cli -storage=json mark-done 1
go run ./cmd/task-cli -mode=cli -storage=json update 1 "Learn Go and PostgreSQL"
go run ./cmd/task-cli -mode=cli -storage=json delete 1
```

### PostgreSQL storage

```bash
go run ./cmd/task-cli -mode=cli -storage=postgres list
```

The default application mode is **Web + PostgreSQL**, so running:

```bash
go run ./cmd/task-cli
```

starts the Web application using PostgreSQL.

## Web API

| Method | Endpoint                      | Description         |
| ------ | ----------------------------- | ------------------- |
| GET    | `/api/tasks`                  | Get tasks           |
| GET    | `/api/tasks?status=done`      | Filter tasks        |
| GET    | `/api/tasks/{id}`             | Get task            |
| POST   | `/api/tasks`                  | Create task         |
| PUT    | `/api/tasks/{id}`             | Update task         |
| DELETE | `/api/tasks/{id}`             | Delete task         |
| POST   | `/api/tasks/{id}/done`        | Mark as done        |
| POST   | `/api/tasks/{id}/in-progress` | Mark as in progress |

## Database

PostgreSQL runs using Docker Compose.

Database connections are managed with `pgxpool`.

Database schema changes are handled by **Goose**. Migrations are also checked automatically when the application starts.

SQL queries are compiled into type-safe Go code using **sqlc**.

## Code Generation

All project code generation is handled through Task:

```bash
task generate
```

This generates:

* PostgreSQL code using sqlc
* Swagger documentation using Swaggo

Generated Swagger files are stored in:

```text
docs/
├── docs.go
├── swagger.json
└── swagger.yaml
```

## Development

Run tests:

```bash
task test
```

Generate code:

```bash
task generate
```

Start PostgreSQL:

```bash
task postgres
```

Start the application:

```bash
go run ./cmd/task-cli
```

## License

Educational and portfolio project.

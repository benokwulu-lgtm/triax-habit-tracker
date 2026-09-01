# Triax Habit Tracker

A habit tracking application built with Go and PostgreSQL, developed as a learning project in backend engineering.

## Tech Stack

- Go (standard library `net/http`, `database/sql`)
- PostgreSQL
- `github.com/lib/pq` (Postgres driver)

## Prerequisites

- Go 1.21+
- PostgreSQL installed and running locally

## Environment Variables

| Variable      | Default                   | Notes                          |
|---------------|----------------------------|---------------------------------|
| `DB_HOST`     | `localhost`                |                                 |
| `DB_PORT`     | `5432`                     |                                 |
| `DB_USER`     | `postgres`                 | Set to your OS/Postgres user   |
| `DB_PASSWORD` | *(empty)*                  | Required on Windows installs   |
| `DB_NAME`     | `triax_habit_tracker`      |                                 |
| `DB_SSLMODE`  | `disable`                  | Local dev only                 |

## Setup

1. Create the database:
```bash
   createdb triax_habit_tracker
```
2. Run migrations:
```bash
   psql triax_habit_tracker -f migrations/0001_create_users_and_habits.sql
```
3. Set environment variables as needed for your machine (see table above).
4. Run the server:
```bash
   go run cmd/api/main.go
```

Server starts on `:8080`.
# Mini Project

A small helpdesk application for learning Go, React, TypeScript, PostgreSQL, Redis, and JWT authentication.

## Stack

- Backend: Go `net/http`, GORM, PostgreSQL
- Authentication: JWT (HS256), Argon2 password hashing, Redis token blacklist
- Frontend: React + TypeScript + Bun

## Features

- User registration and login
- Access-token verification, refresh, and sign-out/revocation
- Create, list, retrieve, update, and delete tickets
- Ticket operations are scoped to the authenticated user's ID
- Request-body and response DTO validation

## Requirements

- Go version declared in `backend/go.mod`
- PostgreSQL
- Redis
- Bun (only if running the frontend)

## Backend setup

1. Create a PostgreSQL database named `mini_project`, or update the connection string.
2. Start Redis.
3. Copy `backend/.env.example` to `backend/.env`.
4. Set `DATABASE_URL`, `REDIS_URL`, and `JWT_SECRET` in `backend/.env`.
   Do not commit `.env` or share its values.
   `JWT_SECRET` must be non-empty; the application does not enforce a 32-character minimum.
5. Run the backend:

   ```sh
   cd backend
   go mod download
   go test ./...
   go vet ./...
   go build ./...
   go run .
   ```

The API listens on `http://localhost:8080`.

GORM `AutoMigrate` runs when the backend starts, so a separate migration command is not currently required. For production, use reviewed versioned migrations rather than relying on automatic schema migration.

## Environment variables

| Variable       | Purpose                      | Example                                                                      |
| -------------- | ---------------------------- | ---------------------------------------------------------------------------- |
| `DATABASE_URL` | PostgreSQL connection string | `postgresql://postgres:postgres@localhost:5432/mini_project?sslmode=disable` |
| `REDIS_URL`    | Redis connection string      | `redis://localhost:6379/0`                                                   |
| `JWT_SECRET`   | HS256 signing secret         | `replace-this-with-a-private-random-secret`                                  |

## API routes

- `POST /api/auth/signup`
- `POST /api/auth/signin`
- `GET /api/auth/verify`
- `POST /api/auth/refresh`
- `POST /api/auth/signout`
- `POST /api/tickets`
- `GET /api/tickets`
- `GET /api/tickets/{id}`
- `PUT /api/tickets/{id}`
- `DELETE /api/tickets/{id}`

Protected ticket routes expect `Authorization: Bearer <token>`. A ticket ID must be a UUID. Empty update bodies are rejected. A ticket belonging to another user is returned as not found.

## Frontend setup

```sh
cd frontend
bun install
bun run dev
```

## Docker

No Dockerfile or Compose configuration is included in this repository version. Start PostgreSQL and Redis separately, or add deployment-specific container configuration before using Docker.

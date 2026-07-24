# GuiSIS API Gateway

The backend API gateway for the PUP Student Guidance System Capstone. Built
with Go, Gin, and MySQL/Redis.

## Architecture

This project is structured around **Purist Clean Architecture** and
**Vertical Slicing**. All modules are grouped inside `internal/features/`
by feature slices rather than technical layers.

For details on coding practices, architectural standards, and naming, please
refer to [CONTRIBUTING.md](CONTRIBUTING.md).

## Tech Stack

- **Language**: Go (v1.22+)
- **HTTP Framework**: Gin
- **Database Access**: sqlx + MySQL
- **Caching & Sessions**: Redis
- **Documentation**: Swagger (swag)
- **Live Reloading**: Air

## Setup Instructions

### Prerequisites

- Go v1.22+ installed
- MySQL Server (v8.0+)
- Redis Server
- `golang-migrate` CLI installed (for migrations)

### Installation

1. Clone the repository and navigate to the `guisis-api` folder.
2. Initialize configuration:
   ```bash
   copy .env.example .env
   ```
3. Update the values in `.env` (such as DB credentials and Redis host).

### Database Migrations

You can manage the database schema using the following `make` commands:

- Run all migrations:
  ```bash
  make migrate-up
  ```
- Rollback migrations:
  ```bash
  make migrate-down
  ```
- Seed database (initial data and locations):
  ```bash
  make seed-up
  make locations
  ```
- Total Database Reset & Reseed:
  ```bash
  make refresh
  ```

### Running Locally

- **With hot-reloading (requires Air installed)**:
  ```bash
  air
  ```
- **Standard execution**:
  ```bash
  go run cmd/api/main.go
  ```

### Docker Deployment

To build and run the backend using Docker Compose:

```bash
# For development/testing
make compose-up

# For staging environment
make compose-staging

# For production environment
make compose-prod
```

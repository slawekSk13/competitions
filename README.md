# Competition Platform

A full-stack application for managing competitions, built with modern technologies and containerized using Docker.

## Project Structure

```
.
├── backend/           # Backend service (Go)
├── frontend/         # Frontend application (React + TypeScript + Vite)
├── database/         # Database initialization scripts
├── secrets/          # Secret files (not tracked in git)
├── docker-compose.yml # Docker services configuration
└── load-secrets.sh   # Script to load secrets and start services
```

## Prerequisites

- Docker and Docker Compose
- Node.js (for local development)
- Git

## Getting Started

### 1. Clone the Repository

```bash
git clone <repository-url>
cd competitions
```

### 2. Set Up Secrets

Create the following files in the `secrets` directory:

```
secrets/
├── db_name.txt       # Database name
├── db_user.txt       # Database user
├── db_password.txt   # Database password
├── pgadmin_email.txt # PgAdmin email
└── pgadmin_password.txt # PgAdmin password
```

Example content for each file:

- `db_name.txt`: `competition_app_db`
- `db_user.txt`: `competition_app_user`
- `db_password.txt`: `your_secure_password`
- `pgadmin_email.txt`: `admin@example.com`
- `pgadmin_password.txt`: `your_pgadmin_password`

### 3. Start the Application

You can start the application in two ways:

#### Option 1: Using the load-secrets script (Recommended)

```bash
# Start with existing data
./load-secrets.sh

# Start with clean data (removes all volumes)
./load-secrets.sh --clean
```

#### Option 2: Using Docker Compose directly

```bash
# Load secrets into environment
export PGADMIN_DEFAULT_EMAIL=$(cat ./secrets/pgadmin_email.txt)
export PGADMIN_DEFAULT_PASSWORD=$(cat ./secrets/pgadmin_password.txt)

# Start services
docker compose up --build
```

## Services

The application consists of the following services:

- **Frontend**: React + TypeScript application running on `http://localhost:7788`
- **Backend**: Go API running on `http://localhost:8080`
- **PostgreSQL**: Database running on port 5432
- **Redis**: Cache server running on port 6379
- **PgAdmin**: Database management interface running on `http://localhost:5050`

## Development

### Frontend Development

The frontend code is located in the `frontend` directory. It's built with React, TypeScript, and Vite. Changes to the code will be automatically reflected in the development environment.

### Backend Development

The backend code is located in the `backend` directory. It's written in Go and uses Air for hot reloading during development. The service will automatically restart when changes are detected.

### Database Management

You can access PgAdmin at `http://localhost:5050` using the credentials specified in your secrets files.

## Environment Variables

The following environment variables are used in the application:

- `PGADMIN_DEFAULT_EMAIL`: Email for PgAdmin login
- `PGADMIN_DEFAULT_PASSWORD`: Password for PgAdmin login
- `DB_HOST`: Database host (default: postgres)
- `DB_PORT`: Database port (default: 5432)
- `REDIS_HOST`: Redis host (default: redis)
- `REDIS_PORT`: Redis port (default: 6379)
- `SERVER_PORT`: Backend server port (default: 8080)
- `VITE_API_URL`: Frontend API URL (default: http://backend:8080)

## Contributing

1. Create a new branch for your feature
2. Make your changes
3. Submit a pull request

## Security Notes

- Never commit the `secrets` directory to version control
- Keep your secret files secure and use strong passwords
- Regularly rotate database and PgAdmin credentials
- Use environment variables for sensitive configuration in production

## Troubleshooting

If you encounter any issues:

1. Check if all services are running: `docker compose ps`
2. View service logs: `docker compose logs [service-name]`
3. Ensure all secret files are properly configured
4. Try cleaning and rebuilding: `./load-secrets.sh --clean`

#!/bin/bash

# Function to show usage
show_usage() {
    echo "Usage: $0 [--clean]"
    echo "  --clean    Remove all volumes before starting"
    exit 1
}

# Parse arguments
CLEAN=false
for arg in "$@"; do
    case $arg in
        --clean)
            CLEAN=true
            shift
            ;;
        *)
            show_usage
            ;;
    esac
done

# Stop and remove containers if running
echo "Stopping containers..."
docker compose down

# Remove volumes if --clean flag is set
if [ "$CLEAN" = true ]; then
    echo "Removing volumes..."
    docker compose down -v
fi

# Load pgAdmin credentials from secrets
echo "Loading secrets..."
export PGADMIN_DEFAULT_EMAIL=$(cat ./secrets/pgadmin_email.txt)
export PGADMIN_DEFAULT_PASSWORD=$(cat ./secrets/pgadmin_password.txt)

# Run docker compose
echo "Starting services..."
docker compose up --build 
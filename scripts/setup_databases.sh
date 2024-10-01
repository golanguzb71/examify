#!/bin/bash
set -e

setup_database() {
    local db_name=$1
    local migration_file=$2

    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
        SELECT 'CREATE DATABASE $db_name'
        WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$db_name')\gexec
EOSQL

    if [ -f "$migration_file" ]; then
        psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$db_name" -f "$migration_file"
        echo "$db_name setup and migrations completed."
    else
        echo "Migration file $migration_file not found. Skipping migrations for $db_name."
    fi
}

setup_database "ielts_database" "/docker-entrypoint-initdb.d/ielts_service_up.sql"

setup_database "user_database" "/docker-entrypoint-initdb.d/user_service_up.sql"

echo "All database setups and migrations completed."
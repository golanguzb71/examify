#!/bin/bash

DB_USER="postgres"
DB_PASSWORD="postgres"

export PGPASSWORD=$DB_PASSWORD

echo "Dropping user_database..."
psql -U $DB_USER -c "DROP DATABASE IF EXISTS user_database;"

echo "Dropping ielts_database..."
psql -U $DB_USER -c "DROP DATABASE IF EXISTS ielts_database;"

echo "Dropping integration_database..."
psql -U $DB_USER -c "DROP DATABASE IF EXISTS integration_database;"

echo "Databases dropped."

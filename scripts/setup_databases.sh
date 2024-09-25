#!/bin/bash
DB_USER="postgres"
DB_PASSWORD="postgres"
export PGPASSWORD=$DB_PASSWORD
echo "Creating user_database..."
psql -U $DB_USER -c "CREATE DATABASE user_database OWNER $DB_USER;"
echo "Creating ielts_database..."
psql -U $DB_USER -c "CREATE DATABASE ielts_database OWNER $DB_USER;"
echo "Creating integration_database..."
psql -U $DB_USER -c "CREATE DATABASE integration_database OWNER $DB_USER"
echo "Applying migrations for user_database..."
psql -U $DB_USER -d user_database -f /home/elon/GolandProjects/examify/UserService/migrations/user_service_up.sql
echo "Applying migrations for ielts_database..."
psql -U $DB_USER -d ielts_database -f /home/elon/GolandProjects/examify/IeltsService/migrations/ielts_service_up.sql
echo "Database setup and migrations completed."
psql -U $DB_USER -d integration_database -f /home/elon/GolandProjects/examify/IntegrationService/migrations/integration_service_up.sql
echo "Database setup and migrations completed."

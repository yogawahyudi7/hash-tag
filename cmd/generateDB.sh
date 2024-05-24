#!/bin/bash
DB_HOST="localhost"
DB_PORT="5433" 
DB_NAME="postgres"
DB_PASSWORD="mocci"
SSL_MODE="disable" 

result=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" --set=sslmode="$SSL_MODE" -f sql/query.sql)

echo "$result"

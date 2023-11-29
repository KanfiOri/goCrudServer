#!/bin/bash

# Wait for the database to be ready
until nc -z my_postgres 5432; do
  echo "Waiting for PostgreSQL to start..."
  sleep 1
done

echo "PostgreSQL started. Starting the application."
./main
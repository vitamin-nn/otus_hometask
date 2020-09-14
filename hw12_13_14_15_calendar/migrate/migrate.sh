#!/bin/sh

DSN="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_DB_HOST:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable"

for i in $(seq 1 5); do
    goose -dir migrations postgres $DSN up && break
    sleep 1
done

#!/bin/sh

: "${DATABASE_USER:?DATABASE_USER not set}"
: "${DATABASE_PASSWORD:?DATABASE_PASSWORD not set}"
: "${DATABASE_HOST:?DATABASE_HOST not set}"
: "${DATABASE_PORT:?DATABASE_PORT not set}"

cat << EOM
# Up:
# migrate -url postgres://$DATABASE_USER:$DATABASE_PASSWORD@$DATABASE_HOST:$DATABASE_PORT?sslmode=disable -path ./migrations up
# Down:
# migrate -url postgres://$DATABASE_USER:$DATABASE_PASSWORD@$DATABASE_HOST:$DATABASE_PORT?sslmode=disable -path ./migrations down
EOM

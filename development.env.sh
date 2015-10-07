#!/bin/sh

env_script="$(basename "$0")"

cat << EOM
export DATABASE_USER="developer"
export DATABASE_PASSWORD="password"
export DATABASE_NAME="developer"
export DATABASE_HOST="amercia_db"
export DATABASE_PORT="5432"
export DATABASE_SSL="disable"
export API_PORT="8088"
export WEB_CLIENT_ID="801574721267-8ocanqgcgln83r5s2bdpk5imu78r2ouk.apps.googleusercontent.com"

# Run this command to configure your environment:
# eval "\$(./$env_script)"
EOM

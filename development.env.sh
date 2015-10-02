#!/bin/sh

u='export DATABASE_USER="developer"\n'
p='export DATABASE_PASSWORD="password"\n'
n='export DATABASE_NAME="develop"\n'
h='export DATABASE_HOST="db"\n'
s='export DATABASE_SSL="disabled"\n'
a='export API_PORT="8080"'

echo $u$p$n$h$s$a;

echo '# Run this command to configure your environment:
# eval "$(./development.env.sh)"';

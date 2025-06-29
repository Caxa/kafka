#!/bin/sh

host_with_port="$1"
host=$(echo "$host_with_port" | cut -d':' -f1)
shift

until ping -c 1 "$host" >/dev/null 2>&1 && PGPASSWORD=postgres psql -h "$host" -U "postgres" -d "ordersdb" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 2
done

>&2 echo "Postgres is up - executing command"
exec "$@"

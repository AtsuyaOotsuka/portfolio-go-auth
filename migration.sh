#!/bin/bash

set -e

mode=$1
opt=$2

# Load .env
source .env
DSN="mysql://${DB_USER}:${DB_PASS}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}"

MIGRATE_CMD="docker compose run --rm migrate"

function run_migrate() {
  $MIGRATE_CMD -path=/migrations -database="${DSN}" "$@"
}

function to_snake_case() {
  echo "$1" | sed -r 's/([A-Z])/_\L\1/g' | sed 's/^_//'
}

case "$mode" in
  "up")
    run_migrate up
    ;;

  "down")
    if [[ "$opt" == "all" ]]; then
      run_migrate down
    else
      : ${opt:=1}
      run_migrate down "$opt"
    fi
    ;;

  "create")
    name="$opt"
    if [ -z "$name" ]; then
      echo "Please provide a name."
      exit 1
    fi
    $MIGRATE_CMD create -ext sql -dir /migrations "$name"
    ;;
  *)
    echo "Usage: $0 {up|down|create} [name|steps]"
    exit 1
    ;;
esac

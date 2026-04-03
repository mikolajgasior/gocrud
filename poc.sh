#!/usr/bin/env bash
#
# poc.sh – Proof of Concept
# ----------------------------------------------
# Commands:
#   create, destroy,
#   check-connection, list-tables,
#   compile, run, test-api
# ----------------------------------------------

COMPOSE_FILE="poc-docker-compose.yml"
SERVICE_NAME="db"          # name of the service in the compose file

# PostgreSQL connection info
PGUSER="myuser"
PGPASSWORD="mypassword"
PGDATABASE="mydb"
PGHOST="localhost"
PGPORT="5432"

# POC app
POC_SRC="./cmd/poc"
POC_BIN="./poc"


# Env vars
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=myuser
export DB_PASS=mypassword
export DB_NAME=mydb
export POC_PORT=8080
export LOG_LEVEL=debug

# Run psql without exposing the password on the command line
run_psql() {
  PGPASSWORD="${PGPASSWORD}" psql \
    -h "${PGHOST}" -p "${PGPORT}" -U "${PGUSER}" -d "${PGDATABASE}" "$@"
}

# Print usage information
usage() {
  cat <<EOF
Usage: $0 <command>

Commands:
  create              docker compose up -d
  destroy             down → delete volume (leaves everything stopped)
  check-connection    test a simple SELECT against the DB
  list-tables         list DB tables
  compile             go build -o ${POC_BIN} ${POC_SRC}
  run                 execute ${POC_BIN}
  test-api            test api
EOF
  exit 1
}

# ------------------ Command handling -------

case "$1" in
  create)
    docker compose -f "${COMPOSE_FILE}" up -d
    ;;

  destroy)
    docker compose -f "${COMPOSE_FILE}" down
    docker volume rm "$(docker volume ls -qf name=pg_data)" 2>/dev/null || true
    ;;

  check-connection)
    if run_psql -c "SELECT 1;" >/dev/null 2>&1; then
      echo "✅ Connection succeeded."
    else
      echo "❌ Unable to connect to the database."
    fi
    ;;

  list-tables)
    echo "Tables in database \"${PGDATABASE}\":"
    run_psql -Atc "\dt"
    ;;

  compile)
    echo "🔨 Compiling POC application..."
    go build -o "${POC_BIN}" "${POC_SRC}"
    if [[ $? -eq 0 ]]; then
      echo "✅ Built ${POC_BIN}"
    else
      echo "❌ Build failed"
      exit 1
    fi
    ;;

  run)
    if [[ ! -x "${POC_BIN}" ]]; then
      echo "⚠️  Binary ${POC_BIN} not found or not executable."
      echo "Run '$0 compile' first."
      exit 1
    fi
    echo "🚀 Running POC via ${POC_BIN} ..."
    export PORT=${POC_PORT}
    export CREATE_TABLES="true"
    "${POC_BIN}"
    ;;

  test-api)
    for i in $(seq 1 321); do
      curl -X PUT \
        -d "{\"title\": \"Title $i\", \"description\": \"Desc $i\", \"uri\": \"page_$i\", \"status_id\": 1, \"created_at\": 0, \"created_by\": 0, \"modified_at\": 0, \"modified_by\": 0}" \
        http://localhost:8080/v0/crudapi/pages/;
    done

    num=$(curl -q http://localhost:8080/v0/crudapi/pages/?limit=13 2>/dev/null | jq '.data[].id' | wc -l)
    if [[ "$num" != "13" ]]; then
      echo "❌ Listing failed"
      exit 1
    fi

    code=$(curl -X DELETE http://localhost:8080/v0/crudapi/pages/11 2>/dev/null | jq -r '.code')
    if [[ "$code" != "SUCCESS" && "$code" != "NOT_FOUND" ]]; then
      echo "❌ Delete failed"
      exit 1
    fi

    for i in $(seq 20 27); do
      curl -X PUT \
        -d "{\"title\": \"Tytul $i\", \"description\": \"Opis $i\", \"uri\": \"page_$i\", \"status_id\": 3, \"created_at\": 0, \"created_by\": 0, \"modified_at\": 0, \"modified_by\": 0}" \
        http://localhost:8080/v0/crudapi/pages/$i;
    done

    num=$(curl -q http://localhost:8080/v0/crudapi/pages/?filter_val_StatusID=3 2>/dev/null | jq '.data[].id' | wc -l)
    if [[ "$num" != "8" ]]; then
      echo "❌ Update failed"
      exit 1
    fi
    ;;
  *)
    usage
    ;;
esac

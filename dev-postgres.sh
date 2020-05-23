# cmd & PostgreSQL (Docker)

# stop docker container
trap 'docker stop postgres-1' EXIT

# run docker container
# docker run --rm -d --name postgres-1 -p 127.0.0.1:5432:5432 --env-file postgres/.env -v "$(pwd)"/postgres/pgdata:/var/lib/postgresql/data postgres:12-alpine &&
docker run --rm -d --name postgres-1 --network host --env-file "$(pwd)"/postgres/.env -v "$(pwd)"/postgres/pgdata:/var/lib/postgresql/data postgres:12-alpine &&

# run go command
export $(cat dev-postgres.env | xargs) &&
export WEB_PUBLIC_DIRECTORY=$DIR/data/public
export SERVER_LOG_FILE=$DIR/data/log/default.log
export API_SERVER_SEED_FILE=$DIR/data/seed/default.json
go run ./cmd

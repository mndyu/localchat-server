# cmd & PostgreSQL (Docker)

# run docker container
# docker run --rm -d --name postgres-1 -p 127.0.0.1:5432:5432 --env-file postgres/.env -v "$(pwd)"/postgres/pgdata:/var/lib/postgresql/data postgres:12-alpine &&
docker run --rm -d --name postgres-1 --network host --env-file "$(pwd)"/postgres/.env -v "$(pwd)"/postgres/pgdata:/var/lib/postgresql/data postgres:12-alpine &&

# run go command
export $(cat dev-postgres.env | xargs) &&
go run ./cmd

# stop docker container
docker stop postgres-1

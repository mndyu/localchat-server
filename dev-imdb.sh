# cmd & SQLite (in-memory)

# run go command
export $(cat dev-imdb.env | xargs) &&
go run ./cmd

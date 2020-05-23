# cmd & SQLite (in-memory)

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )" # このスクリプトの場所

# run go command
export $(cat dev-imdb.env | xargs) &&
export WEB_PUBLIC_DIRECTORY=$DIR/data/public
export SERVER_LOG_FILE=$DIR/data/log/default.log
export API_SERVER_SEED_FILE=$DIR/data/seed/default.json
go run ./cmd

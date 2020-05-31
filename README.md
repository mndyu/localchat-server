# localchat-server

- Go 1.14
- Docker


## 実行
```
# Postgres + Minio
docker-compose up

# SQLite3
./dev-imdb.sh

# Postgres (途中)
./dev-postgres.sh
```


## env
```
DB_TYPE=<RDBMS名>
DB_URL=<接続URL>
DB_ADDRESS=<DBアドレス>
DB_DATABASE=<DB名>
DB_USER=<ユーザ名>
DB_PASSWORD=<ユーザのパスワード>
API_SERVER_ADDRESS=<サーバのアドレス (`:1324`)>
API_SERVER_LOG_FILE=<ログファイル>
API_SERVER_SEED_FILE=<seedファイルの場所>
```


------------------------------------------------------------------------

## TODO
- [ ] DB 接続
- [ ] スキーマ定義
- [ ] APIサーバ
- [ ] websocket
- [ ] Minio
- [ ] Docker Hub & GitHub Actions


------------------------------------------------------------------------
## curl
```
curl \
  -X GET \
  http://localhost:1324/api/v1/users

```

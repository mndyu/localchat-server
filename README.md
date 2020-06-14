# localchat-server

- Go 1.13
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


### Postgres
```
docker exec -ti postgres-1 psql localchat postgres
```


## env
```
APP_MODE=<production or development>
API_SERVER_ADDRESS=<サーバのアドレス (`:1324`), host, port 上書き>
API_SERVER_HOST=<サーバのホスト (`0.0.0.0`)>
API_SERVER_PORT=<サーバのポート (`1324`)>
API_SERVER_LOG_DIR=<ログのディレクトリ>
API_SERVER_SEED_FILE=<seedファイルのパス>
DB_TYPE=<RDBMS名>
DB_URL=<接続URL (address, database, user, password 上書き)>
DB_ADDRESS=<DBアドレス>
DB_DATABASE=<DB名>
DB_USER=<ユーザ名>
DB_PASSWORD=<ユーザのパスワード>
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
  -X POST \
  http://localhost:18000/api/v1/users \
  -H 'Content-Type: application/json' \
  -d '{"ip_address": "::1", "name": "dade", "pc_name": "oeoe"}'

curl \
  -X POST \
  http://localhost:18000/api/v1/groups \
  -H 'Content-Type: application/json' \
  -d '{"name": "huehue"}'

curl \
  -X POST \
  http://localhost:18000/api/v1/groups/1/members \
  -H 'Content-Type: application/json' \
  -d '{"mysqlf": true, "user_id": 1}'

curl \
  -X POST \
  http://localhost:18000/api/v1/messages \
  -H 'Content-Type: application/json' \
  -d '{"group_id": 1, "body": "heheheue", "to": [1,2,3]}'



curl -X GET http://localhost:18000/api/v1/users
curl -X GET http://localhost:18000/api/v1/users/1/groups

curl -X GET http://localhost:18000/api/v1/messages

curl -X GET http://localhost:18000/api/v1/groups
curl -X GET http://localhost:18000/api/v1/groups/1/messages
curl -X GET http://localhost:18000/api/v1/groups/1/members
curl -X GET http://localhost:18000/api/v1/groups/1/members/1/messages

curl -X DELETE http://localhost:18000/api/v1/messages/1


```

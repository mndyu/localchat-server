# localchat-server


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
SERVER_LOG_FILE=<ログファイル>
API_SERVER_SEED_FILE=<seedファイルの場所>
```

------------------------------------------------------------------------
# API client

## client 生成
```
# 流れ: openapi.yaml -> openapi.json -> gen-client
./gen-client.sh
```

## client インストール・更新
```
npm install http://github.com/mndyu/localchat-api-client --save
npm update localchat-api-client
```

## client document 生成
```
// TODO
```

## client 例
```typescript
import { DefaultApi } from "localchat-api-client"

const api = new DefaultApi({}, "http://localhost:1324/api/v1", fetch)
api.usersGet().then(users => {
  console.log("users:", users)
})
```

## server 生成 (適当)
```
# openapi.yaml -> openapi.json -> gen-server
./gen-server.sh
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

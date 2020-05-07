# localchat-server


## 実行
```
# SQLite3
./dev-imdb.sh

# Postgres (途中)
./dev-postgres.sh
```

## client 生成
```
# openapi.yaml -> openapi.json -> gen-client
./gen-client.sh
```

## client インストール・更新
```
npm install http://github.com/mndyu/localchat-api-client --save
npm update localchat-api-client
```

## client document 生成
```
```

## client 例
```typescript
import { DefaultApi } from "localchat-api-client"

const api = new DefaultApi({}, "http://localhost:1324/api/v1", fetch)
api.usersGet().then(users => {
  console.log("users:", users)
})
```

## server 生成
```
# openapi.yaml -> openapi.json -> gen-server
./gen-server.sh
```


## TODO
- [ ] DB 接続
- [ ] スキーマ定義
- [ ] APIサーバ
- [ ] websocket
  - gorilla/websocket
  - REST との連携
  - 登録 & 受け取りの仕組み？
- [ ] Docker Hub & GitHub Actions


## curl
```
curl \
  -X GET \
  http://localhost:1324/api/v1/users

```

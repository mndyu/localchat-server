# API document


## ドキュメント作成
```
npm install create-openapi-repo
npx create-openapi-repo
```

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


#!/bin/bash


# gen (v2)
# yaml to json
# cat openapi.yaml | python3 -c 'import json, sys, yaml ; y=yaml.safe_load(sys.stdin.read()) ; json.dump(y, sys.stdout)' > openapi.json &&
# docker run --rm -v ${PWD}:/local swaggerapi/swagger-codegen-cli generate -c /local/gen-client-config.json -i /local/openapi.json -l typescript-fetch -o /local/gen-client &&

# gen (v3 to v2)
npm install api-spec-converter &&
npm audit fix &&
npx api-spec-converter --from=openapi_3 --to=swagger_2 --syntax=json openapi.yaml > openapi.json &&
docker run --rm -v ${PWD}:/local swaggerapi/swagger-codegen-cli generate -c /local/gen-client-config.json -i /local/openapi.json -l typescript-fetch -o /local/gen-client &&

# gen (v3)
# docker run --rm -v ${PWD}:/local swaggerapi/swagger-codegen-cli-v3 generate -c /local/gen-client-config.json -i /local/openapi.json -l javascript -o /local/gen-client

# config help
# docker run --rm -v ${PWD}:/local swaggerapi/swagger-codegen-cli-v3 config-help  -l typescript-fetch


# build
mkdir gen-client &&
cd gen-client &&
npm install &&
npm audit fix &&
npm run build &&

# git
git init &&
git remote add origin git@github.com:mndyu/localchat-api-client.git || true &&
git add . &&
git commit -m "update" &&
git push --set-upstream origin master &&

echo done

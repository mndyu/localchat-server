#!/bin/bash

# yaml to json
cat openapi.yaml | python3 -c 'import json, sys, yaml ; y=yaml.safe_load(sys.stdin.read()) ; json.dump(y, sys.stdout)' > openapi.json &&

# gen (v2)
docker run --rm -v ${PWD}:/local swaggerapi/swagger-codegen-cli generate -c /local/gen-client-config.json -i /local/openapi.json -l typescript-fetch -o /local/gen-client &&

# gen (v3)
# docker run --rm -v ${PWD}:/local swaggerapi/swagger-codegen-cli-v3 generate -i /local/openapi.json -l typescript-fetch -o /local/client

# config help
# docker run --rm -v ${PWD}:/local swaggerapi/swagger-codegen-cli config-help  -l typescript-fetch


# build
cd gen-client &&
npm install &&
npm audit fix &&
npm run build &&

# git
git init
git remote add origin git@github.com:mndyu/localchat-server.git
git add . &&
git commit -m "update" &&
git push --set-upstream origin master &&

echo done

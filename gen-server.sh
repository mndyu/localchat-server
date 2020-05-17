#!/bin/bash

# yaml to json
cat openapi.yaml | python3 -c 'import json, sys, yaml ; y=yaml.safe_load(sys.stdin.read()) ; json.dump(y, sys.stdout)' > openapi.json &&

# gen (v2)
# docker run --rm -v ${PWD}:/local swaggerapi/swagger-codegen-cli generate -c /local/gen-server-config.json -i /local/openapi.json -l nodejs-server -o /local/gen-server &&

# gen (v3 to v2)
npm install api-spec-converter &&
npx api-spec-converter --from=openapi_3 --to=swagger_2 --syntax=json openapi.yaml > openapi.json &&
docker run --rm -v ${PWD}:/local swaggerapi/swagger-codegen-cli generate -c /local/gen-server-config.json -i /local/openapi.json -l nodejs-server -o /local/gen-server &&

# gen (v3)
# docker run --rm -v ${PWD}:/local swaggerapi/swagger-codegen-cli-v3 generate -i /local/openapi.json -l typescript-fetch -o /local/client

# build
cd gen-server &&
npm install &&
npm audit fix &&

echo done

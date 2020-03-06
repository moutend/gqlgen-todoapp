#!/bin/bash

token="replace me!"

curl -i -X POST \
-H 'Content-Type: application/json' \
-H "Cookie: TOKEN=$token" \
--data '{"query":"query { tasks { id title content }}"}' \
localhost:8080/query

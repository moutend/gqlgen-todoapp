#!/bin/bash

token="replace me!"

curl -i -X POST \
-H 'Content-Type: application/json' \
--data "{\"query\":\"mutation { refreshToken(input: {token:\\\"$token\\\"}) }\"}" \
localhost:8080/query

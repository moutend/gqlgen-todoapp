#!/bin/bash

curl -i -X POST \
-H 'Content-Type: application/json' \
--data '{"query":"mutation { login(input: {name: \"your_name\", password: \"your_password\"})}"}' \
localhost:8080/query

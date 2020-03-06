#!/bin/bash

token="replace me!"
title=`date`
content=`uptime`

curl -i -X POST \
-H 'Content-Type: application/json' \
-H "Cookie: TOKEN=$token" \
--data "{\"query\":\"mutation { createTask(input: {title: \\\"$title\\\", content:\\\"$content\\\"}) { title content user { name }}}\"}" \
localhost:8080/query

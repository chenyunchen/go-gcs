#!/usr/bin/env bash

curl -X POST \
  http://localhost:7890/v1/signurl/ \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 789a1297-4ae1-493a-9cfa-af2b561f61e7' \
  -H 'cache-control: no-cache' \
  -d '{
	"fileName": "abc/123.txt",
	"contentType": "text/plain",
	"tag": "group",
	"payload": "{\r\n  \"from\": \"myawesomeId\",\r\n  \"groupId\": \"groupawesomeId\" \r\n}"
}'

#!/bin/sh

curl \
  --verbose \
  --request POST \
  'http://127.0.0.1:8080/api/v2/apps/ablog/events/mobile-app/9160' \
  --header 'Authorization: 1' \
  --header 'Content-Type: application/json' \
  --data '{}'

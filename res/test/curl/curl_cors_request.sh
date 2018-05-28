#!/bin/sh

curl \
  --verbose \
  --request OPTIONS \
  'http://127.0.0.1:8080/api/v2/apps/ablog/events/mobile-app/9160' \
  --header 'Origin: http://example.com' \
  --header 'Access-Control-Request-Headers: Origin, Accept, Content-Type' \
  --header 'Access-Control-Request-Method: POST' \

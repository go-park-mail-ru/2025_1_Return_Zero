#!/bin/bash

COOKIE_FILE=$(mktemp)

curl -X POST https://returnzero.ru/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email": "test@test.ru", "password": "test1"}' \
    -c $COOKIE_FILE -s

SESSION_ID=$(grep "session_id" $COOKIE_FILE | awk '{print $NF}')
echo "Found session_id: $SESSION_ID"

wrk -d 20m -c 100 -t 100 --header "Cookie: session_id=$SESSION_ID" https://returnzero.ru/api/v1/playlists/me -s fetch_db.lua

rm $COOKIE_FILE

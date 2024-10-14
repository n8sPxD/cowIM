#! /bin/bash

path="$(cd "$(dirname "\$0")" && pwd)"
execute="$path/../../microservices/msgToDB"

cd "$execute" || exit 1
go run .
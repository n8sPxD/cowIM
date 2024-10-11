#! /bin/bash

path="$(cd "$(dirname "\$0")" && pwd)"
execute="$path/../../microservices/user"

cd "$execute" || exit 1
go run .
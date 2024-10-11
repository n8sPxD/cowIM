#! /bin/bash

path="$(cd "$(dirname "\$0")" && pwd)"
execute="$path/../../microservices/msgForward"

cd "$execute" || exit 1
go run .
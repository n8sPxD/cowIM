#! /bin/bash

path="$(cd "$(dirname "\$0")" && pwd)"
execute="$path/../../gateway/http_gateway"

cd "$execute" || exit 1
go run .
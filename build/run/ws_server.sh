#! /bin/bash

path="$(cd "$(dirname "\$0")" && pwd)"
execute="$path/../../im-server"

cd "$execute" || exit 1
go run .

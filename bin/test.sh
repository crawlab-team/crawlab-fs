#!/bin/sh

export CODECOV_TOKEN=`cat /app/.codecov_token`
go test -race -coverprofile=coverage.txt -covermode=atomic
bash <(curl -s https://codecov.io/bash)

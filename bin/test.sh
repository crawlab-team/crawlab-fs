#!/bin/sh

go test -race -coverprofile=coverage.txt -covermode=atomic
bash <(curl -s https://codecov.io/bash)

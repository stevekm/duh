SHELL:=/bin/bash
.ONESHELL:

# go mod init sizemap
# $ go run main.go .
# gofmt -l -w .

format:
	gofmt -l -w main.go

test:
	set -euo pipefail
	go clean -testcache && \
	go test -v ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

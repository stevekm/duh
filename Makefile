SHELL:=/bin/bash
.ONESHELL:

# go mod init sizemap
# go get code.cloudfoundry.org/bytefmt
# go get github.com/google/go-cmp/cmp
# $ go run main.go .
# gofmt -l -w .

format:
	gofmt -l -w main.go

# go test -v ./...
test:
	set -euo pipefail
	go clean -testcache && \
	go test -v . | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

# need to check these cases
# go run main.go ./dir1
# go run main.go ./dir1/
# go run main.go dir1
# go run main.go dir1/
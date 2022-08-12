SHELL:=/bin/bash
.ONESHELL:

# go mod init sizemap
# go get code.cloudfoundry.org/bytefmt
# go get github.com/google/go-cmp/cmp
# $ go run main.go .
# gofmt -l -w .

format:
	gofmt -l -w main.go
	gofmt -l -w main_test.go

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
# go run main.go .
# go run main.go ./


run-all: $(DIR1)
	for i in . ./ ./dir1 ./dir1/ dir1 dir1/ dir1/go; do echo ">>> go run main.go $$i"; go run main.go $$i; done
	

SRC:=main.go
BIN:=duh
build:
	go build -o ./$(BIN) ./$(SRC)
.PHONY:build

# fatal: No names found, cannot describe anything.
GIT_TAG:=$(shell git describe --tags)
build-all:
	mkdir -p build ; \
	for os in darwin linux windows; do \
	for arch in amd64 arm64; do \
	output="build/$(BIN)-v$(GIT_TAG)-$$os-$$arch" ; \
	if [ "$${os}" == "windows" ]; then output="$${output}.exe"; fi ; \
	echo "building: $$output" ; \
	GOOS=$$os GOARCH=$$arch go build -o "$${output}" $(SRC) ; \
	done ; \
	done






# USAGE:
# $ make -f benchmarkdirs.makefile all

# ~~~~~ Set up Benchmark dir ~~~~~ #
# set up a dir with tons of files and some very large duplicate files to test the program against

# all: benchmark-dirs

# https://go.dev/dl/go1.18.3.darwin-amd64.tar.gz
# https://dl.google.com/go/go1.18.3.darwin-amd64.tar.gz

# BENCHDIR:=benchmarkdir
GO_TAR:=go1.18.3.darwin-amd64.tar.gz
$(GO_TAR):
	set -e
	wget https://dl.google.com/go/$(GO_TAR)

DIR1:=dir1
DIR2:=dir2
DIR3:=dir3

$(DIR1) $(DIR2) $(DIR3): $(GO_TAR)
	set -e
	mkdir -p "$(DIR1)"
	mkdir -p "$(DIR2)"
	mkdir -p "$(DIR3)"

benchmark-dirs: $(DIR1) $(DIR2) $(DIR3) $(GO_TAR)
	set -e
	tar -C "$(DIR1)" -xf "$(GO_TAR)"
	tar -C "$(DIR2)" -xf "$(GO_TAR)"
	tar -C "$(DIR3)" -xf "$(GO_TAR)"
	/bin/cp "$(GO_TAR)" $(DIR1)
	/bin/cp "$(GO_TAR)" $(DIR2)/go/
	/bin/cp "$(GO_TAR)" $(DIR2)/copy2.tar.gz
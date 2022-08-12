SHELL:=/bin/bash
.ONESHELL:

# USAGE:
# $ make -f benchmarkdirs.makefile all

# ~~~~~ Set up Benchmark dir ~~~~~ #
# set up a dir with tons of files and some very large duplicate files to test the program against

all: benchmark-dirs

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
	

# for i in $$(seq 1 5); do cat "$(GO_TAR)" >> $(BENCHDIR)/go/bin/foo ; done
# for i in $$(seq 1 5); do cat "$(GO_TAR)" >> $(BENCHDIR)/go/doc/foo2 ; done
# for i in $$(seq 1 10); do cat "$(GO_TAR)" >> $(BENCHDIR)/go/lib/bar ; done
# for i in $$(seq 1 10); do cat "$(GO_TAR)" >> $(BENCHDIR)/go/misc/bar2 ; done
# for i in $$(seq 1 15); do cat "$(GO_TAR)" >> $(BENCHDIR)/go/src/baz ; done
# for i in $$(seq 1 15); do cat "$(GO_TAR)" >> $(BENCHDIR)/go/test/baz2 ; done
# for i in $$(seq 1 20); do cat "$(GO_TAR)" >> $(BENCHDIR)/go/src/buzz ; done
# for i in $$(seq 1 20); do cat "$(GO_TAR)" >> $(BENCHDIR)/go/test/buzz2 ; done
# for i in $$(seq 1 20); do cat "$(GO_TAR)" >> $(BENCHDIR)/go/bin/buzz3 ; done

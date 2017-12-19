BUILD_NUMBER ?= 0

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

bin/tf-state-parser: $(SOURCES)
	go build -ldflags "-X main.build=${BUILD_NUMBER}" -o bin/tf-state-parser

run: bin/tf-state-parser
	bin/tf-state-parser -s testdata/state.json > import.sh
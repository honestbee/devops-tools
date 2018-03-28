BUILD_NUMBER ?= 0

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

# if .env exists, load it
-include .env
export

bin/prefix: $(SOURCES)
	go build -ldflags "-X main.build=${BUILD_NUMBER}" -o bin/prefix

run: bin/prefix
	bin/prefix

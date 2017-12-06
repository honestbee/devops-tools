
BUILD_NUMBER ?= 0

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

bin/github-tf: $(SOURCES)
	go build -ldflags "-X main.build=${BUILD_NUMBER}" -o bin/github-tf

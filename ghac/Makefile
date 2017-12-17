
BUILD_NUMBER ?= 0

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

bin/ghac: $(SOURCES)
	go build -ldflags "-X main.build=${BUILD_NUMBER}" -o bin/ghac

run: bin/ghac
	-rm output/teams-config/*.tf
	bin/ghac -s testdata/teams.yaml

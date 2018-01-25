
BUILD_NUMBER ?= 0

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

bin/drone-kfg: $(SOURCES)
	go build -ldflags "-X main.build=${BUILD_NUMBER}" -o bin/drone-kfg

run: bin/drone-kfg
	bin/drone-kfg -c STAGING_1A=internal.apse1a.staging -c STAGING_1B=internal.apse1b.staging honestbee/testRepo

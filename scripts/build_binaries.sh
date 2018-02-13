#!/bin/bash

set -ex
dep ensure
gox --osarch "darwin/amd64 linux/amd64" --output "../bin/${MBT_MODULE_NAME}-${MBT_MODULE_PROPERTY_VERSION}-{{.OS}}-{{.Arch}}"

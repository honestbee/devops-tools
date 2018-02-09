#!/bin/bash

set -ex
gox --osarch "darwin/amd64 linux/amd64" --output "../bin/${MBT_MODULE_NAME}_{{.OS}}_{{.Arch}}_${MBT_MODULE_PROPERTY_VERSION}"

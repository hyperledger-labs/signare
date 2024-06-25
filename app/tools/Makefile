SHELL = /bin/bash -o pipefail
WIRE_VERSION = "v0.6.0"

.PHONY: build
build:
	GOBIN=${PWD}/bin go install github.com/google/wire/cmd/wire@$(WIRE_VERSION)
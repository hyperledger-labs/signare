SHELL=/bin/bash -o pipefail

define goBuild
        $(eval artifactName = $1)
        $(eval ldflags = $2)
        $(eval os = $3)
        $(eval extension = $4)
        $(eval cgo = $5)
        $(eval osSuffix = _$(os))
        CGO_ENABLED=$(cgo) GOOS=$(os) GOARCH=amd64 go build -ldflags=$(ldflags) -o bin/$(artifactName)$(osSuffix)_amd64$(extension) ./cmd/$(artifactName)/$(artifactName).go
endef

.PHONY: build
build:
	$(eval branch = $(shell git rev-parse --abbrev-ref HEAD))
	$(eval commit = $(shell git rev-parse --short HEAD))
	$(eval build_time = $(shell date -u +%s))
	$(eval tag = $(shell git tag --points-at HEAD))
	$(eval ldflags = "-X main.commitHash=$(commit) -X main.buildTime=$(build_time) -X main.tag=$(tag)")
	$(call goBuild, signare, $(ldflags), linux,,1)

.PHONY: lint
lint:
	@echo "Executing linters"
	golangci-lint run

SHELL=/bin/bash -o pipefail

BASE_IMAGE ?= "scratch"
BUILD_BASE_IMAGE ?= "golang:1.22"

DOCKER_LOCAL_TAG = local
DOCKER_REPO ?= ""
ARTIFACT_NAME ?= "signare"
DOCKER_IMAGE_NAME ?= ${ARTIFACT_NAME}

.PHONY: lint
lint:
	@echo "Executing linters"
	@cd app; $(MAKE) lint
	@cd deployment; $(MAKE) lint

.PHONY: docker_build
docker_build:
	@echo "Building docker image"
	docker build \
		--build-arg BUILD_BASE_IMAGE=${BUILD_BASE_IMAGE} \
		--build-arg BASE_IMAGE=${BASE_IMAGE} \
		--build-arg GOPROXY=${GOPROXY} \
		-t ${DOCKER_IMAGE_NAME}:${DOCKER_LOCAL_TAG} \
		-f ./Dockerfile \
		./


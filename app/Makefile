SHELL=/bin/bash -o pipefail
GIT_ROOT_DIR=$(shell git rev-parse --show-toplevel)
API_LINTER=tools/spectral-openapi-linter
ARTIFACTS_DIR ?= _artifacts
CURDIR := $(shell pwd)
SIGNARE_DOCS="signare-docs"

.PHONY: test
test: lint unit_test

.PHONY: lint
lint:
	@echo "Executing linters"
	golangci-lint run

.PHONY: unit_test
unit_test:
	@echo "Creating artifacts dir"
	mkdir -p ${ARTIFACTS_DIR}
	@echo "Executing tests for application packages"
	gotestsum --jsonfile ${ARTIFACTS_DIR}/unit_tests_results.json --format short-verbose --junitfile ${ARTIFACTS_DIR}/junit.xml -- -p 1 -tags=unit -timeout=5m -coverpkg=github.com/hyperledger-labs/signare/app/pkg/... -coverprofile=${ARTIFACTS_DIR}/coverage.out github.com/hyperledger-labs/signare/app/pkg/...

.PHONY: integration_test
integration_test:
	@echo "TODO"

.PHONY: compatibility_matrix
compatibility_matrix:
	@echo "TODO"

.PHONY: tools.generate
tools.generate:
	@echo "Generating the API spec bundle"
	swagger-cli bundle -t yaml docs/open-api-specs/signare-api/index.yaml -o docs/open-api-specs/signare-api/index_bundled_openapi.yaml
	@echo "Executing linters on API spec"
	@spectral lint ../../docs/open-api-specs/signare-api/index_bundled_openapi.yaml --ruleset ${GIT_ROOT_DIR}/app/spectral-openapi-linter-ruleset.yaml --format pretty --fail-severity=error --verbose
	@echo "Generating HTTP infra"
	cd tools/http-infra-openapi-generator; $(MAKE)
	@echo "Executing go generate"
	@PATH=${PWD}/tools/bin:$$PATH go generate ./...
	@echo "Formatting go code"
	@gofmt -w -s . && goimports -w .

.PHONY: tools.serve_docs
tools.serve_docs:
	@ docker rm -f $(SIGNARE_DOCS) || true
	@cd docs && docker build -t $(SIGNARE_DOCS) .
	@docker run -d --rm --name $(SIGNARE_DOCS) -p "8005:8000" -v $(CURDIR)/docs/mkdocs/docs:/app/docs $(SIGNARE_DOCS)
	@echo 'Hosting documentation on http://localhost:8005, refresh if needed until the documentation is build'
	@xdg-open http://localhost:8005

.PHONY: tools.close_docs
tools.close_docs:
	@docker rm -f $(SIGNARE_DOCS) || true
	@echo 'Closing & deleting documentation container'

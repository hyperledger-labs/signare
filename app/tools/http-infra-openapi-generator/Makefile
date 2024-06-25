SHELL=/bin/bash -o pipefail

GIT_ROOT_DIR:=$(shell git rev-parse --show-toplevel)
TMPDIR := $(shell mktemp -d /tmp/signare.tools.http-infra-openapi-generator.XXXX)

OPENAPI_GENERATOR_POM_VERSION := $(shell cd $(GIT_ROOT_DIR)/app/tools/http-infra-openapi-generator/signare-plugin && mvn help:evaluate -Dexpression=openapi-generator-version -q -DforceStdout)
OPENAPI_GENERATOR_JAR_VERSION := $(shell java -cp $(GIT_ROOT_DIR)/app/tools/http-infra-openapi-generator/openapi-generator-cli.jar org.openapitools.codegen.OpenAPIGenerator version)

API_SPEC_ROOT ?= "$(GIT_ROOT_DIR)/app/docs/open-api-specs/signare-api"
API_SPEC_BUNDLE ?= "$(GIT_ROOT_DIR)/app/docs/open-api-specs/signare-api/index_bundled_openapi.yaml"
API_INFRA_GENERATED_OUT ?= "$(GIT_ROOT_DIR)/app/pkg/infra/generated"

.PHONY: default
default: build_signare_plugin generate_signare_http_infra

.PHONY: build_signare_plugin
build_signare_plugin:
	@cd $(GIT_ROOT_DIR)/app/tools/http-infra-openapi-generator/signare-plugin; mvn package

.PHONY: generate_signare_http_infra
generate_signare_http_infra: generate_http_infra format_go_code format_files

.PHONY: generate_http_infra
generate_http_infra: bundle build_signare_plugin check_signare_plugin_openapi_pom_jar_versions
	@swagger-cli bundle -t yaml $(API_SPEC_BUNDLE) -o $(TMPDIR)/bundle.yaml
	- @rm -rf $(API_INFRA_GENERATED_OUT)
	@GO_POST_PROCESS_FILE="/usr/local/go/bin/gofmt -w" java -cp $(GIT_ROOT_DIR)/app/tools/http-infra-openapi-generator/openapi-generator-cli.jar:$(GIT_ROOT_DIR)/app/tools/http-infra-openapi-generator/signare-plugin/target/signare-plugin-1.0.0.jar org.openapitools.codegen.OpenAPIGenerator generate -g signare-http-infra -i $(TMPDIR)/bundle.yaml --enable-post-process-file --package-name httpinfra --model-package httpinfra --api-package httpinfra --additional-properties=useOneOfDiscriminatorLookup=true,generatedActionsFilePath=$(GIT_ROOT_DIR)/app/include/rbac/actions-generated.yaml -o $(API_INFRA_GENERATED_OUT)

.PHONY: openapi_generator_jar_version
openapi_generator_jar_version:
	@echo $(OPENAPI_GENERATOR_JAR_VERSION)

.PHONY: openapi_generator_version
openapi_generator_version:
	@echo $(OPENAPI_GENERATOR_POM_VERSION)

.PHONY: bundle
bundle:
	@echo "Generating the API spec bundle"
	swagger-cli bundle -t yaml $(API_SPEC_ROOT)/index.yaml -o $(API_SPEC_ROOT)/index_bundled_openapi.yaml

check_signare_plugin_openapi_pom_jar_versions:
ifneq ($(OPENAPI_GENERATOR_POM_VERSION), $(OPENAPI_GENERATOR_JAR_VERSION))
	$(error signare openapi generator version does not match the jar one)
endif

format_go_code:
	gofmt -w -s $(GIT_ROOT_DIR)/app && goimports -w $(GIT_ROOT_DIR)/app

# format_files is needed in order to align the generated code with the current lint configuration
format_files:
	find $(GIT_ROOT_DIR)/app/pkg/infra -type f -name "*.go" -exec sed -i 's/Ip/IP/g' {} \;
	find $(GIT_ROOT_DIR)/app/pkg/infra -type f -name "*.go" -exec sed -i 's/Tls/TLS/g' {} \;
	find $(GIT_ROOT_DIR)/app/pkg/infra -type f -name "*.go" -exec sed -i 's/\([^"]\)Url\([^"]\)/\1URL\2/g' {} \;
	find $(GIT_ROOT_DIR)/app/pkg/infra -type f -name "*.go" -exec sed -i 's/Api/API/g' {} \;
	find $(GIT_ROOT_DIR)/app/pkg/infra -type f -name "*.go" -exec sed -i 's/Http/HTTP/g' {} \;

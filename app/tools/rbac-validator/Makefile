SHELL=/bin/bash -o pipefail

SIGNARE_DIR=$(shell cd ../..; pwd)

.PHONY: tools.run_default
tools.run_default:
	@go run ./cmd/rbac_validator.go \
	  --openapiSpecFiles $(SIGNARE_DIR)/docs/open-api-specs/signare-api/index_bundled_openapi.yaml \
      --rolesFilePath $(SIGNARE_DIR)/include/rbac/roles.yaml \
      --permissionsFilePath $(SIGNARE_DIR)/include/rbac/permissions.yaml \
      --actionsFilesPath $(SIGNARE_DIR)/include/rbac/actions-generated.yaml,$(SIGNARE_DIR)/include/rbac/actions-manual.yaml \
      --operationIdInclusions rpc.method.eth_generateAccount,rpc.method.eth_removeAccount,rpc.method.eth_accounts,rpc.method.eth_signTransaction

.PHONY: tools.help
tools.help:
	@go run ./cmd/rbac_validator.go --help

#!/usr/bin/env bash

set -e
GIT_ROOT_DIR=$(git rev-parse --show-toplevel)

go run "${GIT_ROOT_DIR}"/deployment/cmd/signare/signare.go upgrade --config "${GIT_ROOT_DIR}"/deployment/examples/config

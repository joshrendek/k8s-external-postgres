#!/bin/bash
SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${SCRIPT_ROOT}; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}
${CODEGEN_PKG}/generate-groups.sh all \
github.com/joshrendek/k8s-external-postgres/pkg/client \
github.com/joshrendek/k8s-external-postgres/pkg/apis \
postgresql:v1

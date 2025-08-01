#!/usr/bin/env bash

echo "[golangci-lint] checking go files"

golangci-lint run ./...

# shellcheck disable=SC2181
if [ $? -ne 0 ]; then
  echo "[golangci-lint] found issues that needs to be fixed"
  exit 1
fi



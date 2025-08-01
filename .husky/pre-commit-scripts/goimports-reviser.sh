#!/usr/bin/env bash

files=("$@")
if [ ${#files[@]} -eq 0 ]; then
  echo "[goimports-reviser] no go files to check"
  exit 0
fi

echo "[goimports-reviser] formatting ${#files[@]} go files"

goimports-reviser -format ${files[@]} &> /dev/null

# shellcheck disable=SC2181
if [ $? -ne 0 ]; then
  echo "[goimports-reviser] found issues formatting go files"
  exit 1
fi

git add --all
exit 0

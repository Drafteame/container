#!/usr/bin/env bash

echo "[commitizen] checking commit message with commitizen"

cz check --commit-msg-file .git/COMMIT_EDITMSG

# shellcheck disable=SC2181
if [ $? -ne 0 ]; then
  echo "[commitizen] found issues in commit message"
  exit 1
fi

exit 0
#!/usr/bin/env bash

echo "[go-mod-tidy] checking go.mod files"

function get_base_dir() {
  local dir

  local file="$1"
  dir=$(dirname "$file")

  echo "$dir"
}

# Find and process go.mod files
go_mod_files=$(fd 'go.mod' --glob)
go_sum_files=$(fd 'go.sum' --glob)

if [ -n "$go_mod_files" ]; then
  for file in $go_mod_files; do
    echo "[go-mod-tidy] tidying $file"

    cd "$(get_base_dir "$file")" && go mod tidy -v

    # shellcheck disable=SC2181
    if [ $? -ne 0 ]; then
      exit 2
    fi

    # shellcheck disable=SC2164
    cd - > /dev/null
  done
fi

# add files to git
if [ -n "$go_mod_files" ]; then
  git add $go_mod_files
fi

if [ -n "$go_sum_files" ]; then
  git add $go_sum_files
fi

git commit -m "deps: tidy go.mod files"
exit 0

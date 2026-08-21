#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
mapfile -t files < <(
  find "$root" \
    -path "$root/.git" -prune -o \
    -path "$root/.gocache" -prune -o \
    -path "$root/.gomodcache" -prune -o \
    -path "$root/frontend" -prune -o \
    -path "$root/migrations" -prune -o \
    -path "$root/docs" -prune -o \
    -path "$root/vendor" -prune -o \
    -name '*_test.go' -prune -o \
    -name '*.go' -type f -print | sort
)

total=0
printf 'FILES\n'
for file in "${files[@]}"; do
  rel="${file#${root}/}"
  lines=$(wc -l < "$file")
  printf '%s %s\n' "$lines" "$rel"
  total=$((total + lines))
done
printf 'FILE_COUNT %s\n' "${#files[@]}"
printf 'LINE_COUNT %s\n' "$total"

#!/usr/bin/env bash
set -euo pipefail

version="${VERSION:-dev}"
out_dir="${OUT_DIR:-dist}"
binary_name="${BINARY_NAME:-seekai}"

mkdir -p "$out_dir"

go build \
  -trimpath \
  -ldflags "-s -w -X main.version=${version}" \
  -o "${out_dir}/${binary_name}" \
  .

echo "built ${out_dir}/${binary_name} (${version})"

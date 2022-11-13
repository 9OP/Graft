#!/usr/bin/env bash
set -e -o pipefail

if ! command -v structslop &> /dev/null ; then
    echo "structslop not installed or available in the PATH" >&2
    exit 1
fi

pkg=$(go list)
for dir in $(echo $@|xargs -n1 dirname|sort -u); do
  structslop $pkg/$dir
done

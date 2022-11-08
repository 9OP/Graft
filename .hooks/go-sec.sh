#!/usr/bin/env bash
set -e -o pipefail

if ! command -v gosec &> /dev/null ; then
    echo "gosec not installed or available in the PATH" >&2
    exit 1
fi

pkg=$(go list)
for dir in $(echo $@|xargs -n1 dirname|sort -u); do
  gosec $pkg/$dir
done

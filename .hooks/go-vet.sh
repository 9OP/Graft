#!/usr/bin/env bash
set -e

pkg=$(go list)
for dir in $(echo $@|xargs -n1 dirname|sort -u); do
  # echo "echo path $pkg/$dir"
  go vet $pkg/$dir
done

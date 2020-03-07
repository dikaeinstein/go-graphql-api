#!/bin/bash

set -eo pipefail

#All Packages
PACKAGES=$(find . -name '*.go' -print0 | xargs -0 -n1 dirname | sort --unique)

echo 'packages needing linting'
for pkg in ${PACKAGES[@]}; do
    echo $pkg $(golint $pkg | wc -l) | grep -v -e '\s0$' -e '^./vendor' || test $? == 1
done

echo ''
echo 'packages passing lint'
for pkg in ${PACKAGES[@]}; do
    echo $pkg $(golint $pkg | wc -l) | grep '\s0$' | grep -v '^./vendor' || test $? == 1
done

echo ''
echo 'running `go vet`'
go vet ./...

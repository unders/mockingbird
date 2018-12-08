#!/usr/bin/env bash

set -e

main() {
    echo "Installs the build tool mage"
    rm -rf /tmp/mage
    cd /tmp
    git clone git@github.com:magefile/mage
    cd mage
    go run bootstrap.go
    mage --version
}

main

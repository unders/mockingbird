#!/usr/bin/env bash

set -e

##
## Git and Go must be installed
##

main() {
    echo "bootstraping github.com/unders/mockingbird/tool"
    go version
    go mod download

    ##
    ## Install mage command tool
    ##
    rm -rf /tmp/mage
    cd /tmp
    git clone git@github.com:magefile/mage
    cd mage
    go run bootstrap.go
    mage --version
}

main

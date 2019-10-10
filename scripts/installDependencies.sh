#!/bin/bash

DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/common.sh"

#Go dependencies installation; uses 'go mod'
shw_norm "Installing project dependencies"
export GO111MODULE=on #necessary to install dependencies inside 'GOPATH/src'
go mod vendor
shw_info "Project dependencies installed"
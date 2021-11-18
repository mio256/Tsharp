#!/bin/sh

go build -o tsh main.go

if [[ "$OSTYPE" != ^msys ]]; then
    sudo cp ./tsh /usr/local/bin/tsh
fi

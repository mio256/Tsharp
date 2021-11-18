#!/bin/sh

go build -o tsh main.go

if [[ "$OSTYPE" =~ ^mac ]]; then
    sudo cp ./tsh /usr/local/bin/tsh
fi

if [[ "$OSTYPE" =~ ^linux ]]; then
    sudo cp ./tsh /usr/local/bin/tsh
fi

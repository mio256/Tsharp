#!/bin/sh

go build -o tsh main.go
sudo cp ./tsh /usr/local/bin/tsh
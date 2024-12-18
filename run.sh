#! /bin/bash

SNIPPETBOX_ADDR=":4050"

clear
go run ./cmd/web -addr=$SNIPPETBOX_ADDR -debug
#go run ./cmd/web -addr=$SNIPPETBOX_ADDR >>/tmp/snippetbox_info.log 2>> /tmp/snippetbox_error.log

#! /bin/bash

OUT="/tmp/profile.out"
DIR="./..."

go test -covermode=count -coverprofile=$OUT $DIR
go tool cover -html=$OUT

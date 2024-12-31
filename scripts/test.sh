#! /bin/bash
clear

DIR="./internal/models"
OPC=""

#PORT=4050
#URL="http://localhost:"$PORT

#curl -i -X PUT $URL
#curl -i -X PUT $URL"/snippet/view?id=12"
#curl -i $URL"/snippet/view?id=1"

#go clean -testcache

echo "COVERAGE"
go test -cover $DIR $OPC
#go test -coverprofile=/tmp/profile.out ./...

echo ""
echo "TESTING"
go test $DIR $OPC

echo ""
echo "VERBOSE"
go test -v $DIR $OPC
#go test -v -run="^TestReadableDate/^UTC$" ./cmd/web

#go test -failfast ./cmd/web
#go test -parallel 4 ./...
#go test -race ./...

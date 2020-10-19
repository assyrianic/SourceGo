#!/bin/bash
cd "$(dirname "$0")"

go build go2sp.go
#GOOS=windows GOARCH=386 go build go2sp.go

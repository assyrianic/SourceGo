#!/bin/bash
cd "$(dirname "$0")"

## windows bit
#GOOS=windows GOARCH=amd64 go build -o go2sp_x64.exe   go2sp.go
#GOOS=windows GOARCH=386   go build -o go2sp_ia32.exe  go2sp.go

## linux bit
#GOOS=linux   GOARCH=amd64 go build -o go2sp_x64       go2sp.go
#GOOS=linux   GOARCH=386   go build -o go2sp_ia32      go2sp.go

## mac os bit
#GOOS=darwin  GOARCH=amd64 go build -o go2sp_mac_x64   go2sp.go
#GOOS=darwin  GOARCH=arm64 go build -o go2sp_mac_arm64 go2sp.go
go build go2sp.go

#!/bin/bash -e

if [[ $1 = "-loc" ]]; then
    find . -name '*.go' | xargs wc -l | sort -n
    exit
fi

VER=0.1.0beta
ID=$(git rev-parse HEAD | cut -c1-7)

cd daemon/mysql-cluster

if [[ $1 = "-linux" ]]; then
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/nicholaskh/golib/server.VERSION $VER -X github.com/nicholaskh/golib/server.BuildID $ID -w"
    mv mysql-cluster ../../bin/mysql-cluster.linux
else
    go build -ldflags "-X github.com/nicholaskh/golib/server.VERSION $VER -X github.com/nicholaskh/golib/server.BuildID $ID -w"
    #go build -race -v -ldflags "-X github.com/nicholaskh/golib/server.BuildID $ID -w"
    mv mysql-cluster ../../bin/mysql-cluster.mac
    ../../bin/mysql-cluster.mac -v
fi

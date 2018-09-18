#!/bin/bash
export GOOS=linux
export GOARCH=386
export CGO_ENABLED=0

cd accountservice;go get;go build -o accountservice-linux-amd64;echo built `pwd`;cd ..

docker build -t dasetova/accountservice accountservice/

docker service rm accountservice
docker service create --name=accountservice --replicas=1 --network=my_network -p=6767:6767 dasetova/accountservice
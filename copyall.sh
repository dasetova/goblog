#!/bin/bash
export GOOS=linux
export GOARCH=386
export CGO_ENABLED=0

cd accountservice;go get;go build -o accountservice-linux-amd64;echo built `pwd`;cd ..
cd healthchecker;go get;go build -o healthchecker-linux-amd64;echo built `pwd`;cd ..

cp healthchecker/healthchecker-linux-amd64 accountservice/

sudo docker build -t dasetova/accountservice accountservice/

sudo docker service rm accountservice
sudo docker service create --name=accountservice --replicas=1 --network=my_network -p=6767:6767 dasetova/accountservice
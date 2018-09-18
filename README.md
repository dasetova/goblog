# goblog
Learning repository, following tutorial in https://dzone.com/articles/go-microservices-blog-series-part-1

## Creating docker 
+ env GOOS=linux GOARCH=386 go build -o accountservice-linux-amd64
+ sudo docker build -t dasetova/accountservice accountservice/
+ sudo docker run --rm dasetova/accountservice

## Docker swarm
docker-machine create --driver virtualbox manager1
docker-machine env manager1
docker $(docker-machine config manager1) swarm init --advertise-addr $(docker-machine ip manager1)
sudo docker network create --driver overlay my_network
sudo docker service create --name=accountservice --replicas=1 --network=my_network -p=6767:6767 dasetova/accountservice

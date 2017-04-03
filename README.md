# Syros

A highly available and horizontally scalable DevOps tool for managing microservices across multiple regions and environments. 

Components:

* Syros Agent (collects various system information)
* Syros Indexer (aggregates, transforms and persists collected data)
* Syros App (management UI and API)

Backend:

* NATS (communication backbone)
* MongoDB (persistence layer)

HA Setup:

* Agent: 2 instances per environment or one per host, indexer will do deduplication
* Indexer: 2 instances per environment, NATS will load balance the messages between instances
* App: one per environment, HAProxy or NGNIX can be used but not required
* NATS: 3 instances minimum 
* MongoDB: 3 instances minimum 

### Development 

Syros back-end is written in golang and the front-end in ES6 javascript (VueJs).

Prerequisites:

* macOS or Linux
* golang >= 1.7
* node >= 4.0
* npm >= 3.0
* docker >= 1.13
* make >= 3.81

Local setup:

```sh
# clone the repo into your go PATH under github.com/stefanprodan
$ git clone https://github.com/stefanprodan/syros.git
$ cd syros
# install go dependencies
$ go get -u github.com/kardianos/govendor
$ govendor sync
# install node dependencies
$ cd ui
$ npm install
```

Run locally:

```sh
# start NATS and MongoDB
$ docker-compose up -d
# build and run all services
$ make build run APP_VERSION=0.0.1 MONGO=192.168.1.135:27017 NATS=nats://192.168.1.135:4222
# remove build artifacs 
$ make clean
# remove containers and images
$ make purge APP_VERSION=0.0.1
# run go fmt and go vet
$ make fmt vet
```

### Continuous Integration

The CI pipeline is written in Make and uses Docker containers, 
no external dependencies like go or nodejs are required to build, test and deploy the services.

```sh
# build the UI with webpack and the golang binaries for Alpine
$ make build APP_VERSION=0.0.1
# run integration tests
$ make build test APP_VERSION=0.0.1 MONGO=192.168.1.135:27017 NATS=nats://192.168.1.135:4222
# push Docker images to registry
$ make build pack push APP_VERSION=0.0.1 REGISTRY=index.docker.io REPOSITORY=stefanprodan
# remove test containers and local images
$ make purge APP_VERSION=0.0.1
```

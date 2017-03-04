# Syros

A highly available and horizontally scalable DevOps tool for managing microservices across multiple regions and environments. 

Components:

* Syros Agent (collects various system data)
* Syros Indexer (aggregates, transforms and persists collected data)
* Syros App (management UI and API)

Backend:

* NATS (communication backbone)
* RethinkDB (persistence layer)

### Development 

Syros back-end is written in golang and the front-end in javascript (VueJs).

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
# run NATS and RethinkDB localy 
$ docker-compose up -d
# install node dependencies
$ cd ui
$ npm install
```

Build, pack, test and deploy:

The build system is done with Make and uses Docker containers. 

```sh
# build the UI with webpack and the golang binaries for Alpine
$ make build APP_VERSION=0.0.1
# run services on local containers
$ make run APP_VERSION=0.0.1
# run integration tests localy
$ make run test APP_VERSION=0.0.1 RDB=192.168.1.135:28015 NATS=nats://192.168.1.135:4222
# push Docker images to registry
$ make build pack push APP_VERSION=0.0.1 REGISTRY=index.docker.io REPOSITORY=stefanprodan
# remove containers, images and build artifacs 
$ make purge APP_VERSION=0.0.1
# run go fmt and go vet
$ make fmt vet
```



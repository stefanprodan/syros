#FROM golang:1.8.3
FROM golang:1.8.3-alpine

# services
ADD /models /go/src/github.com/stefanprodan/syros/models
ADD /agent /go/src/github.com/stefanprodan/syros/agent
ADD /indexer /go/src/github.com/stefanprodan/syros/indexer
ADD /api /go/src/github.com/stefanprodan/syros/api

# deps
ADD Gopkg.toml /go/src/github.com/stefanprodan/syros/Gopkg.toml
ADD Gopkg.lock /go/src/github.com/stefanprodan/syros/Gopkg.lock

# solution root
WORKDIR /go/src/github.com/stefanprodan/syros

# pull deps
RUN apk add --no-cache --virtual git
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure

# output
RUN mkdir /go/dist
VOLUME /go/dist

#FROM golang:1.8.3
FROM golang:1.8.3-alpine

# services
ADD /models /go/src/github.com/stefanprodan/syros/models
ADD /agent /go/src/github.com/stefanprodan/syros/agent
ADD /indexer /go/src/github.com/stefanprodan/syros/indexer
ADD /api /go/src/github.com/stefanprodan/syros/api

# deps
ADD /vendor/vendor.json /go/src/github.com/stefanprodan/syros/vendor/vendor.json

# solution root
WORKDIR /go/src/github.com/stefanprodan/syros

# pull deps
RUN apk add --no-cache --virtual git
RUN go get -u github.com/kardianos/govendor
RUN govendor sync

# output
RUN mkdir /go/dist
VOLUME /go/dist

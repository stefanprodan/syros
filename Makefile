SHELL:=/bin/bash

APP_VERSION?=0.1.0

# build vars
DIST:=$$(pwd)/dist
BUILD_DATE:=$(shell date -u +%Y-%m-%d_%H.%M.%S)
GIT_REPOSITORY:=github.com/stefanprodan/syros
GIT_COMMIT:=$(shell git rev-parse HEAD)
GIT_BRANCH:=$(shell git symbolic-ref --short HEAD)
MAINTAINER:="Stefan Prodan"

# go tools
PACKAGES:=$(shell go list ./... | grep -v '/vendor/')
VETARGS:=-asmdecl -atomic -bool -buildtags -copylocks -methods -nilfunc -rangeloops -shift -structtags -unsafeptr

# run vars
MONGO?=192.168.1.135:27017
NATS?=nats://192.168.1.135:4222

#deploy vars
REGISTRY?=index.docker.io
REPOSITORY?=stefanprodan

TIME_START:=$(shell date +%s)
define DURATION
@time_end=`date +%s` ; time_exec=`awk -v "TS=${TIME_START}" -v "TE=$$time_end" 'BEGIN{TD=TE-TS;printf "%02dm:%02ds\n",TD/(60)%60,TD%60}'` ; echo "$@ duration $${time_exec} "
endef

build: clean
	@echo ">>> Building syros-ui-build image"
	@docker build -t syros-ui-build:$(BUILD_DATE) -f build.node.dockerfile .

	@echo ">>> Building syros-ui"
	@docker run --rm  -v "$(DIST)/ui:/usr/src/app/dist" syros-ui-build:$(BUILD_DATE) \
		bash -c "npm run build"
	@docker rmi syros-ui-build:$(BUILD_DATE)

	@echo ">>> Building syros-services-build image"
	@docker build -t syros-services-build:$(BUILD_DATE) -f build.golang.dockerfile .

	@echo ">>> building syros-agent"
	@docker run --rm  -v "$(DIST):/go/dist" syros-services-build:$(BUILD_DATE) \
		go build -o /go/dist/agent github.com/stefanprodan/syros/agent

	@echo ">>> Building syros-indexer"
	@docker run --rm  -v "$(DIST):/go/dist" syros-services-build:$(BUILD_DATE) \
		go build -o /go/dist/indexer github.com/stefanprodan/syros/indexer

	@echo ">>> Building syros-api"
	@docker run --rm  -v "$(DIST):/go/dist" syros-services-build:$(BUILD_DATE) \
		go build -o /go/dist/api github.com/stefanprodan/syros/api

	@docker rmi syros-services-build:$(BUILD_DATE)

	@echo ">>> Build artifacts:"
	@find dist -type f -print0 | xargs -0 ls -t
	$(DURATION)

pack:
	@echo ">>> Building syros-app image for deploy"
	@docker build -t syros-app:$(APP_VERSION) \
	    --build-arg APP_VERSION=$(APP_VERSION) \
	    --build-arg BUILD_DATE=$(BUILD_DATE) \
	    --build-arg GIT_REPOSITORY=$(GIT_REPOSITORY) \
	    --build-arg GIT_BRANCH=$(GIT_BRANCH) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
        --build-arg MAINTAINER=$(MAINTAINER) \
		-f deploy.app.dockerfile .

	@echo ">>> Building syros-indexer image for deploy"
	@docker build -t syros-indexer:$(APP_VERSION) \
	    --build-arg APP_VERSION=$(APP_VERSION) \
	    --build-arg BUILD_DATE=$(BUILD_DATE) \
	    --build-arg GIT_REPOSITORY=$(GIT_REPOSITORY) \
	    --build-arg GIT_BRANCH=$(GIT_BRANCH) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
        --build-arg MAINTAINER=$(MAINTAINER) \
		-f deploy.indexer.dockerfile .

	@echo ">>> Building syros-agent image for deploy"
	@docker build -t syros-agent:$(APP_VERSION) \
	    --build-arg APP_VERSION=$(APP_VERSION) \
	    --build-arg BUILD_DATE=$(BUILD_DATE) \
	    --build-arg GIT_REPOSITORY=$(GIT_REPOSITORY) \
	    --build-arg GIT_BRANCH=$(GIT_BRANCH) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
        --build-arg MAINTAINER=$(MAINTAINER) \
		-f deploy.agent.dockerfile .

	@echo ">>> Images ready for deploy:"
	@docker images | grep syros
	$(DURATION)

run: pack
	@echo ">>> Starting syros-app container"
	@docker run -dp 8888:8888 --name syros-app-$(APP_VERSION) \
	    --restart unless-stopped \
		syros-app:$(APP_VERSION) \
		-MongoDB=$(MONGO) \
		-LogLevel=info

	@echo ">>> Starting syros-indexer container"
	@docker run -dp 8887:8887 --name syros-indexer-$(APP_VERSION) \
	    --restart unless-stopped \
		syros-indexer:$(APP_VERSION) \
		-MongoDB=$(MONGO) \
		-DatabaseStale=0 \
		-Nats=$(NATS) \
		-LogLevel=info

	@echo ">>> Starting syros-agent container"
	@docker run -dp 8886:8886 --name syros-agent-$(APP_VERSION) \
	    --restart unless-stopped \
	    -v /var/run/docker.sock:/var/run/docker.sock:ro \
	    syros-agent:$(APP_VERSION) \
		-DockerApiAddresses=unix:///var/run/docker.sock \
		-Environment=dev \
		-Nats=$(NATS) \
		-LogLevel=info

	@echo ">>> syros-app logs:"
	@docker logs syros-app-$(APP_VERSION)
	@echo ">>> syros-indexer logs:"
	@docker logs syros-indexer-$(APP_VERSION)
	@echo ">>> syros-agent logs:"
	@docker logs syros-agent-$(APP_VERSION)
	$(DURATION)

test: run
	@echo ">>> Checking syros-app status"
	@curl --fail http://localhost:8888/status

	@echo ">>> Checking syros-app auth endpoint"
	$(eval TOKEN := $(shell curl --fail -X POST http://localhost:8888/api/auth/login \
	  -d '{"name":"admin","password":"admin"}' \
	  -H "Content-type: application/json"))
	@echo ">>> JWT token acquired: $(TOKEN)"

	@echo ">>> Checking syros-app hosts endpoint"
	@curl --fail http://localhost:8888/api/docker/hosts \
	  -H "Authorization: Bearer $(TOKEN)" \
	  -H "Content-type: application/json"

	@echo ">>> Checking syros-indexer status"
	@curl --fail http://localhost:8887/status

	@echo ">>> Checking syros-agent status"
	@curl --fail http://localhost:8886/status
	$(DURATION)

push:
	@echo ">>> Pushing syros-app to $(REGISTRY)/$(REPOSITORY)"
	@docker tag syros-app:$(APP_VERSION) $(REGISTRY)/$(REPOSITORY)/syros-app:$(APP_VERSION)
	@docker tag syros-app:$(APP_VERSION) $(REGISTRY)/$(REPOSITORY)/syros-app:latest
	@docker push $(REGISTRY)/$(REPOSITORY)/syros-app:$(APP_VERSION)
	@docker push $(REGISTRY)/$(REPOSITORY)/syros-app:latest

	@echo ">>> Pushing syros-indexer to $(REGISTRY)/$(REPOSITORY)"
	@docker tag syros-indexer:$(APP_VERSION) $(REGISTRY)/$(REPOSITORY)/syros-indexer:$(APP_VERSION)
	@docker tag syros-indexer:$(APP_VERSION) $(REGISTRY)/$(REPOSITORY)/syros-indexer:latest
	@docker push $(REGISTRY)/$(REPOSITORY)/syros-indexer:$(APP_VERSION)
	@docker push $(REGISTRY)/$(REPOSITORY)/syros-indexer:latest

	@echo ">>> Pushing syros-agent to $(REGISTRY)/$(REPOSITORY)"
	@docker tag syros-agent:$(APP_VERSION) $(REGISTRY)/$(REPOSITORY)/syros-agent:$(APP_VERSION)
	@docker tag syros-agent:$(APP_VERSION) $(REGISTRY)/$(REPOSITORY)/syros-agent:latest
	@docker push $(REGISTRY)/$(REPOSITORY)/syros-agent:$(APP_VERSION)
	@docker push $(REGISTRY)/$(REPOSITORY)/syros-agent:latest
	$(DURATION)

fmt:
	@echo ">>> Running go fmt $(PACKAGES)"
	@go fmt $(PACKAGES)
	$(DURATION)

vet:
	@echo ">>> Running go vet $(VETARGS)"
	@go list ./... \
		| grep -v /vendor/ \
		| cut -d '/' -f 4- \
		| xargs -n1 \
			go tool vet $(VETARGS) ;\
	if [ $$? -ne 0 ]; then \
		echo ""; \
		echo "go vet failed"; \
	fi
	$(DURATION)

clean:
	@if [ -d "$(DIST)" ]; then \
		echo "output directory found at $(DIST) removing content"; \
		rm -rf $(DIST); \
	fi
	$(DURATION)

purge:
	@docker rm -f syros-app-$(APP_VERSION) syros-agent-$(APP_VERSION) syros-indexer-$(APP_VERSION) || true
	@docker rmi $$(docker images | awk '$$1 ~ /syros/ { print $$3 }') || true
	$(DURATION)

.PHONY: build

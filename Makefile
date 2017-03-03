SHELL:=/bin/bash

APP_VERSION?="0.0.1"

# build vars
DIST:=$$(pwd)/dist
BUILD_DATE:=$(shell date -u +%Y-%m-%d_%H.%M.%S)
GIT_COMMIT:=$(shell git rev-parse HEAD)
GIT_BRANCH:=$(shell git symbolic-ref --short HEAD)
PACKAGES:=$(shell go list ./... | grep -v '/vendor/')
VETARGS:=-asmdecl -atomic -bool -buildtags -copylocks -methods -nilfunc -rangeloops -shift -structtags -unsafeptr

# run vars
RethinkDB?=192.168.1.135:28015
Nats?=nats://192.168.1.135:4222

TIME_START:=$(shell date +%s)
define DURATION
@time_end=`date +%s` ; time_exec=`awk -v "TS=${TIME_START}" -v "TE=$$time_end" 'BEGIN{TD=TE-TS;printf "%02dm:%02ds\n",TD/(60)%60,TD%60}'` ; echo "$@ duration $${time_exec} "
endef

build: purge
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

	@echo ">>> Build artifact:"
	@find dist -type f -print0 | xargs -0 ls -t
	$(DURATION)

pack: build
	@echo ">>> Building syros-app image for deploy"
	docker build -t syros-app:$(APP_VERSION) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg GIT_BRANCH=$(GIT_BRANCH) \
		--build-arg APP_VERSION=$(APP_VERSION) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-f deploy.app.dockerfile .

	@echo ">>> Building syros-indexer image for deploy"
	docker build -t syros-indexer:$(APP_VERSION) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg GIT_BRANCH=$(GIT_BRANCH) \
		--build-arg APP_VERSION=$(APP_VERSION) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-f deploy.indexer.dockerfile .

	@echo ">>> Building syros-agent image for deploy"
	docker build -t syros-agent:$(APP_VERSION) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg GIT_BRANCH=$(GIT_BRANCH) \
		--build-arg APP_VERSION=$(APP_VERSION) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-f deploy.agent.dockerfile .

	@echo ">>> Images ready for deploy:"
	@docker images | grep syros
	$(DURATION)

run: pack
	@echo ">>> Starting syros-app container"
	@docker run -dp 8888:8888 --name syros-app-$(APP_VERSION) \
		syros-app:$(APP_VERSION) \
		-RethinkDB=$(RethinkDB) \
		-LogLevel=info

	@echo ">>> Starting syros-indexer container"
	@docker run -dp 8887:8887 --name syros-indexer-$(APP_VERSION) \
		syros-indexer:$(APP_VERSION) \
		-RethinkDB=$(RethinkDB) \
		-DatabaseStale=0 \
		-Nats=$(Nats) \
		-LogLevel=info

	@echo ">>> Starting syros-agent container"
	@docker run -dp 8886:8886 --name syros-agent-$(APP_VERSION) \
	    -v /var/run/docker.sock:/var/run/docker.sock:ro \
	    syros-agent:$(APP_VERSION) \
		-DockerApiAddresses=unix:///var/run/docker.sock \
		-Environment=dev \
		-Nats=$(Nats) \
		-LogLevel=info

	@echo ">>> syros-app logs:"
	@docker logs syros-app-$(APP_VERSION)
	@echo ">>> syros-indexer logs:"
	@docker logs syros-indexer-$(APP_VERSION)
	@echo ">>> syros-agent logs:"
	@docker logs syros-agent-$(APP_VERSION)
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

purge: clean
	@docker rm -f syros-app-$(APP_VERSION) syros-agent-$(APP_VERSION) syros-indexer-$(APP_VERSION) || true
	@docker rmi $$(docker images | awk '$$1 ~ /syros/ { print $$3 }') || true
	$(DURATION)

.PHONY: build



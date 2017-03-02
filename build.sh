#!/usr/bin/env bash
set -e

silent=0
output_dir="$(pwd)/dist"
BUILD_DATE=$(date -u +%Y-%m-%d_%H.%M.%S)
GIT_COMMIT=$(git rev-parse HEAD)
GIT_BRANCH=$(git symbolic-ref --short HEAD)
APP_VERSION="0.0.1"
START_TIME=$SECONDS

log() {
    if [ "${silent}" == 0 ]; then
        date=`date +%Y-%m-%d:%H:%M:%S`
        echo "$date INFO $1"
    fi
}

if [ -d "$output_dir" ]; then
    log "output directory found at ${output_dir}, removing content"
    rm -rf "$output_dir"
fi

log "building syros-ui-build image"
docker build -t syros-ui-build:${BUILD_DATE} -f build.deps.node.dockerfile .

log "building syros-ui"
docker run --rm  -v "$output_dir/ui:/usr/src/app/dist" syros-ui-build:${BUILD_DATE} bash -c "npm run build"
docker rmi syros-ui-build:${BUILD_DATE}

log "syros-ui build done"

log "building syros-services-build image"
docker build -t syros-services-build:${BUILD_DATE} -f build.deps.golang.dockerfile .

log "building syros-agent"
docker run --rm  -v "$output_dir:/go/dist" syros-services-build:${BUILD_DATE} go build -o /go/dist/agent github.com/stefanprodan/syros/agent

log "building syros-indexer"
docker run --rm  -v "$output_dir:/go/dist" syros-services-build:${BUILD_DATE} go build -o /go/dist/indexer github.com/stefanprodan/syros/indexer

log "building syros-api"
docker run --rm  -v "$output_dir:/go/dist" syros-services-build:${BUILD_DATE} go build -o /go/dist/api github.com/stefanprodan/syros/api

docker rmi syros-services-build:${BUILD_DATE}

log "syros-services build done"


log "building syros-app image for deploy"
docker build -t syros-app:${APP_VERSION} --build-arg GIT_COMMIT=${GIT_COMMIT} -f deploy.app.dockerfile .

log "syros-app:${APP_VERSION} image ready for deploy"

ELAPSED_TIME=$(($SECONDS - $START_TIME))
echo "$(($ELAPSED_TIME/60)) min $(($ELAPSED_TIME%60)) sec"
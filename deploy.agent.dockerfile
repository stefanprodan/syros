FROM alpine:3.5

ARG GIT_COMMIT=unkown
LABEL syros.revision=$GIT_COMMIT
ARG GIT_BRANCH=unkown
LABEL syros.branch=GIT_BRANCH
ARG APP_VERSION=unkown
LABEL syros.version=$APP_VERSION
ARG BUILD_DATE=unkown
LABEL syros.build=BUILD_DATE
LABEL syros.maintainer "Stefan Prodan"

EXPOSE 8886

COPY /dist/agent /syros/agent
RUN chmod 777 /syros/agent

WORKDIR /syros
ENTRYPOINT ["/syros/agent"]


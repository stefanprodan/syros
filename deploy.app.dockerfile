FROM alpine:3.5

ARG GIT_COMMIT=unkown
LABEL syros.revision=$GIT_COMMIT

ADD /dist /app
RUN chmod 777 /app/api

WORKDIR /app
CMD ["/app/api"]


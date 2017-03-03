FROM node:latest

# Create app directory
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

ADD /ui /usr/src/app

# deps
RUN npm install -q

# output
RUN mkdir -p /usr/src/app/dist
VOLUME /usr/src/app/dist

# #########################################
#  BASE IMAGE
FROM golang:1.14-alpine as base
RUN apk update && apk add --no-cache git

# RABBITMQ ENVS
ARG RABBIT_USER
ENV RABBIT_USER=$RABBIT_USER
ARG RABBIT_PASSWORD
ENV RABBIT_PASSWORD=$RABBIT_PASSWORD
ARG RABBIT_SERVER
ENV RABBIT_SERVER=$RABBIT_SERVER
ARG RABBIT_PORT
ENV RABBIT_PORT=$RABBIT_PORT
ARG RABBIT_VHOST
ENV RABBIT_VHOST=$RABBIT_VHOST

# LOG ENVS
ARG LOG_TYPE
ENV LOG_TYPE=$LOG_TYPE

# APP ENVS
ARG ENVIRONMENT
ENV ENVIRONMENT=$ENVIRONMENT

# #########################################
#  DEV IMAGE
FROM base as dev
# ADD in air file monitor
RUN go get -u github.com/cosmtrek/air
COPY ./docker/app/air.conf /tmp/air.conf
WORKDIR /go/src/github.com/arkrozycki/go-dead-rabbit
# SETUP LIVE RELOAD
CMD air -c /tmp/air.conf
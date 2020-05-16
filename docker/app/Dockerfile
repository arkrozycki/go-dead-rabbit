# #########################################
#  BASE IMAGE
FROM golang:1.14-alpine as base
RUN apk update && apk add --no-cache git

# #########################################
#  DEV IMAGE
FROM base as dev
# ADD in air file monitor
RUN go get -u github.com/cosmtrek/air
COPY ./docker/app/air.conf /tmp/air.conf
WORKDIR /go/src/github.com/arkrozycki/go-dead-rabbit
# SETUP LIVE RELOAD
CMD air -c /tmp/air.conf
FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/github.com/shauncampbell/zigbee2mqtt-mongodb-broker/
COPY . .

RUN go build -o /go/bin/zigbee2mqtt-mongodb-broker ./cmd/zigbee2mqtt-mongodb-broker

FROM alpine:3.12

COPY --from=builder /go/bin/zigbee2mqtt-mongodb-broker /go/bin/zigbee2mqtt-mongodb-broker
LABEL maintainer="Shaun Campbell <docker@shaun.scot>"

ENV BROKER_MONGODB_HOST ""
ENV BROKER_MONGODB_PORT 27017
ENV BROKER_MONGODB_DB "home"
ENV BROKER_MONGODB_USERNAME ""
ENV BROKER_MONGODB_PASSWORD ""
ENV BROKER_MQTT_HOST ""
ENV BROKER_MQTT_PORT 1883
ENV BROKER_MQTT_USERNAME ""
ENV BROKER_MQTT_PASSWORD ""

ENTRYPOINT ["./go/bin/zigbee2mqtt-mongodb-broker"]
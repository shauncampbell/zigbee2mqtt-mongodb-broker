clean:
	rm -rf zigbee2mqtt-mongodb-broker.*

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o zigbee2mqtt-mongodb-broker.linux_amd64 ./cmd/zigbee2mqtt-mongodb-broker
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o zigbee2mqtt-mongodb-broker.darwin_amd64 ./cmd/zigbee2mqtt-mongodb-broker
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o zigbee2mqtt-mongodb-broker.windows_amd64.exe ./cmd/zigbee2mqtt-mongodb-broker

docker:
	docker build -f ./cmd/zigbee2mqtt-mongodb-broker/Dockerfile -t shauncampbell/zigbee2mqtt-mongodb-broker:local .
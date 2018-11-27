BUILD=go build
VERSION := $(shell git describe --abbrev=4 --dirty --always --tags)

BUILD_CLIENT =$(BUILD) -ldflags "-X main.Version=$(VERSION)" -o ioex-wallet log.go config.go client.go
BUILD_SERVICE =$(BUILD) -ldflags "-X main.Version=$(VERSION)" -o service log.go config.go spvwallet.go main.go

all:
	$(BUILD_CLIENT)
	$(BUILD_SERVICE)

client:
	$(BUILD_CLIENT)

service:
	$(BUILD_SERVICE)
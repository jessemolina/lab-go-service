SHELL := /bin/bash

# ================================================================
# GO

go-run:
	go run main.go

go-build:
	go build -ldflags "-X main.build=local"


# ================================================================
# DOCKER

VERSION := 1.0

docker-build: docker-build-service

docker-build-service:
	docker build \
	-f zarf/docker/sales-api.dockerfile \
	-t service-amd64:$(VERSION) \
	--build-arg BUILD_REF=$(VERSION) \
	--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
	.

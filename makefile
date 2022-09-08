SHELL := /bin/bash

# ================================================================
# GO

go-run:
	go run main.go

go-build:
	go build -ldflags "-X main.build=local"

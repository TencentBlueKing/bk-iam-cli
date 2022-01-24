.PHONY: dep lint build linux

init:
	# go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.43.0

dep:
	go mod tidy
	go mod vendor

lint:
	export GOFLAGS=-mod=vendor
	golangci-lint run

build:
	go build -mod=vendor .

build-linux:
	GOOS=linux GOARCH=amd64 go build -mod=vendor .

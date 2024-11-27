build:
	go build -o runtimectl

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o runtimectl

.PHONY: build build-linux
.PHONY: all get clean build

GO ?= go

all: build

build: get
	${GO} build -o bin/license-server main.go

build-all:
    ${echo "building binary for multiple OS"}
    GOOS=linux GOARCH=arm ${GO} build -o bin/license-server-linux-arm main.go
    GOOS=linux GOARCH=arm64 ${GO} build -o bin/license-server-linux-arm64 main.go
    GOOS=linux GOARCH=amd64 ${GO} build -o bin/license-server-linux-amd64 main.go
    GOOS=linux GOARCH=386 ${GO} build -o bin/license-server-linux-386 main.go
    GOOS=darwin GOARCH=amd64 ${GO} build -o bin/license-server-mac-amd64 main.go
    GOOS=darwin GOARCH=386 ${GO} build -o bin/license-server-mac-386 main.go
    GOOS=freebsd GOARCH=386 ${GO} build -o bin/license-server-freebsd-386 main.go
    GOOS=windows GOARCH=386 ${GO} build -o bin/license-server-windows-386 main.go

get:
	dep ensure

clean:
	@rm -rf license-server

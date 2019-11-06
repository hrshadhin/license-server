.PHONY: all get clean build

GO ?= go

all: build

build: get
	${GO} build -o bin/license-server main.go

get:
	dep ensure

clean:
	@rm -rf bin/license-server

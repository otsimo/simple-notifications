.PHONY: default build release clean

default: build

build: clean vet
	script/build

cross: clean vet
	script/build cross

release: clean vet
	script/build docker
	script/release

run: build
	script/run

fmt:
	go fmt ./src/...

vet:
	go vet ./src/...

clean:
	rm -rf bin

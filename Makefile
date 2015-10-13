.PHONY: default build release proto clean

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
	goimports -w src

vet:
	go vet ./src/...

proto:
	script/proto

clean:
	rm -rf bin

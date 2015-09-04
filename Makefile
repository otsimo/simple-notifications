.PHONY: default build release clean

default: build

build: clean
	script/build

cross: clean
	script/build cross

release: clean
	script/build docker
	script/release

run: build
	script/run

clean:
	rm -rf bin

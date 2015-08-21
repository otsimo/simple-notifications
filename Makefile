
default: build

build: clean
	script/build

release: build
	script/release

clean:
	rm -rf bin

all: test

clean:
	rm -f highLowGameServer 

install: prepare
	godep go install

prepare:
	go get github.com/tools/godep

build: prepare
	godep save
	godep go build

test: prepare build
	echo "no tests"

.PHONY: install prepare build test

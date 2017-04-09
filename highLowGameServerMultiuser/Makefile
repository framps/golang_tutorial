GLIDE_VERSION := $(shell glide -v --version 2>/dev/null)

all: test

clean:
	rm -f highLowGameServerMultiuser

install: prepare
	glide update

prepare:
ifndef GLIDE_VERSION
		curl https://glide.sh/get | sh
endif

build: prepare
	go build

test: prepare build
	echo "no tests"

.PHONY: install prepare build test

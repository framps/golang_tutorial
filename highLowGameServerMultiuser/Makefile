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

run: prepare
	go run server.go

deploy:
	cf login
	cf push
	cf apps

.PHONY: install test build

install:
	go install

test:
	go test quad9/*

build:
	go build

run: build
	./q9aas
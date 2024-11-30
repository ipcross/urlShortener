.PHONY: all
all: build

.PHONY: build
build:
	go build -o bin/shortener ./cmd/shortener/

test: build
	shortenertest -test.v -test.run=^TestIteration$(n)$$ -binary-path=bin/shortener
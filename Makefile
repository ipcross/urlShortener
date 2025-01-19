.PHONY: all
all: run

.PHONY: run
run: build
	./bin/shortener -d 'postgres://postgres:qwerty@localhost:5432/go_urlshortener?sslmode=disable'

.PHONY: build
build:
	go build -o bin/shortener ./cmd/shortener/

test: build
	go test ./...

stb: build
	shortenertestbeta -test.v -test.run=^TestIteration$(n)$$ -binary-path=bin/shortener -source-path=. \
      -server-port=3010 -file-storage-path=/tmp/storage.json -database-dsn='postgres://postgres:qwerty@localhost:5432/go_urlshortener?sslmode=disable'

linter:
	go vet -vettool=$$(which statictest) ./...

GOLANGCI_LINT_CACHE?=/tmp/praktikum-golangci-lint-cache

.PHONY: golangci-lint-run
golangci-lint-run: _golangci-lint-rm-unformatted-report

.PHONY: _golangci-lint-reports-mkdir
_golangci-lint-reports-mkdir:
	mkdir -p ./golangci-lint

.PHONY: _golangci-lint-run
_golangci-lint-run: _golangci-lint-reports-mkdir
	-docker run --rm \
    -v $(shell pwd):/app \
    -v $(GOLANGCI_LINT_CACHE):/root/.cache \
    -w /app \
    golangci/golangci-lint:v1.57.2 \
        golangci-lint run \
            -c .golangci.yml \
	> ./golangci-lint/report-unformatted.json

.PHONY: _golangci-lint-format-report
_golangci-lint-format-report: _golangci-lint-run
	cat ./golangci-lint/report-unformatted.json | jq > ./golangci-lint/report.json

.PHONY: _golangci-lint-rm-unformatted-report
_golangci-lint-rm-unformatted-report: _golangci-lint-format-report
	rm ./golangci-lint/report-unformatted.json

.PHONY: golangci-lint-clean
golangci-lint-clean:
	sudo rm -rf ./golangci-lint

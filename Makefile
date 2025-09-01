.PHONY: all
all: vet test build

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: build
build:
	go build ./cmd/gap

.PHONY: lint
lint:
	golangci-lint run

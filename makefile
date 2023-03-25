.DEFAULT_GOAL:= build

fmt:
	go fmt ./...
.PHONY: fmt

lint:
	golangci-lint run
.PHONY: lint

vet:
	go vet ./...
.PHONY: vet

test:
	go test ./...
.PHONY: test

build:
	go build -o bin/ ./...
.PHONY: build

clean:
	rm -rf bin/
.PHONY: clean
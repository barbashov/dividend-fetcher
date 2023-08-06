GOLANGCILINT_VERSION=v1.53.3

all: build

build:
	go build cmd/main.go

lint:
	docker run --rm \
		-v ~/.cache/golangci-lint/$(GOLANGCILINT_VERSION):/root/.cache \
		-v `pwd`:/app -w /app \
		golangci/golangci-lint:$(GOLANGCILINT_VERSION) \
		golangci-lint run

test:
	go test ./...
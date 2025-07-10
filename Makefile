.PHONY: build run test clean docker-build docker-run deps lint

APP_NAME := weather-exporter
DOCKER_IMAGE := $(APP_NAME):latest
GO_FILES := $(shell find . -name '*.go' -type f)

build:
	go build -o bin/$(APP_NAME) ./cmd/$(APP_NAME)

run: build
	./bin/$(APP_NAME) -config configs/config.yaml

deps:
	go mod download
	go mod tidy

test:
	go test -v -race -coverprofile=coverage.out ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/
	rm -f coverage.out

docker-build:
	docker build -t $(DOCKER_IMAGE) -f docker/Dockerfile .

docker-run: docker-build
	docker run -p 8080:8080 -e OPENWEATHER_API_KEY=$(OPENWEATHER_API_KEY) $(DOCKER_IMAGE)

help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  run          - Build and run the application"
	@echo "  deps         - Download and tidy dependencies"
	@echo "  test         - Run tests"
	@echo "  lint         - Run linters"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Build and run Docker container"
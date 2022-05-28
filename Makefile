.PHONY: all build clean run check cover lint docker help
BIN_FILE=basketball-go
all: check build
build:
	@go build -o "${BIN_FILE}"
clean:
	@go clean
	rm --force "xx.out"
test:
	@go test
check:
	@go fmt ./
	@go vet ./
cover:
	@go test -coverprofile xx.out
	@go tool cover -html=xx.out
run:
	./"${BIN_FILE}"
lint:
	golangci-lint run --enable-all
docker:
	@docker build -t lzzzzl/basketball-go:latest .
help:
	@echo "make compile: packages and dependencies"
	@echo "make build: compile go code to generate binary file"
	@echo "make clean: clean up intermediate object files"
	@echo "make test: execute test case"
	@echo "make check: format go code"
	@echo "make cover: check test coverage"
	@echo "make run: run code"
	@echo "make lint: perform code inspection"
	@echo "make docker: build docker image"
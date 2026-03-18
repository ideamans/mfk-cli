BINARY_NAME := mfk
MODULE := github.com/ideamans/mfk-cli

.PHONY: build run test lint clean snapshot

build:
	go build -o $(BINARY_NAME) .

run: build
	./$(BINARY_NAME) $(ARGS)

test:
	go test ./...

lint:
	go vet ./...

clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/

snapshot:
	goreleaser release --snapshot --clean

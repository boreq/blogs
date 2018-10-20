VERSION = `git rev-parse HEAD`
DATE = `date --iso-8601=seconds`
LDFLAGS =  -X github.com/boreq/blogs/cmd/blogs/commands.buildCommit=$(VERSION)
LDFLAGS += -X github.com/boreq/blogs/cmd/blogs/commands.buildDate=$(DATE)

all: build

build:
	mkdir -p build
	go build -ldflags "$(LDFLAGS)" -o ./build/blogs ./cmd/blogs

doc:
	@echo "http://localhost:6060/pkg/github.com/boreq/blogs/"
	godoc -http=:6060

test:
	go test ./...

test-verbose:
	go test -v ./...

test-short:
	go test -short ./...

clean:
	rm -rf ./build

.PHONY: all build doc test test-verbose test-short clean

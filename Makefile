all: build

build:
	mkdir -p build
	go build -o ./build/blogs ./main

run:
	./main/main

doc:
	@echo "http://localhost:6060/pkg/github.com/boreq/blogs/"
	godoc -http=:6060

test:
	go test ./...

test-verbose:
	go test -v ./...

test-short:
	go test -short ./...

bench:
	go test -v -run=XXX -bench=. ./...

clean:
	rm -rf ./build

.PHONY: all build run doc test test-verbose test-short bench clean

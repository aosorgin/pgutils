.PHONY: build-all-swagger

GIT_COMMIT=$(shell git log -1 --pretty=format:"%H")

init:
	# Install golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.27.0

dep:
	go mod vendor

build-pgquery:
	go build -o bin/pgquery -mod vendor ./cmd/main.go

build-all: build-pgquery

fmt:
	bash -c 'diff -u <(echo -n) <(gofmt -l . |grep -v vendor)'

lint:
	go vet -mod vendor ./...
	./bin/golangci-lint run ./...

test:
	go test -mod vendor ./...

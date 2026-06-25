BINARY := seekai
VERSION ?= dev
OUT_DIR ?= dist

.PHONY: all build test fmt fmt-check tidy clean release

all: fmt-check test build

build:
	GOFLAGS="-trimpath" go build -ldflags "-s -w -X main.version=$(VERSION)" -o $(BINARY) .

test:
	go test ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './dist/*')

fmt-check:
	@test -z "$$(gofmt -l $$(find . -name '*.go' -not -path './dist/*'))"

tidy:
	go mod tidy

release:
	VERSION=$(VERSION) OUT_DIR=$(OUT_DIR) scripts/build-release.sh

clean:
	rm -rf $(OUT_DIR) $(BINARY) $(BINARY).exe coverage.out coverage.html

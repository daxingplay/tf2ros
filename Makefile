NAME=tf2ros
ARCH=$(shell uname -m)
VERSION=0.0.1
ITERATION := 1

SOURCE_FILES?=$$(go list ./... | grep -v /vendor/)
TEST_PATTERN?=.
TEST_OPTIONS?=

BIN_DIR := $(CURDIR)/bin

ci: prepare test

prepare: prepare.metalinter
	GOBIN=$(BIN_DIR) go install github.com/buildkite/github-release
	GOBIN=$(BIN_DIR) go install github.com/mitchellh/gox@latest
	GOBIN=$(BIN_DIR) go install github.com/axw/gocov/gocov@latest
	GOBIN=$(BIN_DIR) go install golang.org/x/tools/cmd/cover@latest

# Gometalinter is deprecated and broken dependency so let's use with GO111MODULE=off
prepare.metalinter:
	GO111MODULE=off go get -u github.com/alecthomas/gometalinter
	GO111MODULE=off gometalinter --fast --install

mod:
	@go mod download
	@go mod tidy

compile: mod
	@rm -rf build/
	@$(BIN_DIR)/gox -ldflags "-X main.Version=$(VERSION)" \
	-osarch="darwin/arm" \
	-osarch="darwin/arm64" \
	-osarch="darwin/amd64" \
	-osarch="linux/i386" \
	-osarch="linux/amd64" \
	-osarch="windows/amd64" \
	-osarch="windows/i386" \
	-output "build/{{.Dir}}_$(VERSION)_{{.OS}}_{{.Arch}}/$(NAME)" \
	${SOURCE_FILES}

# Run all the linters
lint:
	@gometalinter --vendor ./...

# gofmt and goimports all go files
fmt:
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

dist:
	$(eval FILES := $(shell ls build))
	@rm -rf dist && mkdir dist
	@for f in $(FILES); do \
		(cd $(shell pwd)/build/$$f && tar -cvzf ../../dist/$$f.tar.gz *); \
		(cd $(shell pwd)/dist && shasum -a 512 $$f.tar.gz > $$f.sha512); \
		echo $$f; \
	done

release:
	@$(BIN_DIR)/github-release "v$(VERSION)" dist/* --commit "$(git rev-parse HEAD)" --github-repository daxingplay/$(NAME)

test:
	@$(BIN_DIR)/gocov test $(SOURCE_FILES) | $(BIN_DIR)/gocov report

clean:
	@rm -fr ./build

generate-mocks:
	mockery -dir pkg/executor --all

.PHONY: default prepare.metalinter prepare mod compile lint fmt dist release test clean generate-mocks

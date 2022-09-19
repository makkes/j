BINARY_NAME := jump
VERSION := $(shell git describe --tags $(git rev-parse HEAD))
PLATFORMS := linux darwin
TEST_FLAGS ?= ""
os = $(word 1, $@)

.PHONY: build
build:
	go build -o jump ./cmd/jump

.PHONY: test
test:
	go test $(TEST_FLAGS) -v ./...

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	GOOS=$(os) GOARCH=amd64 go build -v -o $(BINARY_NAME) ./cmd/jump
	tar -czf j_$(VERSION)_$(os)_x86_64.tar.gz $(BINARY_NAME) j.sh j_completion

.PHONY: release
release: $(PLATFORMS)

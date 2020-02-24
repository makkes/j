BINARY_NAME := jump
VERSION := v1.0.5
PLATFORMS := linux darwin
os = $(word 1, $@)

.PHONY: build
build:
	go build -o jump ./cmd/jump

.PHONY: test
test:
	go test -v ./...

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	GOOS=$(os) GOARCH=amd64 go build -v -o $(BINARY_NAME)_$(VERSION)_$(os)_amd64 ./cmd/jump

.PHONY: release
release: $(PLATFORMS)

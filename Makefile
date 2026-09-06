.PHONY: build build-legacy build-a2a install install-legacy test test-legacy vet vet-legacy lint lint-legacy check check-legacy check-all coverage coverage-legacy

GO ?= go
OFFICIAL_TAGS ?= official_sdk

build:
	GOWORK=off $(GO) build -tags=$(OFFICIAL_TAGS) -o tmux-mcp ./cmd/tmux-mcp

build-legacy:
	GOWORK=off $(GO) build -o tmux-mcp-legacy ./cmd/tmux-mcp

build-a2a:
	GOWORK=off $(GO) build -tags=$(OFFICIAL_TAGS) -o tmux-a2a ./cmd/tmux-a2a

install:
	GOWORK=off $(GO) install -tags=$(OFFICIAL_TAGS) ./cmd/tmux-mcp

install-legacy:
	GOWORK=off $(GO) install ./cmd/tmux-mcp

test:
	GOWORK=off $(GO) test -tags=$(OFFICIAL_TAGS) ./... -count=1

test-legacy:
	GOWORK=off $(GO) test ./... -count=1

vet:
	GOWORK=off $(GO) vet -tags=$(OFFICIAL_TAGS) ./...

vet-legacy:
	GOWORK=off $(GO) vet ./...

lint:
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run --build-tags=$(OFFICIAL_TAGS) ./... || \
	(command -v staticcheck >/dev/null 2>&1 && staticcheck -tags=$(OFFICIAL_TAGS) ./... || echo "no linter installed, skipping")

lint-legacy:
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || \
	(command -v staticcheck >/dev/null 2>&1 && staticcheck ./... || echo "no linter installed, skipping")

check: build vet test

ci: check

check-legacy: build-legacy vet-legacy test-legacy

check-all: check check-legacy

coverage:
	GOWORK=off $(GO) test -tags=$(OFFICIAL_TAGS) ./... -count=1 -coverprofile=coverage.out
	GOWORK=off $(GO) tool cover -func=coverage.out

coverage-legacy:
	GOWORK=off $(GO) test ./... -count=1 -coverprofile=coverage-legacy.out
	GOWORK=off $(GO) tool cover -func=coverage-legacy.out

HG_PIPELINE_MK ?= $(or $(wildcard $(abspath $(CURDIR)/../dotfiles/make/pipeline.mk)),$(wildcard $(HOME)/hairglasses-studio/dotfiles/make/pipeline.mk))
-include $(HG_PIPELINE_MK)


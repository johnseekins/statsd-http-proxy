.PHONY: help all test fmt lint vendor-update build clean

GH_ADDR := github.com/johnseekins
NAME := statsd-http-proxy
GO_VER := 1.18.1
CURRENT_UID := $(shell id -u)
CURRENT_GID := $(shell id -g)
BUILDTIME ?= $(shell date)
BUILDUSER ?= $(shell id -u -n)
PKG_TAG ?= $(shell git tag -l --points-at HEAD)
ifeq ($(PKG_TAG),)
PKG_TAG := $(shell echo $$(git describe --long --all | tr '/' '-')$$(git diff-index --quiet HEAD -- || echo '-dirty-'$$(git diff-index -u HEAD | openssl sha1 | cut -c 10-17)))
endif

all: test build ## format, lint, and build the package

help:
	@echo "Makefile targets:"
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' Makefile \
	| sed -n 's/^\(.*\): \(.*\)##\(.*\)/    \1 :: \3/p' \
	| column -t -c 1  -s '::'

setup:
	docker pull golang:$(GO_VER)
	docker pull golangci/golangci-lint:latest

test: fmt lint ## run tests
	docker run --rm --user $(CURRENT_UID):$(CURRENT_GID) -v $(CURDIR):/app:z -w /app golang:$(GO_VER) go test -mod=vendor ./...

fmt: ## only run gofmt
	docker run --rm --user $(CURRENT_UID):$(CURRENT_GID)-v $(CURDIR):/app:z -w /app golang:$(GO_VER) gofmt -l -w -s *.go

lint: ## run all linting steps
	docker run --rm --user $(CURRENT_UID):$(CURRENT_GID) -v $(CURDIR):/app:z -w /app golangci/golangci-lint:latest golangci-lint run

vendor-update: ## update vendor dependencies
	rm -rf go.mod go.sum vendor/
	docker run --rm --user $(CURRENT_UID):$(CURRENT_GID) -v $(CURDIR):/app:z -w /app golang:$(GO_VER) go mod init $(GH_ADDR)/$(NAME)
	docker run --rm --user $(CURRENT_UID):$(CURRENT_GID) -v $(CURDIR):/app:z -w /app golang:$(GO_VER) go mod tidy
	docker run --rm --user $(CURRENT_UID):$(CURRENT_GID) -v $(CURDIR):/app:z -w /app golang:$(GO_VER) go mod vendor

build: ## actually build package
	docker run --rm --user $(CURRENT_UID):$(CURRENT_GID) -v $(CURDIR):/app:z -w /app golang:$(GO_VER) go build -tags static,netgo -mod=vendor -ldflags="-X 'main.Version=$(PKG_TAG)' -X 'main.BuildUser=$(BUILDUSER)' -X 'main.BuildTime=$(BUILDTIME)'" -o $(NAME) .

clean: ## remove build artifacts
	rm -f $(NAME)

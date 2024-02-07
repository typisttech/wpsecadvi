MAKEFLAGS+=-r -R

# Project
BINARY_NAME=wpsecadvi
DIST=$(CURDIR)/dist
DIST_BINARY=$(DIST)/$(BINARY_NAME)_$(shell $(GOCMD) env GOOS)_$(shell $(GOCMD) env GOARCH)/${BINARY_NAME}
OUTPUT_BINARY=$(CURDIR)/$(BINARY_NAME)

# Go
GOCMD?=$(shell which go)
GOTEST?=$(GOCMD) test
GOBIN?=$(if $(shell go env GOBIN),$(shell go env GOBIN),$(shell go env GOPATH)/bin)

# Tool
GOLANGCI_LINT?=$(GOBIN)/golangci-lint
GORELEASER?=$(GOBIN)/goreleaser
GOVULNCHECK?=$(GOBIN)/govulncheck

# Color
GREEN   := $(shell tput -Txterm setaf 2)
YELLOW  := $(shell tput -Txterm setaf 3)
BLUE    := $(shell tput -Txterm setaf 4)
MAGENTA := $(shell tput -Txterm setaf 5)
CYAN    := $(shell tput -Txterm setaf 6)
WHITE   := $(shell tput -Txterm setaf 7)
RESET   := $(shell tput -Txterm sgr0)


## Common:
.PHONY: default
default: build ## Alias build

.PHONY: all
all: check test build-all ## Alias of check, test, build-all

.PHONY: dev
dev: check test ## Alias of check and test

.PHONY: clean
clean: test-clean build-clean ## Clean up all generated files


## Build:
.PHONY: build
build: build-clean goreleaser-check ## Build binary for current GOOS and GOARCH
	@echo "\n${GREEN}====>${RESET} Building binary for current GOOS and GOARCH to ${CYAN}$(DIST)${RESET}..."
	$(GORELEASER) build --clean --snapshot --single-target
	@echo "\n${GREEN}====>${RESET} Copying binary from ${CYAN}$(DIST_BINARY)${RESET} to ${CYAN}$(OUTPUT_BINARY)${RESET}..."
	cp $(DIST_BINARY) $(OUTPUT_BINARY)
	@echo "\n${GREEN}====>${RESET} Printing binary version information"
	$(OUTPUT_BINARY) --version

.PHONY: build-all
build-all: build-clean goreleaser-check ## Build binaries for all supported targets to ./dist
	@echo "\n${GREEN}====>${RESET} Building binaries to ${CYAN}$(DIST)${RESET}..."
	$(GORELEASER) build --clean --snapshot

.PHONY: build-clean
build-clean: ## Clean up binaries and build outputs
	@echo "\n${GREEN}====>${RESET} Removing ${CYAN}$(DIST)${RESET}..."
	rm -rf $(DIST)
	@echo "\n${GREEN}====>${RESET} Removing ${CYAN}$(OUTPUT_BINARY)${RESET}..."
	rm -f $(OUTPUT_BINARY)


## Test:
.PHONY: test
test: ## Run tests
	@echo "\n${GREEN}====>${RESET} Running tests..."
	$(GOTEST) -failfast -race -shuffle=on ./...

.PHONY: coverage.out
coverage.out:
	@echo "\n${GREEN}====>${RESET} Generating coverage profile to ${CYAN}$(CURDIR)/coverage.out${RESET}..."
	$(GOTEST) -v -count=1 -race -shuffle=on -cover -coverprofile=coverage.out ./...

.PHONY: test-coverage
test-coverage: coverage.out ## Run all tests with coverage, then open the coverage report in default web browser
	@echo "\n${GREEN}====>${RESET} Opening the coverage profile in default web browser..."
	$(GOCMD) tool cover -html=$(CURDIR)/coverage.out

.PHONY: test-clean
test-clean: ## Clean up generated test files
	@echo "\n${GREEN}====>${RESET} Removing ${CYAN}$(CURDIR)/coverage.out${RESET}..."
	rm -f $(CURDIR)/coverage.out


## Check:
.PHONY: check
check: golangci-lint goreleaser-check govulncheck ## Run all checks

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Run golangci-lint linters
	@echo "\n${GREEN}====>${RESET} Running golangci-lint linters..."
	$(GOLANGCI_LINT) run

.PHONY: goreleaser-check
goreleaser-check: $(GORELEASER) ## Valid goreleaser configuration
	@echo "\n${GREEN}====>${RESET} Validating goreleaser configuration..."
	$(GORELEASER) check

.PHONY: govulncheck
govulncheck: $(GOVULNCHECK) ## Scan for vulnerable dependencies
	@echo "\n${GREEN}====>${RESET} Scanning for vulnerable dependencies..."
	$(GOVULNCHECK) ./...


## Tool:
tool: $(GOBIN)/golangci-lint $(GOBIN)/goreleaser $(GOBIN)/govulncheck ## Install all tools to $GOBIN

$(GOBIN)/golangci-lint: PACKAGE=github.com/golangci/golangci-lint/cmd/golangci-lint@latest
$(GOBIN)/goreleaser: PACKAGE=github.com/goreleaser/goreleaser@latest
$(GOBIN)/govulncheck: PACKAGE=golang.org/x/vuln/cmd/govulncheck@latest

$(GOBIN)/%:
	@echo "\n${GREEN}====>${RESET} Installing ${MAGENTA}$(PACKAGE)${RESET} to ${CYAN}$(GOBIN)${RESET}..."
	GOBIN=$(GOBIN) $(GOCMD) install $(PACKAGE)


## Help:
.PHONY: help
help: ## Show this help
	@echo ''
	@echo 'Usage:'
	@echo '  ${CYAN}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${GREEN}%-18s${WHITE}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "\n  ${YELLOW}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

.PHONY: debug
debug: ## List variables
	@echo "BINARY_NAME   ${CYAN}$(BINARY_NAME)${RESET}"
	@echo "DIST          ${CYAN}$(DIST)${RESET}"
	@echo "DIST_BINARY   ${CYAN}$(DIST_BINARY)${RESET}"
	@echo "OUTPUT_BINARY ${CYAN}$(OUTPUT_BINARY)${RESET}"
	@echo "GOCMD         ${CYAN}$(GOCMD)${RESET}"
	@echo "GOTEST        ${CYAN}$(GOTEST)${RESET}"
	@echo "GOBIN         ${CYAN}$(GOBIN)${RESET}"
	@echo "GOLANGCI_LINT ${CYAN}$(GOLANGCI_LINT)${RESET}"
	@echo "GORELEASER    ${CYAN}$(GORELEASER)${RESET}"
	@echo "GOVULNCHECK   ${CYAN}$(GOVULNCHECK)${RESET}"
	@echo "MAKEFLAGS     ${CYAN}$(MAKEFLAGS)${RESET}"

# Makefile

# Default shell
SHELL := /bin/zsh

# Default target
.DEFAULT_GOAL := help

# Colors for output
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: help build clean

help: ## Show this help message
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  ${YELLOW}%-15s${RESET} %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build Rick CLI
	@echo "${GREEN}Building Rick...${RESET}"
	@chmod +x ./scripts/build.sh
	@./scripts/build.sh

clean: ## Clean build artifacts
	@echo "${GREEN}Cleaning...${RESET}"
	@rm -f rick
	@go clean


install: build ## Install Rick CLI
	@echo "${GREEN}Installing Rick...${RESET}"
	@sudo cp rick /usr/local/bin/

uninstall: ## Uninstall Rick CLI
	@echo "${GREEN}Uninstalling Rick...${RESET}"
	@sudo rm -f /usr/local/bin/rick

test: ## Run tests
	@echo "${GREEN}Running tests...${RESET}"
	@go test ./...

fmt: ## Format code
	@echo "${GREEN}Formatting code...${RESET}"
	@go fmt ./...

.PHONY: install uninstall test fmt
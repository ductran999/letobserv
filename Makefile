PKG_SCRIPTS=scripts

default: help

help: ## Show help for each of the Makefile commands
	@awk 'BEGIN \
		{FS = ":.*##"; printf "Usage: make ${cyan}<command>\n${white}Commands:\n"} \
		/^[a-zA-Z_-]+:.*?##/ \
		{ printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } \
		/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' \
		$(MAKEFILE_LIST)

.PHONY: tidy
tidy: ## Tidy up the go.mod
	go mod tidy

.PHONY: lint
lint: ## Run linters
	golangci-lint run --timeout 10m --config .golangci.yml

.PHONY: deps
deps: ## install library
	go install github.com/vektra/mockery/v3@v3.4.0
	go install github.com/wadey/gocovmerge@latest

stop: ## Stop demo 
	docker-compose down

restart: ## Reset demo
	docker-compose down
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Copied .env.example → .env"; \
	else \
		echo ".env already exists, skipping copy."; \
	fi
	docker-compose up -d

setup: ## Setup demo dependencies
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Copied .env.example → .env"; \
	else \
		echo ".env already exists, skipping copy."; \
	fi
	docker-compose up -d
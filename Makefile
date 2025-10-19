.PHONY: help test lint coverage

# Couleurs
GREEN := \033[0;32m
YELLOW := \033[0;33m
NC := \033[0m

help: ## Affiche cette aide
	@echo ""
	@echo "${GREEN}KTN-Linter - Commandes disponibles:${NC}"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { printf "  ${YELLOW}%-15s${NC} %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""

test: ## Exécute les tests avec couverture
	@go test -v ./...
	@$(MAKE) coverage

coverage: ## Génère le rapport de couverture (COVERAGE.MD)
	@./scripts/generate-coverage.sh

lint: ## Lance le linter sur le projet (exclut *_test.go)
	@go run ./cmd/ktn-linter ./...

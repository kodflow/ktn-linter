.PHONY: help test lint coverage build install validate fmt

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

install: ## Installe ktn-linter depuis GitHub releases (fallback: compile)
	@echo "${GREEN}Installation de ktn-linter...${NC}"
	@bash scripts/install-ktn-linter.sh

build: ## Compile le binaire ktn-linter dans builds/
	@echo "${GREEN}Compilation de ktn-linter...${NC}"
	@mkdir -p builds
	@go build -buildvcs=false -o builds/ktn-linter ./cmd/ktn-linter
	@echo "${GREEN}✅ Binaire créé: builds/ktn-linter${NC}"

fmt: ## Formate le code Go avec go fmt
	@echo "${GREEN}Formatage du code...${NC}"
	@go fmt ./...
	@echo "${GREEN}✅ Code formaté${NC}"

test: ## Exécute les tests avec couverture
	@go test -v ./...
	@$(MAKE) coverage

coverage: ## Génère le rapport de couverture (COVERAGE.MD)
	@./scripts/generate-coverage.sh

lint: ## Lance le linter sur le projet (exclut *_test.go)
	@go run -buildvcs=false ./cmd/ktn-linter lint ./...

validate: build ## Valide que tous les testdata good.go/bad.go sont corrects
	@./scripts/validate-testdata.sh

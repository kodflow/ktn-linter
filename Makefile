.PHONY: help test lint lint-testdata coverage build validate fmt install

# Couleurs
GREEN := \033[0;32m
YELLOW := \033[0;33m
NC := \033[0m

# Version from git tags (fallback to dev if no tags)
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")

help: ## Affiche cette aide
	@echo ""
	@echo "${GREEN}KTN-Linter - Commandes disponibles:${NC}"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { printf "  ${YELLOW}%-15s${NC} %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""

build: ## Compile le binaire ktn-linter dans builds/
	@echo "${GREEN}Compilation de ktn-linter ${VERSION}...${NC}"
	@mkdir -p builds
	@go build -buildvcs=false -ldflags="-X main.Version=${VERSION}" -o builds/ktn-linter ./cmd/ktn-linter
	@echo "${GREEN}✅ Binaire créé: builds/ktn-linter (${VERSION})${NC}"

install: build ## Compile et installe ktn-linter dans /usr/local/bin
	@echo "${GREEN}Installation de ktn-linter...${NC}"
	@if [ -w /usr/local/bin ]; then \
		cp builds/ktn-linter /usr/local/bin/ktn-linter && \
		echo "${GREEN}✅ ktn-linter installé dans /usr/local/bin${NC}"; \
	elif sudo -n true 2>/dev/null; then \
		sudo cp builds/ktn-linter /usr/local/bin/ktn-linter && \
		echo "${GREEN}✅ ktn-linter installé dans /usr/local/bin (avec sudo)${NC}"; \
	else \
		mkdir -p ~/.local/bin && \
		cp builds/ktn-linter ~/.local/bin/ktn-linter && \
		echo "${GREEN}✅ ktn-linter installé dans ~/.local/bin${NC}" && \
		echo "${YELLOW}⚠ Assurez-vous que ~/.local/bin est dans votre PATH${NC}"; \
	fi
	@ktn-linter --version

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

lint-testdata: build ## Alias pour validate - vérifie les testdata (good=0 erreur, bad=UNIQUEMENT sa règle)
	@./scripts/validate-testdata.sh

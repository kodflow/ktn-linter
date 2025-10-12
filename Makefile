.PHONY: help check-go build lint test clean deps install-tools

# Couleurs pour le output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

help: ## Affiche cette aide
	@echo ""
	@echo "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
	@echo "${BLUE}║                    KTN-LINTER                              ║${NC}"
	@echo "${BLUE}║          Linter Go pour les bonnes pratiques               ║${NC}"
	@echo "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
	@echo ""
	@echo "${GREEN}Commandes disponibles:${NC}"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { printf "  ${YELLOW}%-20s${NC} %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""
	@echo "${YELLOW}Prérequis:${NC}"
	@echo "  - Go 1.23+"
	@echo "  - golangci-lint v2+ (pour VSCode)"
	@echo ""

check-go: ## Vérifie que Go est installé
	@command -v go >/dev/null 2>&1 || { \
		echo "${RED}❌ Go n'est pas installé${NC}"; \
		echo "Installez Go: https://go.dev/doc/install"; \
		exit 1; \
	}
	@echo "${GREEN}✅ Go $(shell go version | awk '{print $$3}')${NC}"

deps: check-go ## Installe les dépendances Go
	@echo "${BLUE}📦 Installation des dépendances...${NC}"
	@go mod download
	@go mod tidy
	@echo "${GREEN}✅ Dépendances installées${NC}"

build: check-go deps ## Compile le linter (utilisé par golangci-lint-wrapper)
	@echo "${BLUE}🔨 Compilation du linter...${NC}"
	@mkdir -p builds
	@cd src/cmd/ktn-linter && go build -o ../../../builds/ktn-linter
	@echo "${GREEN}✅ Linter compilé: builds/ktn-linter${NC}"
	@echo "${YELLOW}💡 VSCode utilise: bin/golangci-lint-wrapper${NC}"

lint: build ## Exécute le linter sur tests/
	@echo ""
	@echo "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
	@echo "${BLUE}║              ANALYSE DES TESTS                             ║${NC}"
	@echo "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
	@echo ""
	@echo "${YELLOW}📂 Analyse de tests/source/ (code avec erreurs)${NC}"
	@echo ""
	@./builds/ktn-linter ./tests/source/... 2>&1 | head -30 || true
	@echo ""
	@echo "${YELLOW}... (sortie tronquée)${NC}"
	@echo ""
	@echo "${YELLOW}📂 Analyse de tests/target/ (code conforme)${NC}"
	@echo ""
	@./builds/ktn-linter ./tests/target/... && echo "${GREEN}✅ Aucune erreur détectée !${NC}" || echo "${RED}❌ Erreurs détectées${NC}"
	@echo ""
	@echo "${YELLOW}💡 Dans VSCode, golangci-lint utilise bin/golangci-lint-wrapper${NC}"
	@echo ""

test: check-go ## Exécute les tests unitaires
	@echo "${BLUE}🧪 Exécution des tests...${NC}"
	@cd src && go test -v ./...

clean: ## Nettoie les fichiers compilés
	@echo "${BLUE}🧹 Nettoyage...${NC}"
	@rm -rf builds/
	@echo "${GREEN}✅ Nettoyage terminé${NC}"

install-tools: ## Installe golangci-lint (optionnel)
	@echo "${BLUE}📦 Installation de golangci-lint...${NC}"
	@command -v golangci-lint >/dev/null 2>&1 && { \
		echo "${GREEN}✅ golangci-lint déjà installé${NC}"; \
	} || { \
		echo "${YELLOW}Installation de golangci-lint...${NC}"; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin; \
		echo "${GREEN}✅ golangci-lint installé${NC}"; \
	}

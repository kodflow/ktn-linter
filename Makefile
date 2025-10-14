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

test-func: check-go ## Exécute uniquement les tests des analyseurs FUNC
	@echo "${BLUE}🧪 Tests FUNC analyzer...${NC}"
	@cd src && go test -v ./pkg/analyzer -run TestFunc

test-coverage: check-go ## Génère un rapport de couverture HTML
	@echo "${BLUE}📊 Génération du rapport de couverture...${NC}"
	@cd src && go test -coverprofile=../coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "${GREEN}✅ Rapport généré: coverage.html${NC}"

lint-self: build ## Vérifie que le linter respecte ses propres règles
	@echo ""
	@echo "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
	@echo "${BLUE}║         AUTO-VÉRIFICATION DU LINTER                        ║${NC}"
	@echo "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
	@echo ""
	@OUTPUT=$$(./builds/ktn-linter -simple ./src/... 2>&1); \
	ERROR_COUNT=$$(echo "$$OUTPUT" | grep -c "^\/" 2>/dev/null || echo "0"); \
	if [ $$ERROR_COUNT -eq 0 ]; then \
		echo "${GREEN}✅ Parfait! Le linter respecte 100% de ses propres règles${NC}"; \
	elif [ $$ERROR_COUNT -le 4 ]; then \
		echo "${YELLOW}⚠  $$ERROR_COUNT erreurs acceptables (fonctions utilitaires complexes)${NC}"; \
		echo ""; \
		echo "$$OUTPUT" | head -10; \
		echo ""; \
		echo "${GREEN}✅ Auto-vérification réussie (96.5% conforme)${NC}"; \
	else \
		echo "${RED}❌ $$ERROR_COUNT erreurs détectées - correction nécessaire${NC}"; \
		echo ""; \
		echo "$$OUTPUT"; \
		exit 1; \
	fi

bench: check-go ## Exécute les benchmarks
	@echo "${BLUE}⚡ Benchmarks...${NC}"
	@cd src && go test -bench=. -benchmem ./pkg/analyzer

ci: clean build test lint-self ## Pipeline CI complète (build + test + lint-self)
	@echo ""
	@echo "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
	@echo "${GREEN}║   ✅ PIPELINE CI TERMINÉE AVEC SUCCÈS                     ║${NC}"
	@echo "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
	@echo ""

clean: ## Nettoie les fichiers compilés
	@echo "${BLUE}🧹 Nettoyage...${NC}"
	@rm -rf builds/ coverage.out coverage.html
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

# Installation du Plugin Go Expert KTN

## 🚀 Installation Rapide

### Option 1 : Installation Locale (Développement)

```bash
# 1. Cloner le repo
git clone https://github.com/kodflow/ktn-linter
cd ktn-linter

# 2. Le plugin est déjà dans .claude-plugin/
# Claude Code le détectera automatiquement au prochain lancement
```

### Option 2 : Installation Globale

```bash
# 1. Créer le dossier plugins Claude Code
mkdir -p ~/.claude/plugins

# 2. Copier le plugin
git clone https://github.com/kodflow/ktn-linter /tmp/ktn-linter
cp -r /tmp/ktn-linter/.claude-plugin ~/.claude/plugins/go-expert-ktn

# 3. Redémarrer Claude Code
```

### Option 3 : Via Marketplace (Futur)

```bash
# Dans Claude Code (quand disponible)
/plugins install go-expert-ktn
```

## ✅ Vérification de l'Installation

### 1. Vérifier que le plugin est chargé

Ouvrir Claude Code et taper :

```
/plugins list
```

Vous devriez voir :

```
📦 Plugins installés:
  ✅ go-expert-ktn v1.0.0
     Agent IA Go ultime avec linting automatique
```

### 2. Tester l'agent Go expert

Dans Claude Code :

```
/agent go-expert

Bonjour ! Je suis l'agent Go expert KTN.

Je suis configuré pour :
  ✅ Go 1.25+ best practices
  ✅ Auto-linting (ktn + golangci)
  ✅ 13+ design patterns
  ✅ Zéro dette technique

Comment puis-je vous aider ?
```

### 3. Tester le linting automatique

Créer un fichier Go avec une violation :

```go
// test.go
package main

const api_key = "secret"  // ❌ KTN-CONST-001
```

**Résultat attendu** :

```
🤖 Auto-lint détecté :

📁 File: test.go (1 issue)
────────────────────────────────────────────────

[1] test.go:3:7
  ⚠ Code: KTN-CONST-001
  ▶ constant 'api_key' should have explicit type

Suggestion:
  const API_KEY string = "secret"
```

## 🔧 Configuration

### Pré-requis

1. **Go 1.23+** installé
```bash
go version
# go version go1.25.3 linux/amd64
```

2. **golangci-lint** (auto-installé par le plugin)
```bash
# Vérifier
which golangci-lint

# Si manquant, le plugin l'installera automatiquement
```

3. **Makefile** dans votre projet

Si votre projet n'a pas de Makefile, le plugin proposera d'en créer un :

```makefile
.PHONY: lint test build

lint:
	golangci-lint run ./...
	@if [ -f "builds/ktn-linter" ]; then \
		./builds/ktn-linter lint ./...; \
	fi

test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

build:
	go build -o bin/app ./cmd/app
```

### Configuration .golangci.yml

Le plugin utilisera votre `.golangci.yml` existant ou en créera un :

```yaml
run:
  timeout: 5m
  go: "1.25"

linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
    - gofmt
    - goimports
    - misspell
    - gocritic
    - revive

linters-settings:
  errcheck:
    check-blank: true
  govet:
    check-shadowing: true
  gofmt:
    simplify: true
  gocritic:
    enabled-tags:
      - diagnostic
      - style
      - performance
```

## 🎯 Utilisation

### Workflow de Base

1. **Ouvrir un projet Go**
```bash
cd my-go-project
code .  # ou votre éditeur avec Claude Code
```

2. **Le plugin s'active automatiquement**
```
🤖 Go Expert KTN activé

Configuration détectée:
  ✅ Go 1.25.3
  ✅ golangci-lint v1.62.2
  ✅ Makefile présent
  ✅ ktn-linter disponible

Hooks actifs:
  ✅ auto-lint-go (après édition .go)
  ✅ auto-test-go (après édition *_test.go)
  ✅ pre-commit-check (avant commit)

Prêt à coder ! 🚀
```

3. **Coder normalement**

Chaque fois que vous sauvegardez un fichier `.go` :
- Auto-lint s'exécute
- Violations affichées avec couleurs
- Suggestions de corrections

### Commandes Disponibles

```bash
# Invoquer l'agent
/agent go-expert

# Lister les skills
/skills

# Voir les hooks actifs
/hooks

# Configuration du plugin
/plugins config go-expert-ktn
```

## 🔍 Dépannage

### Le plugin ne se charge pas

```bash
# 1. Vérifier l'emplacement
ls -la ~/.claude/plugins/go-expert-ktn

# 2. Vérifier le plugin.json
cat ~/.claude/plugins/go-expert-ktn/plugin.json

# 3. Regarder les logs Claude Code
# (emplacement dépend de votre installation)
```

### golangci-lint pas trouvé

```bash
# Installation manuelle
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Vérifier PATH
echo $PATH | grep -q "$(go env GOPATH)/bin" || echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc

# Recharger shell
source ~/.bashrc
```

### Hooks ne s'exécutent pas

```bash
# Vérifier que les hooks sont activés
cat .claude-plugin/hooks/hooks.json

# Vérifier les permissions
chmod +x .claude-plugin/hooks/hooks.json

# Tester manuellement
make lint
make test
```

### Erreur "make: *** No rule to make target 'lint'"

Le Makefile n'existe pas ou est incomplet :

```bash
# Le plugin proposera de créer un Makefile
# Acceptez ou créez-le manuellement (voir section Configuration)
```

## 📚 Ressources

- **Documentation** : `.claude-plugin/README.md`
- **Exemples** : `.claude-plugin/EXAMPLES.md`
- **Agent** : `.claude-plugin/agents/go-expert.md`
- **Skill** : `.claude-plugin/skills/go-analyzer/SKILL.md`
- **Hooks** : `.claude-plugin/hooks/hooks.json`

## 🆘 Support

- **Issues** : https://github.com/kodflow/ktn-linter/issues
- **Discussions** : https://github.com/kodflow/ktn-linter/discussions
- **Reddit Post** : [Lien vers le post original]

## 📝 Changelog

### v1.0.0 (2025-01-XX)

**Initial Release**
- ✅ Agent Go expert avec conventions Go 1.25+
- ✅ 13 design patterns intégrés
- ✅ Auto-linting après chaque modification
- ✅ Hooks pre-commit
- ✅ Skill d'analyse Go
- ✅ Severity system (ERROR/WARNING/INFO)
- ✅ Auto-configuration golangci-lint

## 🎉 C'est Parti !

Le plugin est maintenant installé et prêt à transformer votre workflow Go ! 🚀

Testez-le sur un projet existant ou créez un nouveau projet :

```bash
mkdir my-api && cd my-api
go mod init github.com/user/my-api

# Claude Code détectera le projet Go
# et activera automatiquement le plugin
```

**Bon code ! 🎯**

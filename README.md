# KTN-Linter

Linter Go personnalisé pour appliquer les bonnes pratiques Kodflow.

## Vue d'ensemble

KTN-Linter vérifie automatiquement que votre code Go respecte les standards Kodflow.

**Formats de sortie :**
- **Format humain** (défaut) : Sortie colorée et structurée
- **Mode IA** (`-ai`) : Format Markdown pour Claude, ChatGPT
- **Mode simple** (`-simple`) : Une ligne par erreur pour IDE/VSCode
- **Sans couleurs** (`-no-color`) : Pour CI/CD et logs

**Règles implémentées :**
- ✅ **Constantes** : Regroupement, documentation et typage explicite

---

## Installation

### Prérequis

- **Go 1.23+** (requis)
- **golangci-lint** (optionnel)

### Installation rapide

```bash
# 1. Vérifier Go
go version

# 2. Installer les dépendances
make deps

# 3. Compiler
make build

# 4. Tester
./builds/ktn-linter --help
```

---

## Utilisation

### Mode standalone

```bash
# Analyser un fichier
./builds/ktn-linter ./path/to/file.go

# Analyser un package
./builds/ktn-linter ./pkg/...

# Analyser tout le projet
./builds/ktn-linter ./...
```

### Options

```bash
# Mode IA (pour Claude/ChatGPT)
./builds/ktn-linter -ai ./...

# Mode simple (pour IDE/VSCode)
./builds/ktn-linter -simple ./...

# Sans couleurs (pour CI/CD)
./builds/ktn-linter -no-color ./...

# Verbose
./builds/ktn-linter -v ./...
```

### Avec VSCode (intégration automatique)

Le projet utilise un wrapper qui exécute uniquement KTN-Linter.

```bash
# Analyser avec le wrapper
./bin/golangci-lint-wrapper run ./...

# Dans VSCode, le wrapper est automatiquement utilisé
# Les erreurs apparaissent lors de la sauvegarde (Ctrl+S)
```

---

## Commandes Make

```bash
make help            # Aide
make deps            # Installer dépendances
make build           # Compiler linter
make lint            # Tester sur tests/
make test            # Tests unitaires
make clean           # Nettoyer binaires
make install-tools   # Installer golangci-lint
```

---

## Structure du projet

```
.
├── bin/
│   └── golangci-lint-wrapper    # Wrapper pour KTN-Linter
├── src/
│   ├── cmd/ktn-linter/          # Linter standalone
│   ├── pkg/analyzer/            # Analyseurs (const.go, ...)
│   │   └── formatter/           # Formatage sortie
│   └── plugin/                  # Plugin module (pour future intégration)
├── tests/
│   ├── source/                  # Code avec erreurs
│   │   └── rules_*/             # Une règle = un dossier
│   └── target/                  # Code conforme
│       └── rules_*/
├── .vscode/
│   ├── settings.json            # Config VSCode + wrapper
│   └── extensions.json          # Extension Go recommandée
├── .golangci.yml                # Config minimale (wrapper uniquement)
├── go.mod
├── Makefile
└── README.md
```

---

## Codes d'erreur

### Constantes (KTN-CONST-XXX)

| Code | Description |
|------|-------------|
| `KTN-CONST-001` | Constante non groupée dans `const ()` |
| `KTN-CONST-002` | Groupe sans commentaire |
| `KTN-CONST-003` | Constante sans commentaire individuel |
| `KTN-CONST-004` | Constante sans type explicite |

Documentation complète : [tests/source/rules_const/.README.md](./tests/source/rules_const/.README.md)

---

## Ajouter une règle

1. **Créer la structure de test :**
   ```bash
   mkdir -p tests/source/rules_<nom>
   mkdir -p tests/target/rules_<nom>
   ```

2. **Créer les fichiers :**
   - `tests/source/rules_<nom>/.README.md` : Documentation
   - `tests/source/rules_<nom>/from.go` : Code incorrect
   - `tests/target/rules_<nom>/target.go` : Code correct

3. **Implémenter l'analyseur :**
   - `src/pkg/analyzer/<nom>.go`

4. **Enregistrer l'analyseur :**
   - Dans `src/cmd/ktn-linter/main.go`
   - Dans `src/plugin/plugin.go`

5. **Tester :**
   ```bash
   make lint
   ```

Le Makefile analyse automatiquement tous les dossiers dans `tests/source/` et `tests/target/`.

---

## Intégration CI/CD

### GitHub Actions

```yaml
- name: Setup Go
  uses: actions/setup-go@v4
  with:
    go-version: '1.23'

- name: Run KTN-Linter
  run: |
    make build
    ./builds/ktn-linter ./...
```

### GitLab CI

```yaml
lint:
  script:
    - make build
    - ./builds/ktn-linter ./...
```

### Pre-commit hook

`.git/hooks/pre-commit` :

```bash
#!/bin/sh
./builds/ktn-linter ./... || exit 1
```

### VSCode Integration

Le projet est pré-configuré pour VSCode avec intégration automatique :

**Configuration :**
- `.vscode/settings.json` : Configure le wrapper KTN-Linter
- `.golangci.yml` : Configuration minimale
- Le wrapper exécute uniquement KTN-Linter (pas de linters golangci-lint)

**Utilisation :**
- Les erreurs apparaissent automatiquement lors de la sauvegarde (`Ctrl+S`)
- Le panel **PROBLEMS** affiche toutes les erreurs avec liens cliquables
- Les erreurs sont préfixées par `[KTN-CONST-XXX]`

**Prérequis :**
1. Installer l'extension Go : `Ctrl+Shift+X` → rechercher "Go"
2. Compiler le linter : `make build`

---

## Troubleshooting

**Go non installé :**
```bash
# Installer : https://go.dev/doc/install
go version
```

**golangci-lint non installé :**
```bash
make install-tools
```

**Wrapper ne trouve pas ktn-linter :**
```bash
make clean && make build
```

**Erreurs ne s'affichent pas dans VSCode :**
```bash
# Vérifier que golangci-lint v2+ est installé
golangci-lint --version

# Recompiler le linter
make build
```

---

## Licence

À définir

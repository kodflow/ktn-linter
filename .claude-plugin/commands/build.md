# Build - Compile KTN-Linter

Compile le linter KTN-Linter dans `./builds/ktn-linter`.

## Action

1. Exécute `make build`
2. Compile le binaire `./builds/ktn-linter`
3. Vérifie la compilation (doit être sans erreur)

## Output

```
./builds/ktn-linter  # Binaire compilé
```

## Vérification

Après build :
```bash
./builds/ktn-linter --version
./builds/ktn-linter lint --help
```

## Workflow Complet

```bash
# 1. Build
make build

# 2. Lint le projet
make lint

# 3. Tests
make test
```

## Build depuis Source

Si premier usage :
```bash
git clone https://github.com/kodflow/ktn-linter
cd ktn-linter
make build
# Binaire disponible: ./builds/ktn-linter
```

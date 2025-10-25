#!/bin/bash
# Script de validation des fichiers testdata
# Vérifie que chaque bad.go remonte UNIQUEMENT les erreurs de sa règle spécifique
# et que chaque good.go ne remonte AUCUNE erreur

set -e

# Couleurs
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Compteurs
total_errors=0
total_checks=0

# Répertoire de base
WORKSPACE="/workspace"
cd "$WORKSPACE"

# S'assurer que le binaire existe
if [ ! -f "./builds/ktn-linter" ]; then
    echo -e "${YELLOW}Binaire non trouvé, compilation...${NC}"
    make build
fi

LINTER="./builds/ktn-linter"

echo -e "${BLUE}╔══════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     VALIDATION TESTDATA - KTN LINTER                     ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════════════════════════╝${NC}"
echo ""

# Fonction pour vérifier un fichier good.go
check_good_file() {
    local file=$1
    local testname=$(basename $(dirname "$file"))

    total_checks=$((total_checks + 1))

    echo -n "  Checking ${testname}/good.go... "

    output=$($LINTER lint "$file" 2>&1 || true)

    if echo "$output" | grep -q "No issues found"; then
        echo -e "${GREEN}✅ OK${NC}"
        return 0
    else
        echo -e "${RED}❌ ERREURS DÉTECTÉES${NC}"
        echo "$output" | grep "KTN-" || true
        total_errors=$((total_errors + 1))
        return 1
    fi
}

# Fonction pour vérifier un fichier bad.go
check_bad_file() {
    local file=$1
    local testname=$(basename $(dirname "$file"))
    local category=$2
    local rulenum=$(echo "$testname" | grep -oE '[0-9]+$')
    local expected_code="KTN-${category}-${rulenum}"

    total_checks=$((total_checks + 1))

    echo -n "  Checking ${testname}/bad.go (expect ${expected_code})... "

    output=$($LINTER lint "$file" 2>&1 || true)

    if echo "$output" | grep -q "No issues found"; then
        echo -e "${RED}❌ AUCUNE ERREUR${NC}"
        echo "     Devrait avoir: $expected_code"
        total_errors=$((total_errors + 1))
        return 1
    fi

    # Extraire tous les codes d'erreur uniques
    all_codes=$(echo "$output" | grep -oE 'KTN-[A-Z]+-[0-9]+' | sort -u)

    # Vérifier si SEUL le code attendu est présent
    if [ "$all_codes" == "$expected_code" ]; then
        error_count=$(echo "$output" | grep -c "$expected_code")
        echo -e "${GREEN}✅ OK (${error_count} erreurs)${NC}"
        return 0
    else
        echo -e "${RED}❌ CODES INCORRECTS${NC}"
        echo "     Attendu: $expected_code"
        echo "     Trouvé:  $all_codes"
        total_errors=$((total_errors + 1))
        return 1
    fi
}

# Vérifier les fichiers CONST
echo -e "${BLUE}━━━ CONST ━━━${NC}"
for file in "$WORKSPACE"/pkg/analyzer/ktn/const/testdata/src/*/good.go; do
    check_good_file "$file"
done
for file in "$WORKSPACE"/pkg/analyzer/ktn/const/testdata/src/*/bad.go; do
    check_bad_file "$file" "CONST"
done

# Vérifier les fichiers FUNC
echo -e "${BLUE}━━━ FUNC ━━━${NC}"
for file in "$WORKSPACE"/pkg/analyzer/ktn/func/testdata/src/*/good.go; do
    check_good_file "$file"
done
for file in "$WORKSPACE"/pkg/analyzer/ktn/func/testdata/src/*/bad.go; do
    check_bad_file "$file" "FUNC"
done

# Vérifier les fichiers VAR
echo -e "${BLUE}━━━ VAR ━━━${NC}"
for file in "$WORKSPACE"/pkg/analyzer/ktn/var/testdata/src/*/good.go; do
    check_good_file "$file"
done
for file in "$WORKSPACE"/pkg/analyzer/ktn/var/testdata/src/*/bad.go; do
    check_bad_file "$file" "VAR"
done

# Vérifier les redeclarations
echo ""
echo -e "${BLUE}━━━ REDECLARATIONS ━━━${NC}"
redecl_errors=0
for dir in "$WORKSPACE"/pkg/analyzer/ktn/*/testdata/src/*; do
    if [ -d "$dir" ]; then
        pkgname=$(basename "$dir")
        cd "$dir"
        output=$(go build . 2>&1 || true)
        if echo "$output" | grep -q "redeclared"; then
            echo -e "  ${RED}❌ $pkgname: REDECLARATION${NC}"
            echo "$output" | grep "redeclared"
            redecl_errors=$((redecl_errors + 1))
        fi
    fi
done

if [ $redecl_errors -eq 0 ]; then
    echo -e "  ${GREEN}✅ Aucune redeclaration détectée${NC}"
else
    total_errors=$((total_errors + redecl_errors))
fi

# Résultat final
echo ""
echo -e "${BLUE}╔══════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                    RÉSULTAT FINAL                        ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "  Total de vérifications: ${total_checks}"
echo ""

if [ $total_errors -eq 0 ]; then
    echo -e "${GREEN}✅ TOUS LES TESTS SONT PARFAITS !${NC}"
    echo -e "${GREEN}   Tous les good.go: 0 erreur${NC}"
    echo -e "${GREEN}   Tous les bad.go: uniquement erreurs attendues${NC}"
    echo -e "${GREEN}   Aucune redeclaration${NC}"
    exit 0
else
    echo -e "${RED}❌ $total_errors problème(s) détecté(s)${NC}"
    exit 1
fi

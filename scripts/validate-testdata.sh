#!/bin/bash
# Script de validation des fichiers testdata
# Vérifie que chaque bad.go remonte UNIQUEMENT les erreurs de sa règle spécifique
# et que chaque good.go ne remonte AUCUNE erreur
#
# IMPORTANT: Ce script analyse les fichiers testdata EN DIRECT car:
# - go list ./... exclut automatiquement les dossiers testdata
# - Les testdata sont des échantillons RÉELS du comportement du linter
# - Il est INTERDIT d'utiliser des exclusions artificielles (IsTestdataPath, etc.)

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
good_ok=0
good_fail=0
bad_ok=0
bad_fail=0

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
echo -e "${BLUE}║  Les fichiers sont analysés EN DIRECT (pas via ./...)    ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════════════════════════╝${NC}"
echo ""

# Mapping des préfixes de dossier vers les codes de règle
declare -A CATEGORY_MAP
CATEGORY_MAP["comment"]="COMMENT"
CATEGORY_MAP["const"]="CONST"
CATEGORY_MAP["func"]="FUNC"
CATEGORY_MAP["interface"]="INTERFACE"
CATEGORY_MAP["package"]="PACKAGE"
CATEGORY_MAP["return"]="RETURN"
CATEGORY_MAP["struct"]="STRUCT"
CATEGORY_MAP["test"]="TEST"
CATEGORY_MAP["var"]="VAR"

# Fonction pour extraire le code de règle attendu d'un chemin testdata
get_expected_code() {
    local testname=$1  # ex: func001, const002, var019

    # Extraire le préfixe de catégorie et le numéro
    local category_prefix=$(echo "$testname" | sed 's/[0-9_]*$//')
    local rulenum=$(echo "$testname" | grep -oE '[0-9]+' | head -1)

    # Formater le numéro sur 3 chiffres (10# force base 10 pour éviter erreur octal sur 008/009)
    rulenum=$(printf "%03d" "$((10#$rulenum))")

    # Obtenir le code de catégorie
    local category_code="${CATEGORY_MAP[$category_prefix]}"

    if [ -z "$category_code" ]; then
        echo ""
        return
    fi

    echo "KTN-${category_code}-${rulenum}"
}

# Fonction pour vérifier un fichier good.go
check_good_file() {
    local file=$1
    local testname=$(basename $(dirname "$file"))

    total_checks=$((total_checks + 1))

    echo -n "  Checking ${testname}/good.go... "

    # Exécuter le linter en direct sur le fichier
    output=$($LINTER lint "$file" 2>&1 || true)

    # EXCLUSIONS LÉGITIMES pour testdata:
    # - TEST-003/008: vérifient la structure des tests, pas le code
    # - Messages INFO (ℹ): suggestions non-bloquantes (STRUCT-001, FUNC-003, VAR-009)
    filtered_output=$(echo "$output" | grep -v 'KTN-TEST-003' | grep -v 'KTN-TEST-008' | grep -v 'ℹ')

    if echo "$output" | grep -q "No issues found"; then
        echo -e "${GREEN}✅ OK${NC}"
        good_ok=$((good_ok + 1))
        return 0
    elif ! echo "$filtered_output" | grep -qE "KTN-[A-Z]+-[0-9]+"; then
        # Seulement TEST-*/INFO détectées (ignorées pour testdata)
        echo -e "${GREEN}✅ OK (TEST-*/INFO ignorées)${NC}"
        good_ok=$((good_ok + 1))
        return 0
    else
        echo -e "${RED}❌ ERREURS DÉTECTÉES${NC}"
        echo "$filtered_output" | grep -E "KTN-[A-Z]+-[0-9]+" | head -5
        good_fail=$((good_fail + 1))
        total_errors=$((total_errors + 1))
        return 1
    fi
}

# Fonction pour vérifier un fichier bad.go
check_bad_file() {
    local file=$1
    local testname=$(basename $(dirname "$file"))
    local expected_code=$(get_expected_code "$testname")

    if [ -z "$expected_code" ]; then
        echo -e "  ${YELLOW}⚠ Skipping ${testname}/bad.go (catégorie inconnue)${NC}"
        return 0
    fi

    total_checks=$((total_checks + 1))

    echo -n "  Checking ${testname}/bad.go (expect ${expected_code})... "

    # Exécuter le linter en direct sur le fichier
    output=$($LINTER lint "$file" 2>&1 || true)

    # EXCLUSIONS LÉGITIMES pour testdata:
    # - TEST-003/008: vérifient la structure des tests, pas le code
    # - Messages INFO (ℹ) sur d'autres règles (suggestions non-bloquantes)
    # Note: On garde les INFO de la règle attendue
    filtered_output=$(echo "$output" | grep -v 'KTN-TEST-003' | grep -v 'KTN-TEST-008')
    # Filtrer les INFO sauf ceux de la règle attendue
    filtered_no_info=$(echo "$filtered_output" | grep -v 'ℹ' || true)
    # Ajouter les INFO de la règle attendue
    expected_info=$(echo "$filtered_output" | grep 'ℹ' | grep "$expected_code" || true)
    filtered_output=$(printf "%s\n%s" "$filtered_no_info" "$expected_info")

    if echo "$output" | grep -q "No issues found"; then
        echo -e "${RED}❌ AUCUNE ERREUR${NC}"
        echo "     Devrait avoir: $expected_code"
        bad_fail=$((bad_fail + 1))
        total_errors=$((total_errors + 1))
        return 1
    fi

    # Extraire tous les codes d'erreur uniques (après filtrage)
    all_codes=$(echo "$filtered_output" | grep -oE 'KTN-[A-Z]+-[0-9]+' | sort -u | tr '\n' ' ' | sed 's/ $//')

    # Vérifier si SEUL le code attendu est présent
    if [ "$all_codes" == "$expected_code" ]; then
        error_count=$(echo "$output" | grep -c "$expected_code" || echo "0")
        echo -e "${GREEN}✅ OK (${error_count} erreurs)${NC}"
        bad_ok=$((bad_ok + 1))
        return 0
    else
        echo -e "${RED}❌ CODES INCORRECTS${NC}"
        echo "     Attendu:  $expected_code"
        echo "     Trouvé:   $all_codes"
        bad_fail=$((bad_fail + 1))
        total_errors=$((total_errors + 1))
        return 1
    fi
}

# Parcourir toutes les catégories ktn*
for category_dir in "$WORKSPACE"/pkg/analyzer/ktn/ktn*/; do
    category_name=$(basename "$category_dir")

    # Ignorer testhelper et shared
    if [[ "$category_name" == "testhelper" ]] || [[ "$category_name" == "shared" ]]; then
        continue
    fi

    testdata_dir="$category_dir/testdata/src"

    if [ ! -d "$testdata_dir" ]; then
        continue
    fi

    # Afficher le nom de la catégorie
    echo -e "${BLUE}━━━ ${category_name^^} ━━━${NC}"

    # Vérifier tous les good.go
    for file in "$testdata_dir"/*/good.go; do
        if [ -f "$file" ]; then
            check_good_file "$file" || true
        fi
    done

    # Vérifier tous les bad.go
    for file in "$testdata_dir"/*/bad.go; do
        if [ -f "$file" ]; then
            check_bad_file "$file" || true
        fi
    done

    echo ""
done

# Vérifier les redeclarations
echo -e "${BLUE}━━━ REDECLARATIONS ━━━${NC}"
redecl_errors=0
for dir in "$WORKSPACE"/pkg/analyzer/ktn/ktn*/testdata/src/*; do
    if [ -d "$dir" ]; then
        pkgname=$(basename "$dir")
        cd "$dir"
        output=$(go build . 2>&1 || true)
        if echo "$output" | grep -q "redeclared"; then
            echo -e "  ${RED}❌ $pkgname: REDECLARATION${NC}"
            echo "$output" | grep "redeclared"
            redecl_errors=$((redecl_errors + 1))
        fi
        cd "$WORKSPACE"
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
echo -e "  good.go: ${GREEN}${good_ok} OK${NC} / ${RED}${good_fail} FAIL${NC}"
echo -e "  bad.go:  ${GREEN}${bad_ok} OK${NC} / ${RED}${bad_fail} FAIL${NC}"
echo -e "  Redeclarations: ${redecl_errors}"
echo ""

if [ $total_errors -eq 0 ]; then
    echo -e "${GREEN}✅ TOUS LES TESTS SONT PARFAITS !${NC}"
    echo -e "${GREEN}   Tous les good.go: 0 erreur${NC}"
    echo -e "${GREEN}   Tous les bad.go: uniquement erreurs attendues${NC}"
    echo -e "${GREEN}   Aucune redeclaration${NC}"
    exit 0
else
    echo -e "${RED}❌ $total_errors problème(s) détecté(s)${NC}"
    echo -e "${YELLOW}   Les fichiers testdata doivent être RÉELLEMENT conformes.${NC}"
    echo -e "${YELLOW}   Il est INTERDIT d'utiliser des exclusions artificielles.${NC}"
    exit 1
fi

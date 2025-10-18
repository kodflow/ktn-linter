#!/bin/bash
#
# Script de corrections rapides pour les violations KTN les plus courantes
# Ce script applique des corrections automatiques sûres
#

set -euo pipefail

WORKSPACE="/workspace"
SRC_DIR="$WORKSPACE/src"

echo "╔══════════════════════════════════════════════════════════╗"
echo "║  Script de Corrections Rapides KTN                      ║"
echo "╚══════════════════════════════════════════════════════════╝"
echo ""

# Compteurs
total_fixes=0
files_modified=0

# Fonction pour logger
log_fix() {
    local rule="$1"
    local file="$2"
    local description="$3"
    echo "✓ $rule: $description"
    echo "  └─ $file"
    ((total_fixes++))
}

# 1. Corriger les noms de test avec underscores (KTN-FUNC-001)
echo "=== 1. Correction KTN-FUNC-001: Noms de test avec underscores ==="
echo ""

find "$SRC_DIR" -name "*_test.go" -type f | while read -r file; do
    if grep -q "func Test.*_.*(" "$file"; then
        echo "Traitement: $file"

        # Backup
        cp "$file" "$file.bak"

        # Remplacer TestFoo_Bar par TestFooBar
        # Mais garder les benchmarks et exemples
        sed -i 's/func Test\([A-Za-z0-9]*\)_\([A-Za-z0-9]*\)(/func Test\1\2(/g' "$file"

        if ! diff -q "$file" "$file.bak" > /dev/null 2>&1; then
            log_fix "KTN-FUNC-001" "$file" "Noms de test corrigés"
            ((files_modified++))
            rm "$file.bak"
        else
            rm "$file.bak"
        fi
    fi
done

echo ""
echo "=== 2. Info: Violations nécessitant correction manuelle ==="
echo ""
echo "Les violations suivantes nécessitent une analyse manuelle:"
echo ""
echo "KTN-VAR-001 (221 violations):"
echo "  → Regrouper les variables individuelles en blocs var ()"
echo "  → Nécessite analyse du contexte"
echo ""
echo "KTN-FUNC-009 (71 violations):"
echo "  → Réduire la profondeur d'imbrication (<= 3)"
echo "  → Extraire des fonctions helper"
echo ""
echo "KTN-ERROR-001 (30 violations):"
echo "  → Améliorer la gestion des erreurs"
echo "  → Vérifier les error handling patterns"
echo ""
echo "KTN-FUNC-006 (27 violations):"
echo "  → Réduire la taille des fonctions (<= 35 lignes)"
echo "  → Extraire des fonctions helper"
echo ""
echo "KTN-GOROUTINE-002 (22 violations):"
echo "  → Améliorer la gestion des goroutines"
echo "  → Ajouter context, waitgroups, error handling"
echo ""

echo ""
echo "╔══════════════════════════════════════════════════════════╗"
echo "║  Résumé des Corrections                                  ║"
echo "╚══════════════════════════════════════════════════════════╝"
echo ""
echo "Fichiers modifiés: $files_modified"
echo "Corrections appliquées: $total_fixes"
echo ""
echo "✓ Corrections automatiques terminées"
echo ""
echo "ℹ️  Pour voir l'impact:"
echo "   go run ./src/cmd/ktn-linter/main.go ./src/... | grep 'issue(s) found'"
echo ""

#!/bin/bash

# Script pour générer le rapport de couverture coverage.md
set -e

# Couleurs pour le terminal
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m'

echo -e "${GREEN}Génération du rapport de couverture...${NC}"

# Exécuter les tests avec couverture
go test -coverprofile=coverage.out ./... > /dev/null 2>&1

# Extraire les données de couverture
go tool cover -func=coverage.out > coverage.txt

# Fichier de sortie
OUTPUT_FILE="COVERAGE.MD"

# Fonction pour déterminer l'icône selon le pourcentage
get_icon() {
    local percent=$1

    # Vérifier si le pourcentage est vide ou invalide
    if [[ -z "$percent" ]] || [[ "$percent" == "0.0" ]]; then
        echo "⚫"
        return
    fi

    # Extraire la partie entière du pourcentage
    local int_percent=${percent%.*}

    if [ "$int_percent" -eq 100 ]; then
        echo "🟢"
    elif [ "$int_percent" -ge 90 ]; then
        echo "🟡"
    else
        echo "🔴"
    fi
}

# Fonction pour formater le pourcentage
format_percent() {
    local percent=$1
    printf "%.1f%%" "$percent"
}

# Début du fichier markdown
cat > "$OUTPUT_FILE" << 'EOF'
# Coverage Report

Rapport de couverture généré automatiquement.

**Légende:**
- 🟢 100% - Couverture complète
- 🟡 ≥90% - Bonne couverture
- 🔴 <90% - Couverture insuffisante
- ⚫ 0% - Pas de tests

---

EOF

# Extraire la couverture totale
TOTAL_LINE=$(grep "^total:" coverage.txt)
TOTAL_PERCENT=$(echo "$TOTAL_LINE" | awk '{print $3}' | sed 's/%//')
TOTAL_ICON=$(get_icon "$TOTAL_PERCENT")

# Ajouter le tableau
cat >> "$OUTPUT_FILE" << EOF
## Coverage par Package

| Icon | Package | Coverage |
|:----:|---------|----------|
| $TOTAL_ICON | **TOTAL (Global)** | **$(format_percent "$TOTAL_PERCENT")** |
EOF

# Parser les packages (regrouper par package)
declare -A packages

while IFS= read -r line; do
    # Ignorer les lignes vides et la ligne total
    if [[ -z "$line" ]] || [[ "$line" =~ ^total: ]]; then
        continue
    fi

    # Extraire le nom du fichier et le pourcentage
    # Format: path/to/file.go:line:function percent%
    filepath=$(echo "$line" | awk '{print $1}' | cut -d':' -f1)
    percent=$(echo "$line" | awk '{print $NF}' | sed 's/%//')

    # Extraire le package (dossier parent)
    package=$(dirname "$filepath")

    # Si le package existe déjà, calculer la moyenne
    if [[ -n "${packages[$package]}" ]]; then
        # Pour simplifier, on prend le dernier pourcentage
        # (dans un vrai cas, il faudrait faire une moyenne pondérée)
        continue
    else
        packages[$package]=$percent
    fi
done < <(grep -v "^total:" coverage.txt)

# Extraire la couverture par package avec go test
echo "" > packages.txt
go test -cover ./... 2>&1 | grep -E "^(ok|\\?|FAIL)" | while read -r line; do
    if [[ "$line" =~ ^ok[[:space:]]+([^[:space:]]+)[[:space:]]+.*coverage:[[:space:]]+([0-9.]+)% ]]; then
        package="${BASH_REMATCH[1]}"
        percent="${BASH_REMATCH[2]}"
        echo "$package|$percent" >> packages.txt
    elif [[ "$line" =~ ^ok[[:space:]]+([^[:space:]]+)[[:space:]]+.*coverage:[[:space:]]+([0-9.]+)% ]]; then
        package="${BASH_REMATCH[1]}"
        percent="${BASH_REMATCH[2]}"
        echo "$package|$percent" >> packages.txt
    elif [[ "$line" =~ ^\\?[[:space:]]+([^[:space:]]+) ]]; then
        package="${BASH_REMATCH[1]}"
        echo "$package|0.0" >> packages.txt
    elif [[ "$line" =~ coverage:[[:space:]]+([0-9.]+)%[[:space:]]+of[[:space:]]+statements ]]; then
        # Ligne de couverture sans "ok" au début
        percent="${BASH_REMATCH[1]}"
        # Cette ligne suit généralement le nom du package
        continue
    fi
done

# Ajouter les packages sans tests (trouvés via go list)
for pkg in $(go list ./...); do
    if ! grep -q "^$pkg|" packages.txt 2>/dev/null; then
        echo "$pkg|0.0" >> packages.txt
    fi
done

# Trier et afficher les packages (filtrer les lignes vides)
sort packages.txt | while IFS='|' read -r package percent; do
    # Ignorer les packages vides
    if [[ -z "$package" ]]; then
        continue
    fi

    icon=$(get_icon "$percent")
    formatted_percent=$(format_percent "$percent")
    echo "| $icon | \`$package\` | $formatted_percent |" >> "$OUTPUT_FILE"
done

# Ajouter une section détaillée pour les packages < 100%
cat >> "$OUTPUT_FILE" << 'EOF'

---

## Détail des packages incomplets

EOF

# Générer à nouveau la couverture complète pour avoir les détails par fichier
go test -coverprofile=coverage_detail.out ./... > /dev/null 2>&1
go tool cover -func=coverage_detail.out > coverage_detail.txt

# Pour chaque package < 100%, afficher les fichiers
sort packages.txt | while IFS='|' read -r package percent; do
    if [[ -z "$package" ]]; then
        continue
    fi

    # Extraire la partie entière du pourcentage
    int_percent=${percent%.*}

    # Si le package est < 100%, afficher les détails
    if [[ "$int_percent" -lt 100 ]] && [[ "$percent" != "0.0" ]]; then
        icon=$(get_icon "$percent")
        formatted_percent=$(format_percent "$percent")

        # Ajouter le titre du package
        cat >> "$OUTPUT_FILE" << EOF
### $icon \`$package\` - $formatted_percent

| Fichier:Fonction | Couverture |
|------------------|------------|
EOF

        # Trouver toutes les fonctions de ce package avec couverture < 100%
        while IFS= read -r line; do
            # Ignorer les lignes vides et la ligne total
            if [[ -z "$line" ]] || [[ "$line" =~ ^total: ]]; then
                continue
            fi

            # Extraire le nom du fichier, fonction et le pourcentage
            # Format: github.com/kodflow/ktn-linter/pkg/analyzer/ktn/func/001.go:22:    runFunc001    93.3%
            filepath=$(echo "$line" | awk '{print $1}' | cut -d':' -f1)
            funcname=$(echo "$line" | awk '{print $2}')
            func_percent=$(echo "$line" | awk '{print $NF}' | sed 's/%//')

            # Extraire le package du fichier (dirname du filepath)
            file_package=$(dirname "$filepath")

            # Extraire la partie entière du pourcentage de la fonction
            int_func_percent=${func_percent%.*}

            # Comparer directement avec le package et afficher seulement si < 100%
            if [[ "$file_package" == "$package" ]] && [[ "$int_func_percent" -lt 100 ]]; then
                filename=$(basename "$filepath")
                func_icon=$(get_icon "$func_percent")
                formatted_func_percent=$(format_percent "$func_percent")
                echo "| $func_icon \`$filename:$funcname\` | $formatted_func_percent |" >> "$OUTPUT_FILE"
            fi
        done < coverage_detail.txt

        echo "" >> "$OUTPUT_FILE"
    fi
done

# Ajouter la date de génération
cat >> "$OUTPUT_FILE" << EOF

---

*Généré le: $(date '+%Y-%m-%d %H:%M:%S')*
EOF

# Nettoyage des fichiers temporaires
rm -f coverage.txt packages.txt coverage.out coverage.html coverage_detail.out coverage_detail.txt

echo -e "${GREEN}✅ Rapport généré: $OUTPUT_FILE${NC}"
echo ""
echo -e "${YELLOW}Couverture globale: $TOTAL_ICON $(format_percent "$TOTAL_PERCENT")${NC}"

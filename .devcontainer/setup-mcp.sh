#!/bin/bash
# Ne pas utiliser set -e pour permettre la récupération d'erreurs

VAULT_ID="ypahjj334ixtiyjkytu5hij2im"
MCP_TPL="/workspace/.devcontainer/mcp.json.tpl"
MCP_OUTPUT="/home/vscode/.devcontainer/mcp.json"

# Créer le répertoire de destination
mkdir -p "$(dirname "$MCP_OUTPUT")"

# Créer un mcp.json vide par défaut (sera écrasé si 1Password fonctionne)
create_empty_mcp() {
    echo '{"mcpServers":{}}' > "$MCP_OUTPUT"
    echo "📝 Fichier mcp.json vide créé (MCP désactivé)"
}

echo "🔐 Récupération des secrets depuis 1Password..."

# Vérifier que op est installé
if ! command -v op &> /dev/null; then
    echo "⚠️  1Password CLI n'est pas installé - MCP désactivé"
    create_empty_mcp
else
    # Récupérer les tokens depuis 1Password
    echo "  → Récupération du token Codacy..."
    CODACY_TOKEN=$(op item get "mcp-codacy" --vault "$VAULT_ID" --fields credential --reveal 2>/dev/null || echo "")

    echo "  → Récupération du token GitHub..."
    GITHUB_TOKEN=$(op item get "mcp-github" --vault "$VAULT_ID" --fields credential --reveal 2>/dev/null || echo "")

    # Si aucun token n'est récupéré, créer un fichier vide
    if [ -z "$CODACY_TOKEN" ] && [ -z "$GITHUB_TOKEN" ]; then
        echo "⚠️  Aucun token récupéré depuis 1Password - MCP désactivé"
        create_empty_mcp
    else
        # Générer le fichier mcp.json à partir du template
        echo "📝 Génération du fichier mcp.json..."
        sed "s|{{ with secret \"secret/mcp/codacy\" }}{{ .Data.data.token }}{{ end }}|${CODACY_TOKEN}|g" "$MCP_TPL" | \
            sed "s|{{ with secret \"secret/mcp/github\" }}{{ .Data.data.token }}{{ end }}|${GITHUB_TOKEN}|g" \
            > "$MCP_OUTPUT"
        echo "✅ Fichier mcp.json généré avec succès!"
    fi
fi

# Configurer les paramètres Claude CLI
echo "⚙️  Configuration de Claude CLI..."
mkdir -p /home/vscode/.claude
cat > /home/vscode/.claude/settings.json <<'EOF'
{
  "enableAllProjectMcpServers": true,
  "alwaysThinkingEnabled": true
}
EOF
echo "✅ Paramètres Claude CLI configurés!"

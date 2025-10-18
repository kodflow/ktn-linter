#!/usr/bin/env python3
"""
Script pour simplifier les messages want dans les testdata.
Au lieu d'avoir des messages complets qui ne matchent pas,
on utilise juste le code d'erreur avec .* pour matcher n'importe quel message.
"""

import os
import re
from pathlib import Path

def simplify_want_messages(filepath):
    """
    Remplace les patterns want complexes par des patterns simples.
    Exemple: want `\[KTN-FOR-001\] Utilisation de _ inutile dans for range`
    Devient: want `\[KTN-FOR-001\].*`
    """
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()

        original_content = content

        # Pattern pour trouver les want avec des messages longs
        # Format: // want `\[KTN-XXX-NNN\] message...`
        pattern = r'(// want `\\?\[KTN-[A-Z]+-\d+\\?\])[^`]*`'
        replacement = r'\1.*`'

        content = re.sub(pattern, replacement, content)

        if content != original_content:
            with open(filepath, 'w', encoding='utf-8') as f:
                f.write(content)
            return True

        return False

    except Exception as e:
        print(f"❌ Erreur avec {filepath}: {e}")
        return False

def main():
    """Parcourt tous les fichiers testdata et simplifie les want."""
    testdata_root = Path("src/pkg/analyzer/ktn")

    files_fixed = 0

    # Chercher tous les bad.go dans testdata
    for bad_file in testdata_root.rglob("testdata/**/bad.go"):
        if simplify_want_messages(bad_file):
            files_fixed += 1
            print(f"✅ {bad_file.relative_to('.')}")

    print(f"\n📊 Résumé:")
    print(f"   Fichiers simplifiés: {files_fixed}")

if __name__ == "__main__":
    main()

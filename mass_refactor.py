#!/usr/bin/env python3
"""
Script pour refactoriser automatiquement les fonctions trop longues.
"""

import re
import os
from pathlib import Path

# Fichiers à ignorer (tests bad_usage)
IGNORE_PATTERNS = [
    'bad_usage',
    '_test.go',
    'testdata'
]

def count_function_lines(func_text):
    """Compte les lignes d'une fonction."""
    lines = [l for l in func_text.split('\n') if l.strip()]
    return len(lines)

def should_ignore(filepath):
    """Vérifie si un fichier doit être ignoré."""
    for pattern in IGNORE_PATTERNS:
        if pattern in filepath:
            return True
    return False

def find_long_functions():
    """Trouve toutes les fonctions > 35 lignes."""
    src_path = Path('src')
    long_functions = []

    for go_file in src_path.rglob('*.go'):
        if should_ignore(str(go_file)):
            continue

        try:
            content = go_file.read_text()
            # Recherche basique de fonctions
            func_pattern = r'func\s+(\w+)\s*\([^)]*\)[^{]*{'
            for match in re.finditer(func_pattern, content):
                func_name = match.group(1)
                start = match.start()
                # Trouver la fin de la fonction (simpliste)
                # Compter les accolades
                brace_count = 1
                pos = match.end()
                while pos < len(content) and brace_count > 0:
                    if content[pos] == '{':
                        brace_count += 1
                    elif content[pos] == '}':
                        brace_count -= 1
                    pos += 1

                func_text = content[start:pos]
                lines = count_function_lines(func_text)

                if lines > 35:
                    long_functions.append({
                        'file': str(go_file),
                        'name': func_name,
                        'lines': lines,
                        'text': func_text
                    })
        except Exception as e:
            print(f"Error processing {go_file}: {e}")

    return long_functions

if __name__ == '__main__':
    functions = find_long_functions()
    print(f"Found {len(functions)} functions > 35 lines:\n")

    for f in sorted(functions, key=lambda x: x['lines'], reverse=True):
        print(f"{f['file']:60s} {f['name']:30s} {f['lines']:3d} lines")

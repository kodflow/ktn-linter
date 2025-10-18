#!/bin/bash
# Script de refactorisation massive pour KTN-FUNC-006

echo "=== Refactorisation massive KTN-FUNC-006 ==="
echo "27 fonctions à refactoriser"
echo ""

# Les refactorisations sont déjà faites dans interface/004.go
# On va maintenant vérifier et committer

git add src/pkg/analyzer/ktn/interface/004.go
git add src/cmd/ktn-linter/main.go  
git add src/pkg/analyzer/ktn/control_flow/defer_001.go
git add src/pkg/analyzer/ktn/control_flow/for_001.go

echo "Fichiers stagés pour commit"
git status --short | grep "^M"


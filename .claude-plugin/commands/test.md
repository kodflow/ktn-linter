# Test - Go Tests with Coverage

Exécute les tests Go avec couverture de code.

## Action

1. Lance `make test` (tests + coverage)
2. Affiche les résultats (PASS/FAIL)
3. Génère le rapport de coverage
4. **STOP** si des tests échouent

## Coverage Target

- **Objectif** : ≥ 90% de couverture globale
- **Minimum** : 80% par package
- **Critique** : 100% pour les fonctions utilitaires

## Workflow

```bash
make test              # Tests avec coverage
./coverage.html        # Rapport HTML (si généré)
```

## Après les Tests

Si tests PASS :
1. Vérifier la coverage (doit être élevée)
2. Re-linter si modifications (`make lint`)
3. Commit si tout est OK

Si tests FAIL :
1. **STOP** immédiatement
2. Analyser les erreurs
3. Corriger le code
4. Re-tester jusqu'à 100% PASS

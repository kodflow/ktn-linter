# Coverage Report

Rapport de couverture généré automatiquement.

**Légende:**
- 🟢 100% - Couverture complète
- 🟡 ≥90% - Bonne couverture
- 🔴 <90% - Couverture insuffisante
- ⚫ 0% - Pas de tests

---

## Coverage par Package

| Icon | Package | Coverage |
|:----:|---------|----------|
| 🟡 | **TOTAL (Global)** | **98.2%** |
| 🟢 | `github.com/kodflow/ktn-linter/cmd/ktn-linter` | 100.0% |
| 🟡 | `github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd` | 95.8% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn` | 100.0% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/const` | 96.6% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/func` | 97.9% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper` | 100.0% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/utils` | 100.0% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/formatter` | 100.0% |

---

## Détail des packages incomplets

### 🟡 `github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd` - 95.8%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `lint.go:loadPackages` | 71.4% |
| 🟡 `lint.go:runAnalyzers` | 95.8% |
| 🟡 `lint.go:createAnalysisPass` | 90.0% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/const` - 96.6%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `002.go:checkConstGrouping` | 84.6% |
| 🟡 `003.go:runConst003` | 93.3% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/func` - 97.9%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `001.go:runFunc001` | 93.3% |
| 🔴 `001.go:isLineToSkip` | 64.3% |
| 🔴 `001.go:countPureCodeLines` | 88.2% |


---

*Généré le: 2025-10-20 22:05:41*

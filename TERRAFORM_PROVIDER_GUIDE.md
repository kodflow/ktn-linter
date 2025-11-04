# Guide KTN-Linter pour Terraform Provider

## Problèmes Identifiés et Solutions

### 1. KTN-VAR-014 sur `schema.Schema`

**Problème** : `schema.Schema` du framework Terraform dépasse 64 bytes.

**Solution immédiate** : Accepter cette erreur (contrainte du framework).

```go
func (p *N8nProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	// NOTE: KTN-VAR-014 ignoré - schema.Schema est une contrainte du framework Terraform
	resp.Schema = schema.Schema{
		MarkdownDescription: "Terraform provider for n8n",
	}
}
```

**Solution future** : Implémenter `.ktnlintrc` avec exemption pour fichiers provider.

---

### 2. KTN-FUNC-007 - Format Params

**Format EXACT requis** :

```go
// Metadata populates the provider metadata including type name and version.
//
// Params:
//   - _: The context for the operation (currently unused)
//   - _: The metadata request from Terraform (currently unused)
//   - resp: The response object to populate with provider metadata
func (p *N8nProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "n8n"
	resp.Version = p.version
}
```

**Points critiques** :
- ✅ Ligne vide après la description (`//`)
- ✅ `Params:` avec majuscule P et deux-points
- ✅ Indentation : 2 espaces + tiret (`  - param:`)
- ✅ Paramètres inutilisés avec underscore documentés comme `_: Description (unused)`

---

### 3. KTN-STRUCT-002 vs KTN-INTERFACE-001

**Problème** : Interface requise mais détectée comme inutilisée.

**Solution** : Utiliser l'interface dans une signature de fonction.

```go
// N8nProvider implements the Terraform provider.
type N8nProvider struct {
	version string
}

// N8nProviderInterface defines the provider contract.
type N8nProviderInterface interface {
	Metadata(context.Context, provider.MetadataRequest, *provider.MetadataResponse)
	Schema(context.Context, provider.SchemaRequest, *provider.SchemaResponse)
	Configure(context.Context, provider.ConfigureRequest, *provider.ConfigureResponse)
	Resources(context.Context) []func() resource.Resource
	DataSources(context.Context) []func() datasource.DataSource
}

// NewN8nProvider creates a provider instance that implements the interface.
//
// Params:
//   - version: Provider version string
//
// Returns:
//   - N8nProviderInterface: Provider interface implementation
func NewN8nProvider(version string) N8nProviderInterface {
	return &N8nProvider{version: version}
}
```

**✅ L'interface n'est plus "inutilisée"** car elle est retournée par `NewN8nProvider()`.

---

### 4. KTN-RETURN-002 vs KTN-VAR-010

**Contradiction identifiée** ✅

**Pour l'instant** :
- Dans les **fonctions** : Retourner slice vide `[]string{}`
- Pour les **variables globales** : Utiliser `nil` ou ignorer l'erreur

```go
// Resources returns supported resources.
//
// Params:
//   - _: Context (unused)
//
// Returns:
//   - []func() resource.Resource: List of resource factories
func (p *N8nProvider) Resources(_ context.Context) []func() resource.Resource {
	// KTN-RETURN-002: Retourner slice vide plutôt que nil
	return []func() resource.Resource{}
}
```

---

## Workflow Recommandé

### Étape 1 : Fixer les erreurs critiques

1. ✅ Corriger le format `Params:` dans toutes les fonctions
2. ✅ Utiliser l'interface `N8nProviderInterface` dans `NewN8nProvider()`
3. ✅ Retourner `[]func(){}` au lieu de `nil`

### Étape 2 : Accepter les erreurs framework

Ajouter des commentaires pour documenter les exceptions :

```go
// NOTE: KTN-VAR-014 ignoré - contrainte du framework Terraform
// NOTE: KTN-STRUCT-002 ignoré - interface imposée par provider.Provider
```

### Étape 3 : Attendre le support de .ktnlintrc

Une fois `.ktnlintrc` implémenté, créer :

```yaml
# .ktnlintrc
exemptions:
  - path: "**/provider/*.go"
    rules:
      - KTN-VAR-014
```

---

## Exemple Complet Conforme

```go
package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = &N8nProvider{}

// N8nProvider implements the Terraform provider for n8n automation platform.
type N8nProvider struct {
	version string
}

// N8nProviderInterface defines the contract for the n8n Terraform provider.
type N8nProviderInterface interface {
	Metadata(context.Context, provider.MetadataRequest, *provider.MetadataResponse)
	Schema(context.Context, provider.SchemaRequest, *provider.SchemaResponse)
	Configure(context.Context, provider.ConfigureRequest, *provider.ConfigureResponse)
	Resources(context.Context) []func() resource.Resource
	DataSources(context.Context) []func() datasource.DataSource
}

// Metadata populates the provider metadata including type name and version.
//
// Params:
//   - _: Context for the operation (unused)
//   - _: Metadata request from Terraform (unused)
//   - resp: Response object to populate
func (p *N8nProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "n8n"
	resp.Version = p.version
}

// Schema defines the provider configuration schema.
//
// Params:
//   - _: Context for the operation (unused)
//   - _: Schema request from Terraform (unused)
//   - resp: Response object to populate
func (p *N8nProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	// NOTE: KTN-VAR-014 ignoré - schema.Schema est imposé par le framework
	resp.Schema = schema.Schema{
		MarkdownDescription: "Terraform provider for n8n automation platform",
	}
}

// Configure initializes the provider with configuration.
//
// Params:
//   - ctx: Operation context
//   - req: Configuration request
//   - resp: Configuration response
func (p *N8nProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config N8nProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	// Exit early if configuration parsing failed
	if resp.Diagnostics.HasError() {
		return
	}
}

// Resources returns supported resources.
//
// Params:
//   - _: Context (unused)
//
// Returns:
//   - []func() resource.Resource: List of resource factories
func (p *N8nProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

// DataSources returns supported data sources.
//
// Params:
//   - _: Context (unused)
//
// Returns:
//   - []func() datasource.DataSource: List of data source factories
func (p *N8nProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// NewN8nProvider creates a new N8nProvider instance.
//
// Params:
//   - version: Provider version string
//
// Returns:
//   - N8nProviderInterface: Provider implementation
func NewN8nProvider(version string) N8nProviderInterface {
	return &N8nProvider{
		version: version,
	}
}

// New returns a provider factory function.
//
// Params:
//   - version: Provider version string
//
// Returns:
//   - func() provider.Provider: Factory function
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return NewN8nProvider(version)
	}
}

// N8nProviderModel defines the provider configuration schema.
type N8nProviderModel struct {
	// Add configuration fields here
}
```

---

## Résumé

### ✅ À Faire Maintenant

1. Corriger le format `Params:` (voir exemples ci-dessus)
2. Utiliser l'interface dans `NewN8nProvider()`
3. Retourner `[]func(){}` au lieu de `nil`

### ⚠️ Accepter Temporairement

1. **KTN-VAR-014** sur `schema.Schema` (contrainte framework)
2. Ajouter des commentaires `// NOTE: KTN-XXX ignoré - raison`

### 🔮 Futur (Implémentation .ktnlintrc)

1. Exemptions par fichier/pattern
2. Désactivation de règles spécifiques
3. Seuils personnalisés

---

**Questions ?** Consulte `EXAMPLE_COMPLIANT_CODE.md` pour plus d'exemples.

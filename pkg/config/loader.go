// Package config provides configuration management for KTN linter rules.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	// DefaultConfigFileName is the default configuration file name.
	DefaultConfigFileName string = ".ktn-linter.yaml"
	// AlternateConfigFileName is an alternate configuration file name.
	AlternateConfigFileName string = ".ktn-linter.yml"
	// filePermReadWrite is the default file permission (0644 = rw-r--r--).
	filePermReadWrite = 0644
)

// Load loads configuration from a file path.
// If path is empty, it searches for default config files in the current directory and parent directories.
//
// Params:
//   - path: File path to load configuration from (empty for default locations)
//
// Returns:
//   - *Config: Loaded configuration
//   - error: Error if loading fails
func Load(path string) (*Config, error) {
	// Vérification si un chemin spécifique est fourni
	if path != "" {
		// Chargement depuis le fichier spécifié
		return loadFromFile(path)
	}

	// Recherche dans les emplacements par défaut
	return loadFromDefaultLocations()
}

// loadFromFile loads configuration from a specific file path.
//
// Params:
//   - path: File path to load configuration from
//
// Returns:
//   - *Config: Loaded configuration
//   - error: Error if loading fails
func loadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	// Vérification si la lecture du fichier a échoué
	if err != nil {
		// Retour d'erreur si impossible de lire le fichier
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	cfg := DefaultConfig()
	// Tentative de désérialisation YAML
	if err := yaml.Unmarshal(data, cfg); err != nil {
		// Retour d'erreur si le parsing YAML échoue
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	// Validate configuration
	// Validation de la configuration chargée
	if err := validateConfig(cfg); err != nil {
		// Retour d'erreur si la validation échoue
		return nil, fmt.Errorf("invalid config in %s: %w", path, err)
	}

	// Retour de la configuration valide
	return cfg, nil
}

// loadFromDefaultLocations searches for config files in default locations.
//
// Returns:
//   - *Config: Loaded configuration or default config
//   - error: Error if loading fails (returns default config on error)
func loadFromDefaultLocations() (*Config, error) {
	cwd, err := os.Getwd()
	// Vérification si impossible d'obtenir le répertoire courant
	if err != nil {
		// Retour de la configuration par défaut si erreur
		return DefaultConfig(), nil
	}

	// Search up the directory tree
	dir := cwd
	for {
		// Try default filename
		path := filepath.Join(dir, DefaultConfigFileName)
		// Vérification si le fichier par défaut existe
		if fileExists(path) {
			// Chargement depuis le fichier trouvé
			return loadFromFile(path)
		}

		// Try alternate filename
		path = filepath.Join(dir, AlternateConfigFileName)
		// Vérification si le fichier alternatif existe
		if fileExists(path) {
			// Chargement depuis le fichier trouvé
			return loadFromFile(path)
		}

		// Move to parent directory
		parent := filepath.Dir(dir)
		// Vérification si on a atteint la racine du système
		if parent == dir {
			// Sortie de la boucle si racine atteinte
			break
		}
		dir = parent
	}

	// No config file found, return default
	// Retour de la configuration par défaut si aucun fichier trouvé
	return DefaultConfig(), nil
}

// fileExists checks if a file exists.
//
// Params:
//   - path: File path to check
//
// Returns:
//   - bool: true if file exists and is not a directory
func fileExists(path string) bool {
	info, err := os.Stat(path)
	// Retour true si le fichier existe et n'est pas un répertoire
	return err == nil && !info.IsDir()
}

// validateConfig validates the configuration.
//
// Params:
//   - cfg: Configuration to validate
//
// Returns:
//   - error: Validation error if any
func validateConfig(cfg *Config) error {
	// Vérification si la configuration est nulle
	if cfg == nil {
		// Retour immédiat si configuration nulle (pas d'erreur)
		return nil
	}

	// Validate version
	// Vérification que la version est supportée
	if cfg.Version != 0 && cfg.Version != 1 {
		// Retour d'erreur si version non supportée
		return fmt.Errorf("unsupported config version: %d", cfg.Version)
	}

	// Validate rules
	for code, ruleCfg := range cfg.Rules {
		// Vérification si la configuration de règle est nulle
		if ruleCfg == nil {
			continue
		}

		// Validate threshold is positive if set
		// Vérification que le seuil est non-négatif si défini
		if ruleCfg.Threshold != nil && *ruleCfg.Threshold < 0 {
			// Retour d'erreur si seuil négatif
			return fmt.Errorf("rule %s: threshold must be non-negative, got %d", code, *ruleCfg.Threshold)
		}

		// Validate exclusion patterns
		for _, pattern := range ruleCfg.Exclude {
			// Vérification si le pattern est vide
			if pattern == "" {
				// Retour d'erreur si pattern vide
				return fmt.Errorf("rule %s: empty exclusion pattern", code)
			}
		}
	}

	// Validate global exclusions
	for _, pattern := range cfg.Exclude {
		// Vérification si le pattern global est vide
		if pattern == "" {
			// Retour d'erreur si pattern vide
			return fmt.Errorf("empty global exclusion pattern")
		}
	}

	// Retour sans erreur si toutes les validations passent
	return nil
}

// LoadAndSet loads configuration and sets it as the global config.
//
// Params:
//   - path: File path to load configuration from
//
// Returns:
//   - error: Error if loading fails
func LoadAndSet(path string) error {
	cfg, err := Load(path)
	// Vérification si le chargement a échoué
	if err != nil {
		// Retour de l'erreur de chargement
		return err
	}

	Set(cfg)

	// Retour sans erreur après configuration globale définie
	return nil
}

// MustLoad loads configuration and panics on error.
//
// Params:
//   - path: File path to load configuration from
//
// Returns:
//   - *Config: Loaded configuration (panics on error)
func MustLoad(path string) *Config {
	cfg, err := Load(path)
	// Vérification si le chargement a échoué
	if err != nil {
		// Panic avec l'erreur si chargement impossible
		panic(err)
	}

	// Retour de la configuration chargée
	return cfg
}

// SaveToFile saves configuration to a file.
//
// Params:
//   - cfg: Configuration to save
//   - path: File path to save configuration to
//
// Returns:
//   - error: Error if saving fails
func SaveToFile(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	// Vérification si la sérialisation YAML a échoué
	if err != nil {
		// Retour d'erreur si impossible de sérialiser
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Tentative d'écriture du fichier avec permissions rw-r--r--
	if err := os.WriteFile(path, data, filePermReadWrite); err != nil {
		// Retour d'erreur si impossible d'écrire le fichier
		return fmt.Errorf("failed to write config file %s: %w", path, err)
	}

	// Retour sans erreur après sauvegarde réussie
	return nil
}

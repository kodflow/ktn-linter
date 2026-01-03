// Package cmd implements the CLI commands for ktn-linter.
package cmd

import "github.com/kodflow/ktn-linter/pkg/updater"

// flagGetter abstracts flag access for testability.
// This interface allows mocking Cobra command flags in tests.
type flagGetter interface {
	GetBool(name string) (bool, error)
	GetString(name string) (string, error)
}

// updaterService abstracts the updater for testability.
// This interface allows mocking the update functionality in tests.
type updaterService interface {
	CheckForUpdate() (updater.UpdateInfo, error)
	Upgrade() (updater.UpdateInfo, error)
}

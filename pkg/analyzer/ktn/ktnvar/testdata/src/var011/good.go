// Package var011 provides good test cases.
package var011

import (
	"context"
	"fmt"
)

const (
	// KeyValue is a constant test value
	KeyValue int = 42
)

// init demonstrates correct usage patterns
func init() {
	// Shadowing de err est autorisé (pattern idiomatique Go)
	var err error
	// Vérification de la condition
	if err != nil {
		err := fmt.Errorf("wrapped: %w", err) // OK: 'err' est exemptée
		_ = err
	}

	// Shadowing de ok est autorisé (map access/type assertion)
	m := map[string]int{"key": KeyValue}
	v, ok := m["key"]
	// Vérification de la condition
	if ok {
		_, ok := m["other"] // OK: 'ok' est exemptée
		_ = ok
	}
	_ = v

	// Shadowing de ctx est autorisé (redéfinition de context)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// Bloc imbriqué
	{
		ctx := context.WithValue(ctx, "key", "value") // OK: 'ctx' est exemptée
		_ = ctx
	}
}

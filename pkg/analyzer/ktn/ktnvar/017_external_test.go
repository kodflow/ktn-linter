package ktnvar_test

import (
	"testing"

	"golang.org/x/tools/go/analysis"

	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

// TestVar017 vérifie la détection des copies de mutex et types sync/atomic.
// Erreurs attendues dans bad.go:
// - counter.increment: copie sync.Mutex (1)
// - safeCounter.read: copie sync.RWMutex (1)
// - config.load: copie atomic.Value (1)
// - container.values: copie atomic.Int32 (1)
// - outer.doSomething: copie sync.Mutex via inner (1)
// - flag.isActive: copie atomic.Bool (1)
// - stats.count: copie atomic.Uint64 (1)
// - badMutexCopy: param sync.Mutex (1)
// - badAssignment: copie sync.Mutex (1)
// - badRWMutexParam: param sync.RWMutex (1)
// - badAtomicValueParam: param atomic.Value (1)
// - badRWMutexAssignment: copie sync.RWMutex (1)
// - badAtomicValueAssignment: copie atomic.Value (1)
// - badWaitGroupParam: param sync.WaitGroup (1)
// - badOnceParam: param sync.Once (1)
// - badCondParam: param sync.Cond (1)
// - badWaitGroupAssignment: copie sync.WaitGroup (1)
// - badOnceAssignment: copie sync.Once (1)
// Total: 30 erreurs (inclut structs avec value receivers + receivers + params + assignments)
func TestVar017(t *testing.T) {
	tests := []struct {
		name           string
		analyzer       *analysis.Analyzer
		testdataDir    string
		expectedErrors int
	}{
		{
			name:           "Mutex and sync type copies",
			analyzer:       ktnvar.Analyzer017,
			testdataDir:    "var017",
			expectedErrors: 30,
		},
		{
			name:           "Valid pointer usage for sync types",
			analyzer:       ktnvar.Analyzer017,
			testdataDir:    "var017",
			expectedErrors: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelper.TestGoodBad(t, tt.analyzer, tt.testdataDir, tt.expectedErrors)
		})
	}
}

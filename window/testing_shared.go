package window

import (
	"sync"
	"testing"

	"github.com/therecipe/qt/core"
)

var tRunner = NewTestRunner(nil)

type testRunner struct {
	core.QObject

	_ func(f func()) `signal:"runOnMain,auto"`
}

func (t *testRunner) runOnMain(f func()) { f() }

// Run doesn't require serialization.
func (t *testRunner) Run(f func()) {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	t.RunOnMain(func() {
		f()
		wg.Done()
	})
	wg.Wait()
}

// RunT returns the return value from the test function. You can serialize your
// tests by returning false.
func (t *testRunner) RunT(tt *testing.T, f func(*testing.T) bool) (stop bool) {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	t.RunOnMain(func() {
		defer wg.Done()
		if !f(tt) {
			stop = true
		}
	})
	wg.Wait()
	return stop
}

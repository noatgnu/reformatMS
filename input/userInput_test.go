package input

import (
	"testing"
	"path/filepath"
)

func TestClean(t *testing.T) {
	result := Clean(`/ ""`)
	f, _ := filepath.Abs(`\`)
	if result != f {
		t.Errorf(`Expected "\" got %v`, result)
	}
}

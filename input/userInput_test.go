package input

import "testing"

func TestClean(t *testing.T) {
	result := Clean(`/ ""`)
	if result != `\` {
		t.Errorf(`Expected "\" got %v`, result)
	}
}

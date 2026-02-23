package retrieval

import "testing"

func TestIsSemanticConfident(t *testing.T) {
	if IsSemanticConfident(0.90, 0.82) {
		t.Fatal("expected low margin to fail")
	}
	if !IsSemanticConfident(0.90, 0.60) {
		t.Fatal("expected clear margin to pass")
	}
}

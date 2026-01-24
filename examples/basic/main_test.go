package main

import (
	"testing"
)

// TestExampleCompilation ensures the example compiles without issues
func TestExampleCompilation(t *testing.T) {
	// This test just ensures the example compiles and imports work correctly
	// The actual main() function is not called to avoid requiring external dependencies
	// (CGO and native libraries are required to run the example)
	t.Log("Example compiles successfully")
}

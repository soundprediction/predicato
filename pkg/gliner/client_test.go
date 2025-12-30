package gliner

import (
	"fmt"
	"testing"
)

func TestClient(t *testing.T) {
	// Use small model
	modelID := "onnx-community/gliner_small-v2.1"

	c, err := NewClient(modelID)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	text := "Apple was founded by Steve Jobs in California."
	labels := []string{"person", "organization", "location"}

	ents, err := c.ExtractEntities(text, labels)
	if err != nil {
		t.Fatalf("Extraction failed: %v", err)
	}

	fmt.Printf("Entities: %+v\n", ents)

	foundOrg := false
	for _, e := range ents {
		if e.Label == "organization" && e.Text == "Apple" {
			foundOrg = true
		}
	}
	if !foundOrg {
		t.Error("Failed to find Apple as organization")
	}
}

package prompts

import (
	"strings"
	"testing"
)

// Minimal test structure for verifying prompt format
func TestGetPromptFormat(t *testing.T) {
	tests := []struct {
		name     string
		context  map[string]interface{}
		expected string
	}{
		{
			name:     "Default TSV",
			context:  map[string]interface{}{},
			expected: "TSV",
		},
		{
			name: "Use YAML",
			context: map[string]interface{}{
				"use_yaml": true,
			},
			expected: "YAML",
		},
		{
			name: "Use TOML",
			context: map[string]interface{}{
				"use_toml": true,
			},
			expected: "TOML",
		},
		{
			name: "Use TOML Priority", // What if both are true? Implementation order matters. YAML is checked first in current impl.
			context: map[string]interface{}{
				"use_yaml": true,
				"use_toml": true,
			},
			expected: "YAML",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format := GetPromptFormat(tt.context)
			if format.Name != tt.expected {
				t.Errorf("GetPromptFormat() = %v, want %v", format.Name, tt.expected)
			}
		})
	}
}

func TestToPromptTOML(t *testing.T) {
	data := []map[string]interface{}{
		{
			"entity":         "John Doe",
			"entity_type_id": 1,
		},
		{
			"entity":         "Jane Smith",
			"entity_type_id": 2,
		},
	}

	tomlStr, err := ToPromptTOML(data)
	if err != nil {
		t.Fatalf("ToPromptTOML failed: %v", err)
	}

	// Verify TOML structure
	// [[entities]]
	// entity = "John Doe"
	// entity_type_id = 1
	// ... or similar array of tables structure
	// go-toml v2 marshals slices as array of tables if they are maps?
	// Actually, top level slice marshals to ... what?
	// TOML spec says top level must be a table.
	// But `ToPromptTOML` uses `toml.Marshal(data)`.
	// If data is a slice, go-toml v2 might error or return a keyless array?
	// Wait, TOML root must be a table (hash). It cannot be an array.
	// So `ToPromptTOML` mightFAIL if passed a slice directly at the root?
	// Let's check this behavior. If it fails, we need to wrap the slice in a map or change how we invoke it.
	// The `extractMessagePrompt` wraps it?
	// No, `ToPromptTOML(filteredEntityTypes)` where `filteredEntityTypes` is `[]map...`.
	// Let's see if this works.

	// If it fails, that's a bug in my implementation and I must fix it.
	// The test will reveal it.
	t.Logf("TOML Output:\n%s", tomlStr)
	if !strings.Contains(tomlStr, "John Doe") {
		t.Errorf("Expected TOML to contain 'John Doe'")
	}
}

func TestFormatContext(t *testing.T) {
	context := map[string]interface{}{
		"use_toml": true,
	}
	format := GetPromptFormat(context)

	data := map[string]interface{}{
		"Test": map[string]string{"foo": "bar"},
	}

	result, err := FormatContext(format, data)
	if err != nil {
		t.Fatalf("FormatContext failed: %v", err)
	}

	if val, ok := result["Test"]; !ok {
		t.Errorf("Expected result to contain key 'Test'")
	} else {
		// Verify TOML formatting
		if !strings.Contains(val, "foo = 'bar'") && !strings.Contains(val, `foo = "bar"`) {
			t.Errorf("Expected TOML formatted string, got: %v", val)
		}
	}
}

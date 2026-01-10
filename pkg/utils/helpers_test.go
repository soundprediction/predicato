package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestItem struct {
	Name  string `yaml:"name"`
	Value int    `yaml:"value"`
}

func TestUnmarshalYAML_Robustness(t *testing.T) {
	// Mixed content:
	// 1. Valid item
	// 2. Invalid item (Value is a string "invalid")
	// 3. Valid item
	yamlData := `
- name: item1
  value: 10
- name: item2
  value: "invalid_int"
- name: item3
  value: 30
`

	items, err := UnmarshalYAML[TestItem](yamlData)

	assert.NoError(t, err)
	assert.Len(t, items, 2, "Should return 2 valid items")

	assert.Equal(t, "item1", items[0].Name)
	assert.Equal(t, 10, items[0].Value)

	assert.Equal(t, "item3", items[1].Name)
	assert.Equal(t, 30, items[1].Value)
}

func TestUnmarshalYAML_AllInvalid(t *testing.T) {
	yamlData := `
- name: item1
  value: "invalid"
`
	items, err := UnmarshalYAML[TestItem](yamlData)

	assert.Error(t, err)
	assert.Nil(t, items)
}

func TestUnmarshalYAML_MalformedStructure(t *testing.T) {
	yamlData := `
this is not a list
`
	items, err := UnmarshalYAML[TestItem](yamlData)

	assert.Error(t, err)
	assert.Nil(t, items)
}

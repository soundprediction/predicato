package nlp_test

import (
	"testing"

	"github.com/soundprediction/predicato/pkg/nlp"
	"github.com/stretchr/testify/assert"
)

func TestGetProvider(t *testing.T) {
	tests := []struct {
		id      nlp.ProviderID
		want    nlp.Provider
		wantErr bool
	}{
		{
			id: nlp.ProviderEmbedEverything,
			want: nlp.Provider{
				ID:          nlp.ProviderEmbedEverything,
				Name:        "EmbedEverything",
				Description: "Local generic embedding models via Rust bindings",
				IsLocal:     true,
			},
			wantErr: false,
		},
		{
			id:      "nonexistent",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.id), func(t *testing.T) {
			got, found := nlp.GetProvider(tt.id)
			if tt.wantErr {
				assert.False(t, found)
			} else {
				assert.True(t, found)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetModel(t *testing.T) {
	// Pick a few representative models
	t.Run("EmbedEverything Model", func(t *testing.T) {
		id := "sentence-transformers/all-MiniLM-L6-v2"
		got, found := nlp.GetModel(id)
		assert.True(t, found)
		assert.Equal(t, id, got.ID)
		assert.Contains(t, got.Capabilities, nlp.TaskEmbedding)
	})

	t.Run("GLiNER Model", func(t *testing.T) {
		id := "urchade/gliner_multi-v2.1"
		got, found := nlp.GetModel(id)
		assert.True(t, found)
		assert.Equal(t, id, got.ID)
		assert.Contains(t, got.Capabilities, nlp.TaskNamedEntityRecognition)
		assert.Contains(t, got.Capabilities, nlp.TaskRelationExtraction)
	})

	t.Run("RustBert Model", func(t *testing.T) {
		id := "bert-base-ner"
		got, found := nlp.GetModel(id)
		assert.True(t, found)
		assert.Equal(t, id, got.ID)
		assert.Contains(t, got.Capabilities, nlp.TaskNamedEntityRecognition)
	})

	t.Run("Nonexistent Model", func(t *testing.T) {
		_, found := nlp.GetModel("fake-model")
		assert.False(t, found)
	})
}

func TestGetModelsByProvider(t *testing.T) {
	models := nlp.GetModelsByProvider(nlp.ProviderEmbedEverything)
	assert.NotEmpty(t, models)
	for _, m := range models {
		assert.Equal(t, nlp.ProviderEmbedEverything, m.ProviderID)
	}
}

func TestGetModelsByCapability(t *testing.T) {
	t.Run("Embedding", func(t *testing.T) {
		models := nlp.GetModelsByCapability(nlp.TaskEmbedding)
		assert.NotEmpty(t, models)
		for _, m := range models {
			assert.Contains(t, m.Capabilities, nlp.TaskEmbedding)
		}
	})

	t.Run("NER", func(t *testing.T) {
		models := nlp.GetModelsByCapability(nlp.TaskNamedEntityRecognition)
		assert.NotEmpty(t, models)
		// Both GLiNER and RustBert models should be here
		hasGliner := false
		hasRustBert := false
		for _, m := range models {
			if m.ProviderID == nlp.ProviderGLiNER {
				hasGliner = true
			}
			if m.ProviderID == nlp.ProviderRustBert {
				hasRustBert = true
			}
		}
		assert.True(t, hasGliner, "Should have GLiNER NER models")
		assert.True(t, hasRustBert, "Should have RustBert NER models")
	})
}

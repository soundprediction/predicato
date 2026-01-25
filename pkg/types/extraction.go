package types

import (
	"time"

	"github.com/soundprediction/predicato/pkg/factstore"
)

// ExtractionResults is returned when ExtractOnly=true in AddEpisodeOptions.
// It contains the raw extractions before graph modeling, allowing for custom
// processing, filtering, or alternative graph modeling logic.
type ExtractionResults struct {
	// SourceID is the episode/document ID in the fact store
	SourceID string `json:"source_id"`

	// ExtractedNodes are raw entities before resolution/deduplication
	ExtractedNodes []*factstore.ExtractedNode `json:"extracted_nodes"`

	// ExtractedEdges are raw relationships before resolution
	ExtractedEdges []*factstore.ExtractedEdge `json:"extracted_edges"`

	// ChunkCount is the number of chunks the episode was split into
	ChunkCount int `json:"chunk_count"`

	// ExtractionTime is how long the extraction took
	ExtractionTime time.Duration `json:"extraction_time"`

	// Metadata contains additional extraction information
	Metadata *ExtractionMetadata `json:"metadata,omitempty"`
}

// ExtractionMetadata contains additional information about the extraction process.
type ExtractionMetadata struct {
	// ModelUsed is the NLP model used for extraction
	ModelUsed string `json:"model_used,omitempty"`

	// EmbeddingModel is the model used for generating embeddings
	EmbeddingModel string `json:"embedding_model,omitempty"`

	// EmbeddingDimensions is the dimension of the embeddings
	EmbeddingDimensions int `json:"embedding_dimensions,omitempty"`

	// TotalTokens is the estimated token count processed
	TotalTokens int `json:"total_tokens,omitempty"`
}

// NodeCount returns the number of extracted nodes.
func (r *ExtractionResults) NodeCount() int {
	if r == nil {
		return 0
	}
	return len(r.ExtractedNodes)
}

// EdgeCount returns the number of extracted edges.
func (r *ExtractionResults) EdgeCount() int {
	if r == nil {
		return 0
	}
	return len(r.ExtractedEdges)
}

// IsEmpty returns true if no entities or relationships were extracted.
func (r *ExtractionResults) IsEmpty() bool {
	return r.NodeCount() == 0 && r.EdgeCount() == 0
}

package factstore

import (
	"fmt"
)

// NewFactsDB creates a new FactsDB instance based on the configuration.
// - FactStoreTypePostgres: Uses external PostgreSQL with VectorChord for native vector search
// - FactStoreTypeDoltGres: Uses DoltGres with in-memory vector search (no VectorChord support)
// If Type is empty, defaults to DoltGres.
func NewFactsDB(config *FactStoreConfig) (FactsDB, error) {
	if config == nil {
		return nil, fmt.Errorf("factstore config is required")
	}

	// Default values
	if config.EmbeddingDimensions <= 0 {
		config.EmbeddingDimensions = 1024 // Default for qwen3-embedding
	}

	if config.ConnectionString == "" {
		return nil, fmt.Errorf("connection string is required")
	}

	switch config.Type {
	case FactStoreTypePostgres:
		// External PostgreSQL with VectorChord for native vector search
		return NewPostgresDB(config.ConnectionString, config.EmbeddingDimensions)

	case FactStoreTypeDoltGres, "":
		// DoltGres without VectorChord - uses in-memory vector search
		return NewDoltGresDB(config.ConnectionString, config.EmbeddingDimensions)

	default:
		return nil, fmt.Errorf("unsupported factstore type: %s (supported: postgres, doltgres)", config.Type)
	}
}

// NewFactsDBFromURL creates a FactsDB from a simple connection URL.
// This is a convenience function for quick setup with external PostgreSQL + VectorChord.
// For PostgreSQL: "postgres://user:pass@host:5432/dbname"
// For DoltGres, use NewFactsDB with FactStoreTypeDoltGres explicitly.
func NewFactsDBFromURL(connectionURL string, embeddingDimensions int) (FactsDB, error) {
	return NewFactsDB(&FactStoreConfig{
		Type:                FactStoreTypePostgres,
		ConnectionString:    connectionURL,
		EmbeddingDimensions: embeddingDimensions,
	})
}

// NewDoltGresFactsDB creates a FactsDB for DoltGres (without VectorChord).
// This is a convenience function for DoltGres setup.
func NewDoltGresFactsDB(connectionURL string, embeddingDimensions int) (FactsDB, error) {
	return NewFactsDB(&FactStoreConfig{
		Type:                FactStoreTypeDoltGres,
		ConnectionString:    connectionURL,
		EmbeddingDimensions: embeddingDimensions,
	})
}

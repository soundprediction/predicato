// Package factstore provides persistent storage for extracted facts and entities.
//
// This package implements the FactsDB interface for storing and searching
// extracted nodes, edges, and sources from the knowledge graph processing.
//
// # Supported Backends
//
// The following storage backends are supported:
//   - PostgresDB: PostgreSQL with optional pgvector extension for vector search
//   - DoltDB: Dolt SQL database (deprecated, use PostgresDB with DoltGres)
//
// # Usage
//
//	config := &factstore.FactStoreConfig{
//	    Type:             "postgres",
//	    ConnectionString: "postgres://user:pass@localhost:5432/facts",
//	    UsePgVector:      true,
//	}
//	db, err := factstore.NewFactsDB(config)
//	if err != nil {
//	    return err
//	}
//	defer db.Close()
//
//	// Initialize tables
//	if err := db.Initialize(ctx); err != nil {
//	    return err
//	}
//
// # Search Capabilities
//
// The package provides hybrid search combining:
//   - Vector similarity search (using embeddings)
//   - Keyword/fulltext search
//   - Reciprocal Rank Fusion (RRF) for result merging
//
// # Vector Search
//
// For best performance, use PostgreSQL with pgvector extension.
// Without pgvector (e.g., DoltGres), vector search falls back to
// in-memory cosine similarity calculation.
package factstore

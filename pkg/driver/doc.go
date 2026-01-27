// Package driver provides graph database driver implementations for predicato.
//
// This package defines the GraphDriver interface and provides implementations
// for various graph databases including Neo4j, Memgraph, FalkorDB, and Ladybug.
//
// # Supported Databases
//
// The following graph databases are supported:
//   - Neo4j: Full-featured graph database with native vector search
//   - Memgraph: In-memory graph database for high performance
//   - FalkorDB: Redis-based graph database
//   - Ladybug: Embedded graph database (requires CGO)
//
// # Usage
//
// Create a driver using the appropriate constructor:
//
//	// Neo4j
//	driver, err := driver.NewNeo4jDriver(ctx, uri, username, password)
//
//	// Ladybug (embedded)
//	driver, err := driver.NewLadybugDriver(dbPath)
//
// # Thread Safety
//
// All driver implementations are safe for concurrent use from multiple goroutines.
// Database connections are managed internally and pooled where appropriate.
//
// # Type Helpers
//
// The package provides safe type conversion helpers in type_helpers.go for
// converting database results to Go types without panicking on type assertion
// failures.
package driver

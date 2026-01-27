// Package search provides search functionality for the predicato knowledge graph.
//
// This package implements semantic search across the knowledge graph, combining
// vector similarity search with graph traversal to find relevant nodes and edges.
//
// # Search Methods
//
// The Searcher supports multiple search methods:
//   - Vector search: Find nodes/edges by embedding similarity
//   - Fulltext search: Find nodes/edges by keyword matching
//   - Hybrid search: Combine vector and fulltext results
//   - Graph traversal: Expand results by following relationships
//
// # Usage
//
//	searcher := search.NewSearcher(driver, embedder, nlpClient)
//
//	config := &types.SearchConfig{
//	    Limit:        10,
//	    MinScore:     0.7,
//	    IncludeEdges: true,
//	}
//
//	results, err := searcher.Search(ctx, "query text", config)
//
// # Reranking
//
// Search results can be optionally reranked using cross-encoder models
// for improved relevance. Enable reranking via SearchConfig.Rerank.
//
// # Filtering
//
// Results can be filtered by:
//   - GroupIDs: Limit to specific groups
//   - NodeTypes: Filter by node type (entity, episode, community)
//   - TimeRange: Filter by temporal validity
//
// # Internal Type Design
//
// This package defines its own SearchConfig, NodeSearchConfig, EdgeSearchConfig,
// and SearchFilters types that are separate from pkg/types. This is intentional:
//
//   - pkg/types provides a simplified public API with string-based configuration
//   - pkg/search provides a richer internal implementation with typed enums
//   - Conversion happens in retrieval.go when calling the Searcher
//
// This separation allows the public API to remain stable while the internal
// implementation can evolve. New search methods and rerankers can be added
// internally without changing the public interface.
package search

# go-predicato API Reference

This document provides a comprehensive reference for the `go-predicato` API, covering advanced configuration, search capabilities, and graph maintenance operations.

## Client Configuration

The `predicato.Client` is the main entry point. It is initialized via `NewClient`.

### Config Structure

```go
type Config struct {
    // GroupID isolates data for multi-tenant scenarios.
    GroupID string

    // TimeZone for temporal operations (default: UTC).
    TimeZone *time.Location

    // DefaultEntityTypes defines entity types to extract when not specified in AddEpisodeOptions.
    EntityTypes map[string]interface{}

    // SearchConfig sets the default search parameters.
    SearchConfig *types.SearchConfig
}
```

## Advanced Search

The `Client.Search` method supports a powerful hybrid search engine with configurable reranking.

### Search Configuration

The `SearchConfig` allows fine-grained control over how Nodes, Edges, and Communities are retrieved.

```go
type SearchConfig struct {
    NodeConfig      *NodeSearchConfig
    EdgeConfig      *EdgeSearchConfig
    CommunityConfig *CommunitySearchConfig
    Limit           int     // Max total results
    MinScore        float64 // Minimum relevance score (0.0 - 1.0)
}
```

### Search Methods

Each config type (`NodeConfig`, etc.) accepts a list of `SearchMethods`:

*   `"cosine_similarity"`: Vector processing using embeddings. Requires an `embedder`.
*   `"bm25"`: Keyword-based full-text search. High precision for exact term matches.
*   `"bfs"`: Breadth-First Search. Expands the result set by traversing edges from initial matches. Useful for finding related concepts.

### Reranking Strategies

`go-predicato` supports several reranking algorithms to refine search results:

*   `"rrf"` (**Reciprocal Rank Fusion**): Combines results from multiple search methods (e.g., BM25 + Cosine) without requiring improved embeddings. Robust and zero-config.
*   `"mmr"` (**Maximal Marginal Relevance**): Optimizes for diversity. detailed by `MMRLambda` (0.0 = diverse, 1.0 = relevant). Great for RAG to avoid redundant context.
*   `"cross_encoder"`: Uses a heavy-duty BERT-like model to score relevance. Highest accuracy but slower. Requires `crossencoder` client.
*   `"node_distance"`: Reranks based on graph distance from specific nodes.
*   "episode_mentions"`: Boosts nodes mentioned in specific episodes.

### Search Filters

Refine the scope of your search:

```go
type SearchFilters struct {
    GroupIDs    []string
    NodeTypes   []types.NodeType // e.g., types.EntityNodeType
    EdgeTypes   []types.EdgeType
    EntityTypes []string         // e.g., "Person", "Location"
    TimeRange   *types.TimeRange // Filter by ValidFrom/ValidTo
}
```

## Graph Maintenance

### Community Detection

`go-predicato` allows you to detect and materialize communities within your graph using the **Label Propagation Algorithm (LPA)**.

```go
// Detect communities, summarize them using LLM, and write Community Nodes back to the graph.
nodes, edges, err := client.UpdateCommunities(ctx, episodeUUID, groupID)
```

**Process:**
1.  **Clustering**: Runs Label Propagation to identify dense clusters of entities.
2.  **Summarization**: Uses the LLM to hierarchically summarize the entities in each cluster.
3.  **Materialization**: Creates `CommunityNode`s with these summaries and links them to member entities via `HAS_MEMBER` edges.

### Index Management

To ensure performance, explicitly create indices after bulk loading:

```go
// Creates vector indices and full-text search indices
err := client.CreateIndices(ctx)
```

## Robustness & Recovery

`go-predicato` includes built-in mechanisms to handle real-world failure modes.

### WAL Recovery

When using the embedded `ladybug` driver, the system automatically detects database corruption (e.g., from a hard crash or power loss).

*   **Detection**: If the DB fails to open with a corruption error.
*   **Recovery**: The driver automatically identifies the corrupt Write-Ahead Log (WAL), moves it to a backup location (e.g., `wal.12345678.corrupt`), and starts a fresh session.
*   **No Intervention Required**: Your application will simply start successfully after a crash.

### Concurrency Safety

The `ladybug` driver ensures thread safety regardless of your application's concurrency model:
*   **Single Connection**: Uses a single serialized connection for all read/write operations.
*   **Write Queue**: Async writes are buffered and executed sequentially by a dedicated worker.
*   **Read Locks**: Executions are mutex-protected.

You can safely share a single `Client` instance across hundreds of goroutines.

## Low-Level Graph Operations

For manual graph manipulation:

*   **`AddTriplet`**: Insert a Subject-Predicate-Object triple directly.
*   **`GetNode` / `GetEdge`**: Retrieve specific elements by UUID.
*   **`GetNodesAndEdgesByEpisode`**: Retrieve the subgraph associated with a specific ingestion event.

# predicato

A temporal information extraction and knowledge graph library for Go with internal NLP model capabilities (or support for external services).

## What Makes Predicato Different

Most agentic memory libraries require external services (LLMs, vector databases, graph databasesj).
Predicato has implemented embeedded alternatives for all components so it can be used without external service dependencies.
 Predicato is **modular by design** - every component can run locally OR connect to external services. Start with the internal stack for development, then swap in cloud services for production without changing your code.

| Component | Internal (No API) | External Services |
|-----------|-----------------|---------------------|
| **Graph Database** | Ladybug (embedded) | Neo4j, Memgraph |
| **Embeddings** | go-embedeverything | OpenAI compatible APIs, AWS bedrock, Gemini |
| **Reranking** | go-embedeverything | Jina, Cohere |
| **Text Generation** | go-rust-bert (BERT models) | OpenAI compatible APIs |
| **Entity Extraction** | GLiNER (ONNX) | LLM-based extraction |
| **Fact Storage** | DoltGres (embedded) | PostgreSQL + pgvector |

**Why choose Predicato:**
- **Security** - Don't expose your data to external services
- **Run offline** - Embedded database + local ML models = no network required
- **Swap components freely** - Same code works with local models or cloud APIs
- **Bi-temporal knowledge** - Track when facts were recorded AND when they were valid
- **Hybrid search** - Semantic + BM25 keyword + graph traversal in one query
- **Production hardened** - WAL recovery, circuit breakers, cost tracking, telemetry

## Quick Start (Internal Stack)

No API keys. No external services. Just Go and CGO.

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/soundprediction/predicato"
    "github.com/soundprediction/predicato/pkg/crossencoder"
    "github.com/soundprediction/predicato/pkg/driver"
    "github.com/soundprediction/predicato/pkg/embedder"
    "github.com/soundprediction/predicato/pkg/rustbert"
)

func main() {
    ctx := context.Background()

    // Embedded graph database (no server required)
    db, _ := driver.NewLadybugDriver("./knowledge.db", 1)
    defer db.Close(ctx)

    // Local text generation (GPT-2, no API)
    rustbertClient := rustbert.NewClient(rustbert.Config{})
    llmClient := rustbert.NewLLMAdapter(rustbertClient, "text_generation")
    defer rustbertClient.Close()

    // Local embeddings (no API)
    embedderClient, _ := embedder.NewEmbedEverythingClient(&embedder.EmbedEverythingConfig{
        Config: &embedder.Config{
            Model:      "qwen/qwen3-embedding-0.6b",
            Dimensions: 1024,
        },
    })
    defer embedderClient.Close()

    // Local reranking (no API)
    reranker, _ := crossencoder.NewEmbedEverythingClient(&crossencoder.EmbedEverythingConfig{
        Config: &crossencoder.Config{
            Model: "zhiqing/Qwen3-Reranker-0.6B-ONNX",
        },
    })
    defer reranker.Close()

    // Create client
    client, _ := predicato.NewClient(db, llmClient, embedderClient, &predicato.Config{
        GroupID:  "my-app",
        TimeZone: time.UTC,
    }, nil)
    defer client.Close(ctx)

    // Add knowledge
    client.Add(ctx, []predicato.Episode{{
        ID:        "meeting-1",
        Name:      "Team Standup",
        Content:   "Alice mentioned the API redesign is blocked on the auth team.",
        Reference: time.Now(),
        CreatedAt: time.Now(),
        GroupID:   "my-app",
    }})

    // Search with reranking
    results, _ := client.Search(ctx, "API redesign status", nil)
    
    // Rerank for better relevance
    passages := make([]string, len(results.Nodes))
    for i, node := range results.Nodes {
        passages[i] = node.Summary
    }
    ranked, _ := reranker.Rank(ctx, "API redesign status", passages)
    
    log.Printf("Top result: %s (score: %.2f)", ranked[0].Passage, ranked[0].Score)
}
```

First run downloads models (~1.7GB total). Subsequent runs use cached models.

## Quick Start (External APIs)

The same interfaces work with cloud services - just swap the implementations:

```go
package main

import (
    "context"
    "log"
    "os"
    "time"

    "github.com/soundprediction/predicato"
    "github.com/soundprediction/predicato/pkg/driver"
    "github.com/soundprediction/predicato/pkg/embedder"
    "github.com/soundprediction/predicato/pkg/nlp"
)

func main() {
    ctx := context.Background()

    // Neo4j database
    db, _ := driver.NewNeo4jDriver(
        os.Getenv("NEO4J_URI"),
        os.Getenv("NEO4J_USER"),
        os.Getenv("NEO4J_PASSWORD"),
    )
    defer db.Close(ctx)

    // OpenAI for LLM and embeddings
    apiKey := os.Getenv("OPENAI_API_KEY")
    llmClient, _ := nlp.NewOpenAIClient(apiKey, nlp.Config{Model: "gpt-4o-mini"})
    embedderClient := embedder.NewOpenAIEmbedder(apiKey, embedder.Config{
        Model: "text-embedding-3-small",
    })

    client, _ := predicato.NewClient(db, llmClient, embedderClient, &predicato.Config{
        GroupID:  "my-app",
        TimeZone: time.UTC,
    }, nil)
    defer client.Close(ctx)

    // Same API as internal stack
    client.Add(ctx, []predicato.Episode{{
        ID:      "meeting-1",
        Content: "Alice mentioned the API redesign is blocked.",
        // ...
    }})
}
```

## Installation

```bash
go get github.com/soundprediction/predicato
```

For internal services (requires CGO):

```bash
# Download native libraries
go generate ./...

# Build with Ladybug support
go build -tags system_ladybug ./...
```

### Prerequisites

**Internal stack:**
- Go 1.24+
- CGO enabled (`go env CGO_ENABLED` should return `1`)
- ~4GB RAM for local models

**External APIs:**
- Go 1.24+
- API keys for your chosen providers

## Examples

| Example | Description |
|---------|-------------|
| [`examples/basic/`](examples/basic/) | Full internal stack - Ladybug + RustBert + EmbedEverything + Reranking |
| [`examples/chat/`](examples/chat/) | Interactive chat with local models |
| [`examples/external_apis/`](examples/external_apis/) | Neo4j + OpenAI integration |

## Architecture

```
predicato/
├── pkg/driver/        # Graph databases (Ladybug, Neo4j, Memgraph)
├── pkg/embedder/      # Embedding providers (EmbedEverything, OpenAI, Gemini)
├── pkg/crossencoder/  # Reranking (EmbedEverything, Jina, LLM-based)
├── pkg/rustbert/      # Local text generation (GPT-2, NER, summarization)
├── pkg/nlp/           # LLM clients (OpenAI-compatible APIs)
├── pkg/search/        # Hybrid search (semantic + BM25 + graph traversal)
├── pkg/factstore/     # Versioned fact storage (PostgreSQL/DoltGres + pgvector)
└── pkg/types/         # Core types (nodes, edges, episodes)
```

## Internal Services Stack

| Component | Model | Download Size |
|-----------|-------|---------------|
| Embeddings | `qwen/qwen3-embedding-0.6b` | ~600MB |
| Reranking | `zhiqing/Qwen3-Reranker-0.6B-ONNX` | ~600MB |
| Text Generation | GPT-2 | ~500MB |

Models download automatically on first use and cache to `~/.cache/huggingface/`.

## Key Features

**Temporal Knowledge Graph**
- Bi-temporal model: `created_at` (when recorded) vs `valid_from/valid_to` (when true)
- Automatic invalidation of contradicting facts
- Historical queries: "what did we know about X as of date Y?"

**Hybrid Search**
- Semantic similarity (cosine distance on embeddings)
- BM25 keyword matching
- Graph traversal (BFS expansion through relationships)
- 5 reranking strategies: RRF, MMR, cross-encoder, node distance, episode mentions

**Production Ready**
- Circuit breakers with provider fallback
- Token usage tracking and cost calculation
- Error telemetry with DB persistence

## Fact Storage & RAG

Predicato includes a **fact storage system** for extracted entities and relationships that can be used independently for RAG (Retrieval-Augmented Generation) without requiring graph queries.

### PostgreSQL Backend

The fact storage uses PostgreSQL-compatible databases with **pgvector** for native vector similarity search:

| Mode | Database | Use Case |
|------|----------|----------|
| **Embedded** | DoltGres | Development, single-node deployment |
| **External** | PostgreSQL + pgvector | Production, managed databases (RDS, Cloud SQL) |

If no external PostgreSQL is configured, Predicato automatically uses **DoltGres** (embedded PostgreSQL-compatible database with git-like versioning).

### Configuration

```go
// Option 1: Automatic embedded DoltGres (no config needed)
client, _ := predicato.NewClient(db, llm, embedder, &predicato.Config{
    GroupID: "my-app",
}, nil)

// Option 2: External PostgreSQL with pgvector
client, _ := predicato.NewClient(db, llm, embedder, &predicato.Config{
    GroupID: "my-app",
    FactStoreConfig: &factstore.FactStoreConfig{
        Type:             "postgres",
        ConnectionString: "postgres://user:pass@localhost:5432/facts?sslmode=disable",
    },
}, nil)
```

### RAG Search (without Graph)

For simpler RAG use cases that don't need relationship traversal:

```go
// Search extracted facts directly (no graph queries)
results, _ := client.SearchFacts(ctx, "API design patterns", &types.SearchConfig{
    Limit:    10,
    MinScore: 0.7,
})

for _, node := range results.Nodes {
    fmt.Printf("Found: %s (score: %.2f)\n", node.Name, node.Score)
}
```

This performs hybrid search (vector similarity + keyword matching) using pgvector and PostgreSQL full-text search.

## CLI & Server

```bash
# Build CLI
make build-cli

# Start HTTP server
./bin/predicato server --port 8080

# API endpoints
POST /api/v1/ingest/messages  # Add content
POST /api/v1/search           # Search knowledge graph
GET  /api/v1/episodes/:id     # Get episodes
```

## Documentation

- [Getting Started](docs/GETTING_STARTED.md)
- [API Reference](docs/API_REFERENCE.md)
- [Ladybug Setup](docs/ladybug_SETUP.md)
- [FAQ](docs/FAQ.md)

## License

Apache 2.0

## Acknowledgments

Inspired by [Graphiti](https://github.com/getzep/graphiti) by Zep.

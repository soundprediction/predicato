# predicato

A temporal knowledge graph library for Go with a fully local ML stack - no API keys required.

## What Makes Predicato Different

Most agentic memory libraries require external services (OpenAI, Pinecone, Neo4j). Predicato is **modular by design** - every component can run locally OR connect to external services. Start with the internal stack for development, then swap in cloud services for production without changing your code.

| Component | Internal (No API) | External (Cloud) |
|-----------|-----------------|---------------------|
| **Graph Database** | Ladybug (embedded) | Neo4j, Memgraph |
| **Embeddings** | go-embedeverything | OpenAI, Voyage, Gemini |
| **Reranking** | go-embedeverything | Jina, Cohere |
| **Text Generation** | go-rust-bert (GPT-2) | OpenAI, Anthropic, Ollama |
| **Entity Extraction** | GLiNER (ONNX) | LLM-based extraction |
| **Fact Storage** | DoltDB (embedded) | DoltDB server |

**Why choose Predicato:**
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
├── pkg/factstore/     # Versioned fact storage (DoltDB)
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

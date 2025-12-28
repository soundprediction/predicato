# go-predicato

A temporal knowledge graph library for Go that extracts, organizes, and queries predicate logic from unstructured content.

## Key Capabilities

* **LLM Integration**: Supports OpenAI-compatible APIs including OpenAI, Anthropic, Gemini, Together AI, Ollama, and vLLM.
* **Cost Tracking**: Token usage tracking and cost calculation with serverless pricing models.
* **Routing**: Provider fallback, circuit breaking, and configurable routing rules.
* **Storage Options**:
    * Embedded ladybugDB for in-process graph storage
    * BadgerDB caching layer for embeddings and LLM responses
    * TSV-based prompting format
* **Telemetry**: Error tracking with DuckDB persistence.



## Features

- **Temporal Knowledge Graphs**: Bi-temporal data model tracking when data was recorded and when it is valid
- **Hybrid Search**: Combines semantic embeddings, keyword search (BM25), and graph traversal
- **Graph Database Support**: Embedded ladybugDB, Memgraph, and Neo4j
- **LLM Compatibility**: Works with OpenAI-compatible APIs (OpenAI, Anthropic, Gemini, Ollama, LocalAI, vLLM)
- **Optional Dependencies**: Can run with embedded database and without LLM features
- **CLI Tool**: Command-line interface for server management and graph operations
- **HTTP Server**: REST API endpoints for integration
- **MCP Protocol**: Model Context Protocol support for Claude Desktop and MCP clients
- **Cross-Encoder Reranking**: Multiple reranking implementations (Jina API, embedding similarity, LLM-based)


## Installation

```bash
go get github.com/soundprediction/go-predicato
```

**Note:** If building from source, you must run `go generate` to download the ladybug library and use the `system_ladybug` build tag:

```bash
go generate ./...
go build -tags system_ladybug ./...
```


## Quick Start

### Prerequisites

- Go 1.24+
- **Optional**: Graph database (ladybug embedded by default, or external Memgraph/Neo4j)
- **Optional**: LLM API access (OpenAI, Ollama, vLLM, or any OpenAI-compatible service)

### Environment Variables

**Basic Setup (Local/Embedded):**
```bash
# No environment variables required for basic usage with ladybug embedded database
# and without LLM features
```

**With OpenAI-compatible LLM (optional):**
```bash
export OPENAI_API_KEY="your-api-key"           # For OpenAI
export LLM_BASE_URL="http://localhost:11434"   # For local LLMs like Ollama
```

**With External Graph Database (optional):**
```bash
# For Memgraph
export MEMGRAPH_URI="bolt://localhost:7687"
export MEMGRAPH_USER="memgraph"
export MEMGRAPH_PASSWORD="your-password"

# For Neo4j
export NEO4J_URI="bolt://localhost:7687"
export NEO4J_USER="neo4j"
export NEO4J_PASSWORD="your-neo4j-password"

# Or for embedded ladybug (default)
export ladybug_DB_PATH="./ladybug_db"  # Optional: defaults to "./ladybug_db"
```

### Basic Usage

**Basic Example (ladybug + No LLM):**

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/soundprediction/go-predicato"
    "github.com/soundprediction/go-predicato/pkg/driver"
)

func main() {
    ctx := context.Background()

    // Create ladybug driver (embedded database)
    ladybugDriver, err := driver.NewLadybugDriver("./ladybug_db")
    if err != nil {
        log.Fatal(err)
    }
    defer ladybugDriver.Close(ctx)

    // Create Predicato client (LLM and embedder are optional)
    config := &predicato.Config{
        GroupID:  "my-group",
        TimeZone: time.UTC,
    }
    client := predicato.NewClient(ladybugDriver, nil, nil, config)
    defer client.Close(ctx)

    // Add episodes
    episodes := []predicato.Episode{
        {
            ID:        "meeting-1",
            Name:      "Team Meeting",
            Content:   "Discussed project timeline and resource allocation",
            Reference: time.Now(),
            CreatedAt: time.Now(),
            GroupID:   "my-group",
        },
    }

    err = client.Add(ctx, episodes)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Episode added to knowledge graph")
}
```

**With OpenAI-Compatible LLM:**

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/soundprediction/go-predicato"
    "github.com/soundprediction/go-predicato/pkg/driver"
    "github.com/soundprediction/go-predicato/pkg/embedder"
    "github.com/soundprediction/go-predicato/pkg/llm"
)

func main() {
    ctx := context.Background()

    // Create ladybug driver (embedded database)
    ladybugDriver, err := driver.NewLadybugDriver("./ladybug_db")
    if err != nil {
        log.Fatal(err)
    }
    defer ladybugDriver.Close(ctx)

    // Create LLM client (works with any OpenAI-compatible API)
    llmConfig := llm.Config{
        Model:       "gpt-4o-mini",  // Or "llama3", "mistral", etc.
        Temperature: &[]float32{0.7}[0],
        BaseURL:     "http://localhost:11434",  // Optional: for local LLMs
    }
    llmClient, err := llm.NewOpenAIClient("your-api-key", llmConfig)
    if err != nil {
        log.Fatal(err)
    }

    // Create embedder (optional, but recommended for semantic search)
    embedderConfig := embedder.Config{
        Model:     "text-embedding-3-small",  // Or local embedding model
        BaseURL:   "http://localhost:11434",  // Optional: for local embeddings
    }
    embedderClient := embedder.NewOpenAIEmbedder("your-api-key", embedderConfig)

    // Create Predicato client
    config := &predicato.Config{
        GroupID:  "my-group",
        TimeZone: time.UTC,
    }
    client := predicato.NewClient(ladybugDriver, llmClient, embedderClient, config)
    defer client.Close(ctx)

    // Add episodes
    episodes := []predicato.Episode{
        {
            ID:        "meeting-1",
            Name:      "Team Meeting",
            Content:   "Discussed project timeline and resource allocation",
            Reference: time.Now(),
            CreatedAt: time.Now(),
            GroupID:   "my-group",
        },
    }

    err = client.Add(ctx, episodes)
    if err != nil {
        log.Fatal(err)
    }

    // Search the knowledge graph (requires embedder for semantic search)
    results, err := client.Search(ctx, "project timeline", nil)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Found %d nodes", len(results.Nodes))
}
```

## CLI Tool

Go-Predicato includes a command-line interface for managing the knowledge graph and running servers.

### Installation

```bash
# Build from source using Makefile (recommended)
make build-cli

# Or build manually
go generate ./cmd/main.go
go build -tags system_ladybug -o bin/predicato ./cmd/main.go
```

### Server Command

Start the HTTP server:

```bash
./bin/predicato server
```

With custom configuration:

```bash
./bin/predicato server --port 9090 --llm-api-key your-key-here
```

### Configuration

Create a configuration file:

```bash
cp .predicato.example.yaml .predicato.yaml
# Edit the configuration as needed
```

The server provides REST API endpoints:

- `GET /health` - Health check
- `POST /api/v1/ingest/messages` - Add messages to knowledge graph
- `POST /api/v1/search` - Search the knowledge graph
- `GET /api/v1/episodes/:group_id` - Get episodes for a group
- `POST /api/v1/get-memory` - Get memory based on messages

### API Examples

Add messages:
```bash
curl -X POST http://localhost:8080/api/v1/ingest/messages \
  -H "Content-Type: application/json" \
  -d '{
    "group_id": "user123",
    "messages": [{"role": "user", "content": "Hello, I work at Acme Corp"}]
  }'
```

Search:
```bash
curl -X POST http://localhost:8080/api/v1/search \
  -H "Content-Type: application/json" \
  -d '{
    "query": "Acme Corp",
    "group_ids": ["user123"],
    "max_facts": 10
  }'
```

See [cmd/README.md](cmd/README.md) for detailed CLI documentation.

## Architecture

The library is structured into several key packages:

- **`predicato.go`**: Main client interface and configuration
- **`pkg/driver/`**: Graph database drivers (ladybug, Memgraph, Neo4j)
- **`pkg/llm/`**: Language model clients (OpenAI-compatible APIs)
- **`pkg/embedder/`**: Embedding model clients (OpenAI, Gemini, Voyage)
- **`pkg/search/`**: Hybrid search functionality
- **`pkg/types/`**: Core types for nodes, edges, and data structures
- **`pkg/models/`**: Database query builders for nodes and edges
- **`pkg/prompts/`**: LLM prompts for extraction and processing
- **`pkg/crossencoder/`**: Cross-encoder reranking for improved relevance
- **`pkg/community/`**: Community detection and management
- **`pkg/utils/`**: Utility functions for maintenance and operations

## Node Types

- **EntityNode**: Represents entities extracted from content
- **EpisodicNode**: Represents episodic memories or events  
- **CommunityNode**: Represents communities of related entities
- **SourceNode**: Represents source nodes where content originates

## Edge Types

- **EntityEdge**: Relationships between entities
- **EpisodicEdge**: Episodic relationships
- **CommunityEdge**: Community relationships

## Implementation Status

Implemented features:

- Entity and relationship extraction
- Node and edge deduplication
- Embedding generation and storage
- Hybrid search
- Community detection
- Temporal operations
- Bulk operations
- Error tracking and telemetry
- Cost calculation
- Router with provider fallback
- Caching layer (BadgerDB)
- Circuit breaker and email alerts

## Documentation

- [Getting Started](docs/GETTING_STARTED.md): Setup guide and first steps
- [Examples](docs/EXAMPLES.md): Usage examples
- [ladybug Setup Guide](docs/ladybug_SETUP.md): Embedded ladybug graph database configuration
- [FAQ](docs/FAQ.md): Common questions and troubleshooting

## Examples

See the `examples/` directory for usage examples:

- `examples/basic/`: Minimal setup with ladybug embedded database
- `examples/ladybug_ollama/`: Local setup with ladybug and Ollama
- `examples/openai_compatible/`: OpenAI-compatible service integration
- `examples/chat/`: Chat interface implementation
- `examples/prompts/`: Prompt engineering
- Additional examples in [docs/EXAMPLES.md](docs/EXAMPLES.md)

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build ./...
```

### Running Examples

```bash
# Basic example (no external dependencies)
cd examples/basic
go generate ./...
go run -tags system_ladybug main.go

# Or with local LLM
cd examples/ladybug_ollama
go generate ./...
go run -tags system_ladybug main.go

# Chat interface example
cd examples/chat
go generate ./...
go run -tags system_ladybug main.go
```


## License

Apache 2.0

## Acknowledgments

- This package takes inspiration from the original [Graphiti](https://github.com/getzep/graphiti) Python library by Zep


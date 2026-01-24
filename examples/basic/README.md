# Basic Predicato Example (Internal Services Stack)

This example demonstrates using predicato with **only internal services** - no external APIs or servers required!

## Internal Services Stack

| Component | Service | Model |
|-----------|---------|-------|
| **Database** | Ladybug | Embedded graph database |
| **Embeddings** | go-embedeverything | `qwen/qwen3-embedding-0.6b` |
| **Reranking** | go-embedeverything | `qwen/qwen3-reranker-0.6b` |
| **Text Generation** | go-rust-bert | GPT-2 |

## Features

This example shows how to:
- Create and configure a Predicato client using only internal services
- Use Ladybug embedded database (no external database server)
- Use go-rust-bert GPT-2 for text generation (no external LLM API)
- Use go-embedeverything with qwen3-embedding for embeddings (no external API)
- Use go-embedeverything with qwen3-reranker for reranking search results
- Add episodes (data) to the knowledge graph
- Search and rerank results from the knowledge graph

**No API keys or external services required!**

## Prerequisites

### Required
- Go 1.21 or later
- CGO enabled (required for Rust FFI bindings)
- ~4GB RAM minimum

### CGO Setup

CGO is required for the Rust-based ML libraries. Ensure you have:

**macOS:**
```bash
# Xcode command line tools (includes clang)
xcode-select --install
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt-get install build-essential
```

**Linux (Fedora/RHEL):**
```bash
sudo dnf install gcc gcc-c++
```

**Windows:**
Install MinGW-w64 or use WSL2 with Linux instructions.

### Verify CGO is enabled
```bash
go env CGO_ENABLED
# Should output: 1
```

If CGO is disabled, enable it:
```bash
export CGO_ENABLED=1
```

## First Run - Model Downloads

On first run, the example will automatically download the required models:

| Model | Size | Purpose |
|-------|------|---------|
| `qwen/qwen3-embedding-0.6b` | ~600MB | Text embeddings |
| `qwen/qwen3-reranker-0.6b` | ~600MB | Result reranking |
| GPT-2 | ~500MB | Text generation |

**Total: ~1.7GB**

Models are cached in `~/.cache/huggingface/` after the first download.

## Usage

1. Navigate to the example directory:
   ```bash
   cd examples/basic
   ```

2. Download the Ladybug native library (first time only):
   ```bash
   go generate ./...
   ```

3. Run the example:
   ```bash
   go run .
   ```

4. Or build and run:
   ```bash
   go build -o basic_example .
   ./basic_example
   ```

## Example Output

```
================================================================================
Predicato Basic Example - Internal Services Stack
================================================================================

This example uses predicato's internal services:
  - Ladybug: embedded graph database (no server required)
  - RustBert GPT-2: local text generation (no API required)
  - EmbedEverything: local embeddings with qwen/qwen3-embedding-0.6b
  - EmbedEverything: local reranking with qwen/qwen3-reranker-0.6b

No API keys or external services needed!

[1/5] Setting up Ladybug embedded graph database...
      Ladybug driver created (embedded database at ./example_graph.db)
[2/5] Setting up RustBert GPT-2 for text generation...
      RustBert GPT-2 text generation model loaded
[3/5] Setting up EmbedEverything embedder with qwen/qwen3-embedding-0.6b...
      EmbedEverything embedder created (model: qwen/qwen3-embedding-0.6b)
[4/5] Setting up EmbedEverything reranker with qwen/qwen3-reranker-0.6b...
      EmbedEverything reranker created (model: qwen/qwen3-reranker-0.6b)
[5/5] Creating Predicato client...
      Predicato client created (group: example-group)

================================================================================
All components initialized successfully!
================================================================================

Adding sample episodes to the knowledge graph...
Added 3 episodes to the knowledge graph

Searching the knowledge graph for: "API design and deadlines"
Found 5 nodes and 3 edges

Search results (before reranking):
----------------------------------
  1. Meeting with Alice (episode)
     Had a productive meeting with Alice about the new project...
  2. Project Research (episode)
     Researched various approaches for implementing the API...

Reranking results with qwen/qwen3-reranker-0.6b...

Search results (after reranking):
---------------------------------
  1. (score: 0.892) Had a productive meeting with Alice about the deadline...
  2. (score: 0.756) Researched various approaches for implementing the API...

Demonstrating text generation with RustBert GPT-2...
Prompt: The advantages of using a knowledge graph are
Generated: The advantages of using a knowledge graph are numerous and...

================================================================================
Example completed successfully!
================================================================================

Summary:
  - Used Ladybug embedded database (no Neo4j server)
  - Used RustBert GPT-2 for text generation (no OpenAI API)
  - Used qwen/qwen3-embedding-0.6b for embeddings (no API)
  - Used qwen/qwen3-reranker-0.6b for reranking (no API)

For external API examples, see: examples/external_apis/
```

## Troubleshooting

### CGO not enabled
```
# Error: undefined: embedder.NewEmbedEverythingClient
export CGO_ENABLED=1
go run .
```

### Model download fails
```
# Check internet connection and retry
# Models are downloaded from Hugging Face Hub

# If behind a proxy:
export HTTP_PROXY=http://proxy:port
export HTTPS_PROXY=http://proxy:port
```

### Out of memory
```
# The example requires ~4GB RAM
# Close other applications or use a machine with more memory

# Alternatively, use the external_apis example which offloads
# ML to cloud services
```

### Slow first run
First run downloads ~1.7GB of models. Subsequent runs use cached models and start much faster.

## What This Example Demonstrates

1. **Zero External Dependencies**: No API keys, no database servers, no cloud services
2. **Embedded Database**: Ladybug stores the knowledge graph locally
3. **Local ML**: All embeddings, reranking, and text generation run locally
4. **Privacy**: All data stays on your machine
5. **Reranking**: Demonstrates how reranking improves search result relevance

## Alternative: External APIs

For production deployments or if you prefer cloud services, see:
- `examples/external_apis/` - Uses Neo4j and OpenAI

## Related Examples

- **[External APIs Example](../external_apis/)**: Neo4j + OpenAI setup
- **[Chat Example](../chat/)**: Interactive chat with internal services

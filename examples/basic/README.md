# Basic Example - Internal Services Stack

This example demonstrates Predicato's fully local setup using only internal services - **no API keys or external services required**.

## What This Example Shows

- Creating a Predicato client with the internal services stack
- Using Ladybug embedded database (no server required)
- Using RustBert GPT-2 for text generation (local, no API)
- Using EmbedEverything for embeddings (qwen/qwen3-embedding-0.6b)
- Using EmbedEverything for reranking (zhiqing/Qwen3-Reranker-0.6B-ONNX)
- Adding episodes to the knowledge graph
- Searching and reranking results

## Prerequisites

- **Go 1.21+**
- **GCC** (for CGO compilation)
- **~4GB RAM** minimum (for loading ML models)
- **~2GB disk space** (for model downloads)

No API keys needed!

## Model Downloads

First run will automatically download models to `~/.cache/huggingface/`:

| Component | Model | Size |
|-----------|-------|------|
| Embeddings | qwen/qwen3-embedding-0.6b | ~600MB |
| Reranking | zhiqing/Qwen3-Reranker-0.6B-ONNX | ~600MB |
| Text Generation | GPT-2 | ~500MB |

## Build & Run

```bash
# From the repository root
cd examples/basic

# Download native library (first time only)
go generate

# Build
go build -o basic_example .

# Run
./basic_example
```

Or use Make from the repository root:

```bash
make build
cd examples/basic
go run .
```

## Expected Output

```
================================================================================
Predicato Basic Example - Internal Services Stack
================================================================================

This example uses predicato's internal services:
  - Ladybug: embedded graph database (no server required)
  - RustBert GPT-2: local text generation (no API required)
  - EmbedEverything: local embeddings with qwen/qwen3-embedding-0.6b
  - EmbedEverything: local reranking with zhiqing/Qwen3-Reranker-0.6B-ONNX

No API keys or external services needed!

[1/5] Setting up Ladybug embedded graph database...
      Ladybug driver created (embedded database at ./example_graph.db)
[2/5] Setting up RustBert GPT-2 for text generation...
      RustBert GPT-2 text generation model loaded
[3/5] Setting up EmbedEverything embedder with qwen/qwen3-embedding-0.6b...
      EmbedEverything embedder created
[4/5] Setting up EmbedEverything reranker with zhiqing/Qwen3-Reranker-0.6B-ONNX...
      EmbedEverything reranker created
[5/5] Creating Predicato client...
      Predicato client created (group: example-group)

================================================================================
All components initialized successfully!
================================================================================

Adding sample episodes to the knowledge graph...
Added 3 episodes to the knowledge graph

Searching the knowledge graph for: "API design and deadlines"
Found 3 nodes and 5 edges

Search results (before reranking):
----------------------------------
  1. Meeting with Alice (episodic)
     Had a productive meeting with Alice about the new project...
  2. Project Research (episodic)
     Researched various approaches for implementing the API...

Reranking results with zhiqing/Qwen3-Reranker-0.6B-ONNX...

Search results (after reranking):
---------------------------------
  1. (score: 0.892) Had a productive meeting with Alice about the new project...
  2. (score: 0.756) Researched various approaches for implementing the API...

Demonstrating text generation with RustBert GPT-2...
Prompt: The advantages of using a knowledge graph are
Generated: that it can be used to represent the relationships between...

================================================================================
Example completed successfully!
================================================================================
```

## Files Created

After running, you'll see:
- `./example_graph.db` - Ladybug database directory

## Troubleshooting

### "cannot find -llbug"

The Ladybug native library hasn't been downloaded:

```bash
go generate
```

### CGO errors

Ensure GCC is installed:

```bash
# Ubuntu/Debian
sudo apt install build-essential

# macOS
xcode-select --install
```

### Out of memory

The example requires ~4GB RAM. Close other applications or try a machine with more memory.

### Model download fails

Check your internet connection and disk space. Models download to `~/.cache/huggingface/`.

## Next Steps

- See `examples/chat/` for an interactive chat application
- See `examples/external_apis/` for using cloud services (OpenAI, Neo4j)
- Read the [Getting Started Guide](../../docs/GETTING_STARTED.md)

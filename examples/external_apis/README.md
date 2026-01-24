# External APIs Example (Cloud and Server-Based Services)

This example demonstrates how to use predicato with **external APIs and services** for production deployments and cloud-scale applications.

> **Note**: For a fully local setup with no external dependencies, see the [basic example](../basic/) which uses only internal services.

## When to Use External APIs

Use this approach when you need:
- **Cloud-scale databases**: Neo4j for shared, distributed knowledge graphs
- **High-quality LLM responses**: OpenAI GPT-4 or other cloud LLM services
- **GPU-accelerated inference**: vLLM or other inference servers
- **Team collaboration**: Shared knowledge graphs across multiple users/services
- **Managed services**: Reduced operational overhead

## Supported Services

### LLM Providers (OpenAI-Compatible API)

| Service | Description | API Key Required |
|---------|-------------|------------------|
| **Ollama** | Local LLM server | No |
| **LocalAI** | Local LLM with model flexibility | No |
| **vLLM** | High-performance inference server | No |
| **OpenAI** | Cloud LLM service | Yes |

### Database

| Service | Description |
|---------|-------------|
| **Neo4j** | External graph database server |
| **Memgraph** | Neo4j-compatible graph database |

### Embeddings

| Service | Description | API Key Required |
|---------|-------------|------------------|
| **OpenAI** | Cloud embedding service | Yes |

## Prerequisites

### For Ollama Examples
```bash
# Install Ollama (macOS)
brew install ollama

# Install Ollama (Linux)
curl -fsSL https://ollama.ai/install.sh | sh

# Start Ollama server
ollama serve

# Pull a model
ollama pull llama2:7b
```

### For Neo4j Examples
```bash
# Run Neo4j in Docker
docker run --name neo4j \
  -p 7687:7687 -p 7474:7474 \
  -e NEO4J_AUTH=neo4j/password \
  neo4j:latest

# Set environment variable
export NEO4J_PASSWORD=password
```

### For OpenAI Examples
```bash
export OPENAI_API_KEY=your_api_key_here
```

## Usage

1. Navigate to the example directory:
   ```bash
   cd examples/external_apis
   ```

2. Download the Ladybug native library (first time only):
   ```bash
   go generate ./...
   ```

3. Set required environment variables (see Prerequisites above)

4. Run the example:
   ```bash
   go run .
   ```

## Example Output

```
OpenAI-Compatible Client Examples
=================================

This example demonstrates how to use predicato with various
OpenAI-compatible services. Make sure you have the following
services running:

1. Ollama: Install and run 'ollama serve', then 'ollama pull llama2:7b'
2. LocalAI: Run LocalAI server on http://localhost:8080
3. vLLM: Run vLLM server on the specified URL
4. Neo4j: Required for full Predicato integration

Set these environment variables:
- NEO4J_URI (default: bolt://localhost:7687)
- NEO4J_USER (default: neo4j)
- NEO4J_PASSWORD (required for Predicato integration)
- OPENAI_API_KEY (optional, for embeddings)

=== Ollama Example ===
Creating Ollama client...
Ollama Response: A knowledge graph is a structured representation of...
Tokens used: 42

=== LocalAI Example ===
Creating LocalAI client...
LocalAI Response: Neo4j provides native graph storage and processing...

=== Full Predicato Integration Example ===
Creating Predicato client with Ollama LLM...
Adding episodes to knowledge graph...
Successfully processed episodes with local LLM!
Found 3 relevant nodes in knowledge graph
```

## Code Examples

### Ollama Integration

```go
client, err := nlp.NewOpenAIClient(
    "", // No API key needed for Ollama
    nlp.Config{
        BaseURL:     "http://localhost:11434",
        Model:       "llama2:7b",
        Temperature: &[]float32{0.7}[0],
        MaxTokens:   &[]int{1000}[0],
    },
)
```

### OpenAI Integration

```go
client, err := nlp.NewOpenAIClient(
    os.Getenv("OPENAI_API_KEY"),
    nlp.Config{
        Model:       "gpt-4o-mini",
        Temperature: &[]float32{0.7}[0],
        MaxTokens:   &[]int{2000}[0],
    },
)
```

### Neo4j + Predicato Integration

```go
// Create Neo4j driver
neo4jDriver, err := driver.NewNeo4jDriver(
    "bolt://localhost:7687",
    "neo4j",
    os.Getenv("NEO4J_PASSWORD"),
    "neo4j",
)

// Create embedder
embedderClient := embedder.NewOpenAIEmbedder(
    os.Getenv("OPENAI_API_KEY"),
    embedder.Config{
        Model:     "text-embedding-3-small",
        BatchSize: 100,
    },
)

// Create Predicato client
predicatoClient, err := predicato.NewClient(
    neo4jDriver,
    llmClient,
    embedderClient,
    config,
    nil,
)
```

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `NEO4J_URI` | No | `bolt://localhost:7687` | Neo4j connection URI |
| `NEO4J_USER` | No | `neo4j` | Neo4j username |
| `NEO4J_PASSWORD` | For Neo4j | - | Neo4j password |
| `OPENAI_API_KEY` | For OpenAI | - | OpenAI API key |

## Comparison: External APIs vs Internal Services

| Aspect | External APIs | Internal Services |
|--------|---------------|-------------------|
| **Setup** | Requires external services | No external dependencies |
| **Cost** | API usage fees | Free (CPU only) |
| **Quality** | Higher (GPT-4, etc.) | Good (GPT-2, qwen) |
| **Speed** | Depends on network | Depends on CPU |
| **Privacy** | Data sent to cloud | All data stays local |
| **Scalability** | High (cloud) | Limited (local machine) |

## Related Examples

- **[Basic Example](../basic/)**: Fully local with internal services
- **[Chat Example](../chat/)**: Interactive chat with internal services

# GLInER2 Integration

## Overview

This package adds GLInER2 support to Predicato through a unified client architecture that supports multiple providers:

- **Local API**: FastAPI service mirroring Fastino's `/gliner-2` endpoint
- **Fastino API**: Official Fastino hosted GLInER2 service  
- **Future Native**: Direct go-gline-rs GLInER2 when available

## Architecture

```
pkg/gliner2/                     # GLInER2 package (new)
├── client.go                    # Unified client with provider switching
├── config.go                    # Provider configurations
├── types.go                      # GLInER2 data structures  
├── http_client.go               # HTTP client for local/Fastino APIs
├── native_client.go             # Placeholder for future go-gline-rs GLInER2
└── utils.go                      # Helper functions

cmd/gliner2-api/               # Local FastAPI service (optional)
└── main.go                       # HTTP server implementing /gliner-2 endpoint
```

## Features

### Entity Extraction (NER)
- Custom entity types with descriptions
- Confidence scores and span positions
- Batch processing support

### Fact Extraction (Relations)  
- GLInER2 relation extraction mapped to fact triples
- Support for custom relation types
- Compatible with existing Predicato fact format

### Text Classification (Future)
- Single and multi-label classification
- Confidence scores
- Custom category definitions

### Structured Data Extraction (Future)
- JSON schema-driven extraction
- Field-level validation
- Nested object support

## Usage

### Basic Entity Extraction

```go
config := gliner2.Config{
    Provider: gliner2.ProviderLocal,
    Local: &gliner2.LocalConfig{
        Endpoint: "http://localhost:8000/gliner-2", 
        Timeout: 30 * time.Second,
    },
}

client, err := gliner2.NewClient(config)
if err != nil {
    log.Fatal(err)
}

entities, err := client.ExtractEntities(ctx, text, []string{"person", "company", "location"})
```

### Fastino API Usage

```go
config := gliner2.Config{
    Provider: gliner2.ProviderFastino,
    Fastino: &gliner2.FastinoConfig{
        Endpoint: "https://api.fastino.ai/gliner-2",
        APIKey:   os.Getenv("FASTINO_API_KEY"),
        Timeout:  30 * time.Second,
    },
}
```

### Local API Service

Run the local FastAPI service:

```bash
# Build the service
go build ./cmd/gliner2-api

# Run with default port 8000
./gliner2-api

# Or custom port
PORT=8080 ./gliner2-api
```

The local service mirrors Fastino's API exactly:

```
POST /gliner-2
{
  "task": "extract_entities",
  "text": "Apple CEO Tim Cook announced iPhone 15 in Cupertino.", 
  "schema": ["person", "company", "product", "location"],
  "threshold": 0.5
}

Response:
{
  "result": {
    "entities": {
      "person": [{"text": "Tim Cook", "confidence": 0.92}],
      "company": [{"text": "Apple", "confidence": 0.95}],
      "product": [{"text": "iPhone 15", "confidence": 0.88}], 
      "location": [{"text": "Cupertino", "confidence": 0.90}]
    }
  }
}
```

## Configuration Examples

### Local Development
```yaml
gliner2:
  provider: "local"
  local:
    endpoint: "http://localhost:8000/gliner-2"
    timeout: "30s"
```

### Production (Fastino)
```yaml
gliner2:
  provider: "fastino"  
  fastino:
    endpoint: "https://api.fastino.ai/gliner-2"
    api_key: "${FASTINO_API_KEY}"
    timeout: "30s"
```

## Integration with Predicato

The GLInER2 package integrates with existing NLP pipelines through the standard `nlp.Client` interface, supporting entity and fact extraction prompts that match the existing GLInER adapter patterns.

### Registry Support

GLInER2 models are registered in the NLP registry:
- `fastino/gliner2-base-v1` - Base GLInER2 model
- `fastino/gliner2-large-v1` - Large GLInER2 model

Both models support:
- TaskNamedEntityRecognition  
- TaskRelationExtraction (fact extraction)

## Future Native Support

When go-gline-rs adds GLInER2 support, the native client will enable:
- Zero HTTP overhead and maximum performance
- Direct Go bindings without Python dependencies
- Same unified provider switching interface

The `native_client.go` file contains the placeholder structure for this future integration.
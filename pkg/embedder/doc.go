// Package embedder provides text embedding clients for vector representations.
//
// This package defines the Client interface and provides implementations for
// various embedding providers including OpenAI and local embedding services.
//
// # Supported Providers
//
// The following embedding providers are supported:
//   - OpenAI: text-embedding-3-small, text-embedding-3-large, text-embedding-ada-002
//   - Local: Local embedding services via HTTP API
//
// # Usage
//
//	// Create an OpenAI embedder
//	embedder := embedder.NewOpenAIEmbedder(apiKey, embedder.Config{
//	    Model:     "text-embedding-3-small",
//	    BatchSize: 100,
//	})
//
//	// Embed text
//	embeddings, err := embedder.Embed(ctx, []string{"hello world"})
//
// # Batch Processing
//
// The Client interface supports batch embedding for efficiency:
//   - Embed(): Embed multiple texts in a single request
//   - EmbedSingle(): Convenience method for single text
//
// Implementations handle batching internally based on provider limits.
package embedder

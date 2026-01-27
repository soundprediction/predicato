// Package nlp provides natural language processing clients for LLM interactions.
//
// This package defines the Client interface and provides implementations for
// various LLM providers including OpenAI, Anthropic, Gemini, and OpenAI-compatible
// APIs (Ollama, vLLM, etc.).
//
// # Supported Providers
//
// The following LLM providers are supported:
//   - OpenAI: GPT-4, GPT-3.5, and other OpenAI models
//   - Anthropic: Claude models
//   - Gemini: Google's Gemini models
//   - OpenAI-compatible: Any API following OpenAI's format (Ollama, vLLM, etc.)
//
// # Client Wrappers
//
// The package provides several wrapper clients for enhanced functionality:
//   - RetryClient: Automatic retry with exponential backoff and jitter
//   - TokenTrackingClient: Track token usage across requests
//   - CircuitBreakerClient: Circuit breaker pattern for fault tolerance
//   - RouterClient: Route requests to different providers based on criteria
//
// # Usage
//
//	// Create a base client
//	client, err := nlp.NewOpenAIClient(apiKey, config)
//
//	// Wrap with retry logic
//	retryClient, err := nlp.NewRetryClient(client, nlp.DefaultRetryConfig())
//
//	// Use the client
//	response, err := retryClient.Generate(ctx, prompt, nil)
//
// # Error Handling
//
// The package defines specific error types for common failure modes:
//   - RateLimitError: API rate limit exceeded
//   - RefusalError: Model refused to generate content
//   - EmptyResponseError: Model returned empty response
//
// These errors support errors.Is() for type checking.
package nlp

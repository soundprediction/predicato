package crossencoder

import (
	"context"
	"testing"
)

func TestMockRerankerClient(t *testing.T) {
	client := NewMockRerankerClient(DefaultConfig(ProviderMock))
	defer client.Close()

	ctx := context.Background()
	query := "artificial intelligence machine learning"
	passages := []string{
		"Machine learning is a subset of artificial intelligence",
		"The weather is nice today",
		"Deep learning models use neural networks",
		"Cats are cute animals",
		"AI and ML are transforming technology",
	}

	results, err := client.Rank(ctx, query, passages)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(results) != len(passages) {
		t.Fatalf("Expected %d results, got %d", len(passages), len(results))
	}

	// Verify results are sorted by score (descending)
	for i := 1; i < len(results); i++ {
		if results[i-1].Score < results[i].Score {
			t.Errorf("Results not sorted by score: %f < %f", results[i-1].Score, results[i].Score)
		}
	}

	// The first passage should likely rank highest due to keyword overlap
	if results[0].Passage != "Machine learning is a subset of artificial intelligence" &&
		results[0].Passage != "AI and ML are transforming technology" {
		t.Logf("Top result: %s (score: %f)", results[0].Passage, results[0].Score)
		t.Logf("All results:")
		for i, result := range results {
			t.Logf("  %d. %s (%.3f)", i+1, result.Passage, result.Score)
		}
	}
}

func TestLocalRerankerClient(t *testing.T) {
	client := NewLocalRerankerClient(DefaultConfig(ProviderLocal))
	defer client.Close()

	ctx := context.Background()
	query := "machine learning algorithms"
	passages := []string{
		"Machine learning algorithms are used in data science",
		"Cooking recipes for dinner tonight",
		"Neural networks and deep learning",
		"Sports scores and statistics",
		"Supervised learning algorithms like decision trees",
	}

	results, err := client.Rank(ctx, query, passages)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(results) != len(passages) {
		t.Fatalf("Expected %d results, got %d", len(passages), len(results))
	}

	// Verify results are sorted by score (descending)
	for i := 1; i < len(results); i++ {
		if results[i-1].Score < results[i].Score {
			t.Errorf("Results not sorted by score: %f < %f", results[i-1].Score, results[i].Score)
		}
	}
}

func TestEmptyPassages(t *testing.T) {
	client := NewMockRerankerClient(DefaultConfig(ProviderMock))
	defer client.Close()

	ctx := context.Background()
	results, err := client.Rank(ctx, "test query", []string{})
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("Expected 0 results for empty passages, got %d", len(results))
	}
}

func TestEmbedEverythingClient(t *testing.T) {
	// This test requires model downloads from Hugging Face and may fail if:
	// 1. No internet connection
	// 2. Model URL is not accessible
	// 3. Model format is not compatible
	// Skip if client creation fails
	config := &EmbedEverythingConfig{
		Config: &Config{
			Model: "BAAI/bge-reranker-base",
		},
	}

	client, err := NewEmbedEverythingClient(config)
	if err != nil {
		t.Skipf("Skipping EmbedEverything test: %v", err)
		return
	}
	defer client.Close()

	ctx := context.Background()
	query := "machine learning algorithms"
	passages := []string{
		"Machine learning algorithms are used in data science",
		"Cooking recipes for dinner tonight",
		"Neural networks and deep learning",
	}

	results, err := client.Rank(ctx, query, passages)
	if err != nil {
		t.Fatalf("Expected no error during ranking, got: %v", err)
	}

	if len(results) != len(passages) {
		t.Fatalf("Expected %d results, got %d", len(passages), len(results))
	}

	// Verify results are sorted by score (descending)
	for i := 1; i < len(results); i++ {
		if results[i-1].Score < results[i].Score {
			t.Errorf("Results not sorted by score: %f < %f", results[i-1].Score, results[i].Score)
		}
	}
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name         string
		config       ClientConfig
		expectError  bool
		expectedType string
	}{
		{
			name: "mock provider",
			config: ClientConfig{
				Provider: ProviderMock,
				Config:   DefaultConfig(ProviderMock),
			},
			expectError:  false,
			expectedType: "*crossencoder.MockRerankerClient",
		},
		{
			name: "local provider",
			config: ClientConfig{
				Provider: ProviderLocal,
				Config:   DefaultConfig(ProviderLocal),
			},
			expectError:  false,
			expectedType: "*crossencoder.LocalRerankerClient",
		},
		{
			name: "openai provider without llm client",
			config: ClientConfig{
				Provider: ProviderOpenAI,
				Config:   DefaultConfig(ProviderOpenAI),
			},
			expectError: true,
		},
		// Note: embedeverything provider test is skipped here as it requires model downloads
		// See TestEmbedEverythingClient for a dedicated test with skip logic
		{
			name: "unknown provider",
			config: ClientConfig{
				Provider: "unknown",
				Config:   Config{},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
				return
			}

			if client == nil {
				t.Errorf("Expected client, got nil")
				return
			}

			// Clean up
			client.Close()
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	tests := []struct {
		provider Provider
		expected Config
	}{
		{
			provider: ProviderOpenAI,
			expected: Config{
				Model:          "gpt-4o-mini",
				BatchSize:      10,
				MaxConcurrency: 5,
			},
		},
		{
			provider: ProviderLocal,
			expected: Config{
				BatchSize: 100,
			},
		},
		{
			provider: ProviderMock,
			expected: Config{
				BatchSize: 100,
			},
		},
		{
			provider: ProviderEmbedEverything,
			expected: Config{
				Model:          "BAAI/bge-reranker-base",
				BatchSize:      100,
				MaxConcurrency: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.provider), func(t *testing.T) {
			config := DefaultConfig(tt.provider)

			if config.Model != tt.expected.Model {
				t.Errorf("Expected model %s, got %s", tt.expected.Model, config.Model)
			}
			if config.BatchSize != tt.expected.BatchSize {
				t.Errorf("Expected batch size %d, got %d", tt.expected.BatchSize, config.BatchSize)
			}
			if config.MaxConcurrency != tt.expected.MaxConcurrency {
				t.Errorf("Expected max concurrency %d, got %d", tt.expected.MaxConcurrency, config.MaxConcurrency)
			}
		})
	}
}

func TestEmbedEverythingQwenReranker(t *testing.T) {
	// This test validates the qwen/qwen3-reranker-0.6b model which is the
	// recommended internal reranker for predicato examples.
	// It requires model downloads from Hugging Face and may be skipped if:
	// 1. No internet connection
	// 2. Model URL is not accessible
	// 3. Insufficient disk space (~600MB)
	config := &EmbedEverythingConfig{
		Config: &Config{
			Model:     "qwen/qwen3-reranker-0.6b",
			BatchSize: 100,
		},
	}

	client, err := NewEmbedEverythingClient(config)
	if err != nil {
		t.Skipf("Skipping qwen reranker test (model may need download): %v", err)
		return
	}
	defer client.Close()

	ctx := context.Background()

	t.Run("basic ranking", func(t *testing.T) {
		query := "What is machine learning?"
		passages := []string{
			"Machine learning is a branch of artificial intelligence that enables computers to learn from data.",
			"The weather forecast predicts rain tomorrow afternoon.",
			"Deep learning uses neural networks with many layers to learn complex patterns.",
			"Cooking pasta requires boiling water and adding salt.",
			"Supervised learning algorithms learn from labeled training examples.",
		}

		results, err := client.Rank(ctx, query, passages)
		if err != nil {
			t.Fatalf("Expected no error during ranking, got: %v", err)
		}

		if len(results) != len(passages) {
			t.Fatalf("Expected %d results, got %d", len(passages), len(results))
		}

		// Verify results are sorted by score (descending)
		for i := 1; i < len(results); i++ {
			if results[i-1].Score < results[i].Score {
				t.Errorf("Results not sorted by score: %f < %f", results[i-1].Score, results[i].Score)
			}
		}

		// Log results for debugging
		t.Logf("Query: %s", query)
		for i, result := range results {
			t.Logf("  %d. [%.4f] %s", i+1, result.Score, result.Passage)
		}

		// The ML-related passages should rank higher than unrelated ones
		// Check that the top result contains relevant keywords
		topResult := results[0].Passage
		if !containsAny(topResult, []string{"machine learning", "learning", "neural", "supervised"}) {
			t.Errorf("Expected top result to be ML-related, got: %s", topResult)
		}
	})

	t.Run("empty passages", func(t *testing.T) {
		results, err := client.Rank(ctx, "test query", []string{})
		if err != nil {
			t.Fatalf("Expected no error for empty passages, got: %v", err)
		}
		if len(results) != 0 {
			t.Fatalf("Expected 0 results for empty passages, got %d", len(results))
		}
	})

	t.Run("single passage", func(t *testing.T) {
		results, err := client.Rank(ctx, "test query", []string{"Single passage to rank"})
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if len(results) != 1 {
			t.Fatalf("Expected 1 result, got %d", len(results))
		}
	})

	t.Run("relevance ordering", func(t *testing.T) {
		// Test that clearly relevant passages rank higher than irrelevant ones
		query := "capital cities of European countries"
		passages := []string{
			"Paris is the capital of France and is known for the Eiffel Tower.",
			"The Amazon rainforest is the largest tropical rainforest in the world.",
			"Berlin is the capital and largest city of Germany.",
			"Pizza originated in Naples, Italy.",
			"London is the capital of the United Kingdom.",
		}

		results, err := client.Rank(ctx, query, passages)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Log results
		t.Logf("Query: %s", query)
		for i, result := range results {
			t.Logf("  %d. [%.4f] %s", i+1, result.Score, result.Passage)
		}

		// Count how many of the top 3 results are about capital cities
		capitalCount := 0
		for i := 0; i < 3 && i < len(results); i++ {
			if containsAny(results[i].Passage, []string{"capital", "Paris", "Berlin", "London"}) {
				capitalCount++
			}
		}
		if capitalCount < 2 {
			t.Errorf("Expected at least 2 capital city passages in top 3, got %d", capitalCount)
		}
	})
}

// containsAny checks if s contains any of the substrings
func containsAny(s string, substrings []string) bool {
	for _, sub := range substrings {
		if contains(s, sub) {
			return true
		}
	}
	return false
}

// contains checks if s contains substr (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && containsLower(toLower(s), toLower(substr))))
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

func containsLower(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Benchmark tests
func BenchmarkMockReranker(b *testing.B) {
	client := NewMockRerankerClient(DefaultConfig(ProviderMock))
	defer client.Close()

	ctx := context.Background()
	query := "machine learning artificial intelligence"
	passages := []string{
		"Machine learning is a subset of artificial intelligence",
		"Deep learning models use neural networks",
		"Natural language processing applications",
		"Computer vision and image recognition",
		"Reinforcement learning algorithms",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.Rank(ctx, query, passages)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

func BenchmarkLocalReranker(b *testing.B) {
	client := NewLocalRerankerClient(DefaultConfig(ProviderLocal))
	defer client.Close()

	ctx := context.Background()
	query := "machine learning artificial intelligence"
	passages := []string{
		"Machine learning is a subset of artificial intelligence",
		"Deep learning models use neural networks",
		"Natural language processing applications",
		"Computer vision and image recognition",
		"Reinforcement learning algorithms",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.Rank(ctx, query, passages)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

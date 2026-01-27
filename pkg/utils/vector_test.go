package utils

import (
	"math"
	"testing"
)

func TestCosineSimilarity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		a        []float32
		b        []float32
		expected float64
	}{
		{
			name:     "identical vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{1, 0, 0},
			expected: 1.0,
		},
		{
			name:     "opposite vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{-1, 0, 0},
			expected: -1.0,
		},
		{
			name:     "orthogonal vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{0, 1, 0},
			expected: 0.0,
		},
		{
			name:     "similar vectors",
			a:        []float32{1, 2, 3},
			b:        []float32{1, 2, 3},
			expected: 1.0,
		},
		{
			name:     "different lengths",
			a:        []float32{1, 2, 3},
			b:        []float32{1, 2},
			expected: 0.0,
		},
		{
			name:     "empty vectors",
			a:        []float32{},
			b:        []float32{},
			expected: 0.0,
		},
		{
			name:     "zero vector a",
			a:        []float32{0, 0, 0},
			b:        []float32{1, 2, 3},
			expected: 0.0,
		},
		{
			name:     "zero vector b",
			a:        []float32{1, 2, 3},
			b:        []float32{0, 0, 0},
			expected: 0.0,
		},
		{
			name:     "nil vectors",
			a:        nil,
			b:        nil,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CosineSimilarity(tt.a, tt.b)
			if math.Abs(result-tt.expected) > 1e-6 {
				t.Errorf("CosineSimilarity(%v, %v) = %v, expected %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestCosineSimilarity32(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		a        []float32
		b        []float32
		expected float32
	}{
		{
			name:     "identical vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{1, 0, 0},
			expected: 1.0,
		},
		{
			name:     "orthogonal vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{0, 1, 0},
			expected: 0.0,
		},
		{
			name:     "different lengths",
			a:        []float32{1, 2, 3},
			b:        []float32{1, 2},
			expected: 0.0,
		},
		{
			name:     "empty vectors",
			a:        []float32{},
			b:        []float32{},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CosineSimilarity32(tt.a, tt.b)
			if math.Abs(float64(result-tt.expected)) > 1e-6 {
				t.Errorf("CosineSimilarity32(%v, %v) = %v, expected %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestCosineSimilarity64(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		a        []float64
		b        []float64
		expected float64
	}{
		{
			name:     "identical vectors",
			a:        []float64{1, 0, 0},
			b:        []float64{1, 0, 0},
			expected: 1.0,
		},
		{
			name:     "orthogonal vectors",
			a:        []float64{1, 0, 0},
			b:        []float64{0, 1, 0},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CosineSimilarity64(tt.a, tt.b)
			if math.Abs(result-tt.expected) > 1e-6 {
				t.Errorf("CosineSimilarity64(%v, %v) = %v, expected %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestDotProduct(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		a        []float32
		b        []float32
		expected float64
	}{
		{
			name:     "simple dot product",
			a:        []float32{1, 2, 3},
			b:        []float32{4, 5, 6},
			expected: 32.0, // 1*4 + 2*5 + 3*6 = 4 + 10 + 18 = 32
		},
		{
			name:     "orthogonal vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{0, 1, 0},
			expected: 0.0,
		},
		{
			name:     "different lengths",
			a:        []float32{1, 2, 3},
			b:        []float32{1, 2},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DotProduct(tt.a, tt.b)
			if math.Abs(result-tt.expected) > 1e-6 {
				t.Errorf("DotProduct(%v, %v) = %v, expected %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestMagnitude(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		v        []float32
		expected float64
	}{
		{
			name:     "unit vector x",
			v:        []float32{1, 0, 0},
			expected: 1.0,
		},
		{
			name:     "3-4-5 triangle",
			v:        []float32{3, 4},
			expected: 5.0,
		},
		{
			name:     "zero vector",
			v:        []float32{0, 0, 0},
			expected: 0.0,
		},
		{
			name:     "empty vector",
			v:        []float32{},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Magnitude(tt.v)
			if math.Abs(result-tt.expected) > 1e-6 {
				t.Errorf("Magnitude(%v) = %v, expected %v", tt.v, result, tt.expected)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	t.Parallel()
	t.Run("normalize unit vector", func(t *testing.T) {
		v := []float32{1, 0, 0}
		result := Normalize(v)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if len(result) != 3 {
			t.Fatalf("expected length 3, got %d", len(result))
		}
		if math.Abs(float64(result[0])-1.0) > 1e-6 {
			t.Errorf("expected [1,0,0], got %v", result)
		}
	})

	t.Run("normalize non-unit vector", func(t *testing.T) {
		v := []float32{3, 4}
		result := Normalize(v)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		// Should be [0.6, 0.8]
		if math.Abs(float64(result[0])-0.6) > 1e-6 || math.Abs(float64(result[1])-0.8) > 1e-6 {
			t.Errorf("expected [0.6, 0.8], got %v", result)
		}
		// Magnitude should be 1
		mag := Magnitude(result)
		if math.Abs(mag-1.0) > 1e-6 {
			t.Errorf("expected magnitude 1.0, got %v", mag)
		}
	})

	t.Run("normalize zero vector", func(t *testing.T) {
		v := []float32{0, 0, 0}
		result := Normalize(v)
		if result != nil {
			t.Errorf("expected nil for zero vector, got %v", result)
		}
	})

	t.Run("normalize empty vector", func(t *testing.T) {
		v := []float32{}
		result := Normalize(v)
		if result != nil {
			t.Errorf("expected nil for empty vector, got %v", result)
		}
	})
}

func BenchmarkCosineSimilarity(b *testing.B) {
	// Create random-ish vectors of typical embedding size (1536 for OpenAI)
	a := make([]float32, 1536)
	bVec := make([]float32, 1536)
	for i := range a {
		a[i] = float32(i) / 1536.0
		bVec[i] = float32(1536-i) / 1536.0
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CosineSimilarity(a, bVec)
	}
}

func BenchmarkCosineSimilarity32(b *testing.B) {
	a := make([]float32, 1536)
	bVec := make([]float32, 1536)
	for i := range a {
		a[i] = float32(i) / 1536.0
		bVec[i] = float32(1536-i) / 1536.0
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CosineSimilarity32(a, bVec)
	}
}

func TestTopKByScore(t *testing.T) {
	t.Parallel()
	t.Run("basic top k", func(t *testing.T) {
		items := []ScoredItem[string]{
			{Item: "a", Score: 0.5},
			{Item: "b", Score: 0.9},
			{Item: "c", Score: 0.3},
			{Item: "d", Score: 0.7},
			{Item: "e", Score: 0.1},
		}

		result := TopKByScore(items, 3)
		if len(result) != 3 {
			t.Fatalf("expected 3 items, got %d", len(result))
		}

		// Should be sorted descending
		if result[0].Score != 0.9 || result[0].Item != "b" {
			t.Errorf("expected first item to be b with score 0.9, got %v", result[0])
		}
		if result[1].Score != 0.7 || result[1].Item != "d" {
			t.Errorf("expected second item to be d with score 0.7, got %v", result[1])
		}
		if result[2].Score != 0.5 || result[2].Item != "a" {
			t.Errorf("expected third item to be a with score 0.5, got %v", result[2])
		}
	})

	t.Run("k greater than length", func(t *testing.T) {
		items := []ScoredItem[int]{
			{Item: 1, Score: 0.5},
			{Item: 2, Score: 0.9},
		}

		result := TopKByScore(items, 10)
		if len(result) != 2 {
			t.Fatalf("expected 2 items, got %d", len(result))
		}
		if result[0].Score != 0.9 {
			t.Errorf("expected first score 0.9, got %f", result[0].Score)
		}
	})

	t.Run("k equals length", func(t *testing.T) {
		items := []ScoredItem[int]{
			{Item: 1, Score: 0.3},
			{Item: 2, Score: 0.9},
			{Item: 3, Score: 0.6},
		}

		result := TopKByScore(items, 3)
		if len(result) != 3 {
			t.Fatalf("expected 3 items, got %d", len(result))
		}
	})

	t.Run("k is zero", func(t *testing.T) {
		items := []ScoredItem[int]{
			{Item: 1, Score: 0.5},
		}

		result := TopKByScore(items, 0)
		if result != nil {
			t.Errorf("expected nil for k=0, got %v", result)
		}
	})

	t.Run("empty items", func(t *testing.T) {
		var items []ScoredItem[int]

		result := TopKByScore(items, 5)
		if result != nil {
			t.Errorf("expected nil for empty items, got %v", result)
		}
	})

	t.Run("k is one", func(t *testing.T) {
		items := []ScoredItem[string]{
			{Item: "low", Score: 0.1},
			{Item: "high", Score: 0.9},
			{Item: "mid", Score: 0.5},
		}

		result := TopKByScore(items, 1)
		if len(result) != 1 {
			t.Fatalf("expected 1 item, got %d", len(result))
		}
		if result[0].Item != "high" || result[0].Score != 0.9 {
			t.Errorf("expected high with 0.9, got %v", result[0])
		}
	})

	t.Run("duplicate scores", func(t *testing.T) {
		items := []ScoredItem[int]{
			{Item: 1, Score: 0.5},
			{Item: 2, Score: 0.5},
			{Item: 3, Score: 0.9},
			{Item: 4, Score: 0.5},
		}

		result := TopKByScore(items, 2)
		if len(result) != 2 {
			t.Fatalf("expected 2 items, got %d", len(result))
		}
		if result[0].Score != 0.9 {
			t.Errorf("expected first score 0.9, got %f", result[0].Score)
		}
		if result[1].Score != 0.5 {
			t.Errorf("expected second score 0.5, got %f", result[1].Score)
		}
	})
}

func TestTopKIndicesByScore(t *testing.T) {
	t.Parallel()
	t.Run("basic indices", func(t *testing.T) {
		scores := []float64{0.3, 0.9, 0.5, 0.1, 0.7}

		result := TopKIndicesByScore(scores, 3)
		if len(result) != 3 {
			t.Fatalf("expected 3 indices, got %d", len(result))
		}

		// Index 1 has score 0.9 (highest)
		// Index 4 has score 0.7 (second)
		// Index 2 has score 0.5 (third)
		if result[0] != 1 {
			t.Errorf("expected first index 1, got %d", result[0])
		}
		if result[1] != 4 {
			t.Errorf("expected second index 4, got %d", result[1])
		}
		if result[2] != 2 {
			t.Errorf("expected third index 2, got %d", result[2])
		}
	})

	t.Run("empty scores", func(t *testing.T) {
		result := TopKIndicesByScore([]float64{}, 5)
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})

	t.Run("k zero", func(t *testing.T) {
		result := TopKIndicesByScore([]float64{0.5, 0.3}, 0)
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func BenchmarkTopKByScore(b *testing.B) {
	// Simulate 10,000 items (typical for in-memory search)
	items := make([]ScoredItem[int], 10000)
	for i := range items {
		items[i] = ScoredItem[int]{
			Item:  i,
			Score: float64(i%1000) / 1000.0,
		}
	}

	b.Run("k=10", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			TopKByScore(items, 10)
		}
	})

	b.Run("k=100", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			TopKByScore(items, 100)
		}
	})

	b.Run("k=1000", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			TopKByScore(items, 1000)
		}
	})
}

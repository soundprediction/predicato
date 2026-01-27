package utils

import (
	"math"
	"testing"
)

func TestCosineSimilarity(t *testing.T) {
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

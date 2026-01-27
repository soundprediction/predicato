// Package utils provides common utility functions for the predicato project.
package utils

import "math"

// CosineSimilarity calculates the cosine similarity between two float32 vectors.
// Returns 0 if vectors have different lengths, are empty, or either has zero magnitude.
// The result is in the range [-1, 1], where 1 means identical direction,
// 0 means orthogonal, and -1 means opposite direction.
func CosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dotProduct, normA, normB float64

	for i := range a {
		dotProduct += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// CosineSimilarity32 calculates the cosine similarity between two float32 vectors,
// returning a float32 result. This is useful when working with float32 throughout.
// Returns 0 if vectors have different lengths, are empty, or either has zero magnitude.
func CosineSimilarity32(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dotProduct, normA, normB float32

	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}

// CosineSimilarity64 calculates the cosine similarity between two float64 vectors.
// Returns 0 if vectors have different lengths, are empty, or either has zero magnitude.
func CosineSimilarity64(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dotProduct, normA, normB float64

	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// DotProduct calculates the dot product of two float32 vectors.
// Returns 0 if vectors have different lengths.
func DotProduct(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}

	var result float64
	for i := range a {
		result += float64(a[i]) * float64(b[i])
	}
	return result
}

// Magnitude calculates the Euclidean magnitude (L2 norm) of a float32 vector.
func Magnitude(v []float32) float64 {
	var sum float64
	for _, x := range v {
		sum += float64(x) * float64(x)
	}
	return math.Sqrt(sum)
}

// Normalize normalizes a float32 vector to unit length.
// Returns nil if the input is empty or has zero magnitude.
func Normalize(v []float32) []float32 {
	if len(v) == 0 {
		return nil
	}

	mag := Magnitude(v)
	if mag == 0 {
		return nil
	}

	result := make([]float32, len(v))
	for i, x := range v {
		result[i] = float32(float64(x) / mag)
	}
	return result
}

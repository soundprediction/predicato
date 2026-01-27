// Package utils provides common utility functions for the predicato project.
package utils

import (
	"container/heap"
	"math"
)

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

// ScoredItem represents an item with a score for top-K selection.
type ScoredItem[T any] struct {
	Item  T
	Score float64
}

// minHeap implements a min-heap for ScoredItem.
// We use a min-heap to efficiently maintain top-K highest scores:
// the smallest score in the heap is always at the root, making it
// easy to decide if a new item should replace it.
type minHeap[T any] []ScoredItem[T]

func (h minHeap[T]) Len() int           { return len(h) }
func (h minHeap[T]) Less(i, j int) bool { return h[i].Score < h[j].Score } // min-heap
func (h minHeap[T]) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *minHeap[T]) Push(x any) {
	*h = append(*h, x.(ScoredItem[T]))
}

func (h *minHeap[T]) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// TopKByScore returns the top K items with the highest scores using a heap.
// This is O(n log k) which is more efficient than sorting O(n log n) when k << n.
// The returned slice is sorted in descending order by score.
func TopKByScore[T any](items []ScoredItem[T], k int) []ScoredItem[T] {
	if k <= 0 || len(items) == 0 {
		return nil
	}

	if k >= len(items) {
		// If k >= n, just sort and return all
		result := make([]ScoredItem[T], len(items))
		copy(result, items)
		// Sort descending by score
		for i := 0; i < len(result)-1; i++ {
			for j := i + 1; j < len(result); j++ {
				if result[j].Score > result[i].Score {
					result[i], result[j] = result[j], result[i]
				}
			}
		}
		return result
	}

	// Use a min-heap of size k to track the top k items
	h := make(minHeap[T], 0, k)
	heap.Init(&h)

	for _, item := range items {
		if h.Len() < k {
			heap.Push(&h, item)
		} else if item.Score > h[0].Score {
			// Replace the smallest item in heap if current item has higher score
			heap.Pop(&h)
			heap.Push(&h, item)
		}
	}

	// Extract items from heap and reverse to get descending order
	result := make([]ScoredItem[T], h.Len())
	for i := len(result) - 1; i >= 0; i-- {
		result[i] = heap.Pop(&h).(ScoredItem[T])
	}

	return result
}

// TopKIndicesByScore returns the indices of the top K items with the highest scores.
// Useful when you need to reference back to the original slice.
// Returns indices in descending order by score.
func TopKIndicesByScore(scores []float64, k int) []int {
	if k <= 0 || len(scores) == 0 {
		return nil
	}

	// Create scored items with indices
	items := make([]ScoredItem[int], len(scores))
	for i, score := range scores {
		items[i] = ScoredItem[int]{Item: i, Score: score}
	}

	topK := TopKByScore(items, k)
	indices := make([]int, len(topK))
	for i, item := range topK {
		indices[i] = item.Item
	}
	return indices
}

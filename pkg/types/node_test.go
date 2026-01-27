package types

import (
	"context"
	"errors"
	"testing"
	"time"
)

// mockNodeOperations implements NodeOperations for testing
type mockNodeOperations struct {
	results   interface{}
	execError error
	queries   []string
	params    []map[string]interface{}
}

func (m *mockNodeOperations) ExecuteQuery(query string, params map[string]interface{}) (interface{}, interface{}, interface{}, error) {
	m.queries = append(m.queries, query)
	m.params = append(m.params, params)
	return m.results, nil, nil, m.execError
}

func TestGetEpisodicNodeByUUID(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		mock := &mockNodeOperations{
			results: []map[string]interface{}{
				{
					"uuid":               "test-uuid-123",
					"name":               "test-episode",
					"content":            "This is test content",
					"group_id":           "group-1",
					"source_description": "Test source description",
					"entity_edges":       []interface{}{"edge-1", "edge-2"},
				},
			},
		}

		node, err := GetEpisodicNodeByUUID(context.Background(), mock, "test-uuid-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if node.Uuid != "test-uuid-123" {
			t.Errorf("expected uuid 'test-uuid-123', got '%s'", node.Uuid)
		}
		if node.Name != "test-episode" {
			t.Errorf("expected name 'test-episode', got '%s'", node.Name)
		}
		if node.Content != "This is test content" {
			t.Errorf("expected content 'This is test content', got '%s'", node.Content)
		}
		if node.GroupID != "group-1" {
			t.Errorf("expected group_id 'group-1', got '%s'", node.GroupID)
		}
		if node.Summary != "Test source description" {
			t.Errorf("expected summary 'Test source description', got '%s'", node.Summary)
		}
		if node.Type != EpisodicNodeType {
			t.Errorf("expected type '%s', got '%s'", EpisodicNodeType, node.Type)
		}
		if len(node.EntityEdges) != 2 {
			t.Errorf("expected 2 entity edges, got %d", len(node.EntityEdges))
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &mockNodeOperations{
			results: []map[string]interface{}{},
		}

		_, err := GetEpisodicNodeByUUID(context.Background(), mock, "non-existent")
		if err == nil {
			t.Fatal("expected error for non-existent UUID")
		}
		if err.Error() != "episode with UUID non-existent not found" {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("query error", func(t *testing.T) {
		mock := &mockNodeOperations{
			execError: errors.New("database connection failed"),
		}

		_, err := GetEpisodicNodeByUUID(context.Background(), mock, "test-uuid")
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != "database connection failed" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("invalid result type", func(t *testing.T) {
		mock := &mockNodeOperations{
			results: "invalid type",
		}

		_, err := GetEpisodicNodeByUUID(context.Background(), mock, "test-uuid")
		if err == nil {
			t.Fatal("expected error for invalid result type")
		}
	})
}

func TestDeleteNode(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		mock := &mockNodeOperations{}
		node := &Node{Uuid: "test-uuid-123"}

		err := DeleteNode(context.Background(), mock, node)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(mock.queries) != 1 {
			t.Fatalf("expected 1 query, got %d", len(mock.queries))
		}
		if mock.params[0]["uuid"] != "test-uuid-123" {
			t.Errorf("expected uuid parameter 'test-uuid-123', got '%v'", mock.params[0]["uuid"])
		}
	})

	t.Run("query error", func(t *testing.T) {
		mock := &mockNodeOperations{
			execError: errors.New("delete failed"),
		}
		node := &Node{Uuid: "test-uuid"}

		err := DeleteNode(context.Background(), mock, node)
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != "delete failed" {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestDeleteNodesByUUIDs(t *testing.T) {
	t.Parallel()
	t.Run("empty list", func(t *testing.T) {
		mock := &mockNodeOperations{}

		err := DeleteNodesByUUIDs(context.Background(), mock, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(mock.queries) != 0 {
			t.Errorf("expected 0 queries for empty list, got %d", len(mock.queries))
		}
	})

	t.Run("success with multiple UUIDs", func(t *testing.T) {
		mock := &mockNodeOperations{}
		uuids := []string{"uuid-1", "uuid-2", "uuid-3"}

		err := DeleteNodesByUUIDs(context.Background(), mock, uuids)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should run 3 queries (one for each label: Entity, Episodic, Community)
		if len(mock.queries) != 3 {
			t.Errorf("expected 3 queries (one per label), got %d", len(mock.queries))
		}

		// Check that UUIDs are passed correctly
		for _, params := range mock.params {
			paramUUIDs, ok := params["uuids"].([]string)
			if !ok {
				t.Fatal("expected uuids parameter to be []string")
			}
			if len(paramUUIDs) != 3 {
				t.Errorf("expected 3 uuids, got %d", len(paramUUIDs))
			}
		}
	})

	t.Run("query error", func(t *testing.T) {
		mock := &mockNodeOperations{
			execError: errors.New("batch delete failed"),
		}

		err := DeleteNodesByUUIDs(context.Background(), mock, []string{"uuid-1"})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestGetMentionedNodes(t *testing.T) {
	t.Parallel()
	t.Run("empty episodes", func(t *testing.T) {
		mock := &mockNodeOperations{}

		nodes, err := GetMentionedNodes(context.Background(), mock, []*Node{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(nodes) != 0 {
			t.Errorf("expected empty result, got %d nodes", len(nodes))
		}
	})

	t.Run("success with results", func(t *testing.T) {
		mock := &mockNodeOperations{
			results: []map[string]interface{}{
				{
					"uuid":        "entity-1",
					"name":        "Entity One",
					"entity_type": "Person",
					"summary":     "A test entity",
					"group_id":    "group-1",
				},
				{
					"uuid":        "entity-2",
					"name":        "Entity Two",
					"entity_type": "Organization",
					"summary":     "Another entity",
					"group_id":    "group-1",
				},
			},
		}

		episodes := []*Node{
			{Uuid: "episode-1"},
			{Uuid: "episode-2"},
		}

		nodes, err := GetMentionedNodes(context.Background(), mock, episodes)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(nodes) != 2 {
			t.Fatalf("expected 2 nodes, got %d", len(nodes))
		}

		if nodes[0].Uuid != "entity-1" {
			t.Errorf("expected first node uuid 'entity-1', got '%s'", nodes[0].Uuid)
		}
		if nodes[0].Type != EntityNodeType {
			t.Errorf("expected node type '%s', got '%s'", EntityNodeType, nodes[0].Type)
		}
		if nodes[1].EntityType != "Organization" {
			t.Errorf("expected entity type 'Organization', got '%s'", nodes[1].EntityType)
		}
	})

	t.Run("query error", func(t *testing.T) {
		mock := &mockNodeOperations{
			execError: errors.New("query failed"),
		}

		episodes := []*Node{{Uuid: "episode-1"}}

		_, err := GetMentionedNodes(context.Background(), mock, episodes)
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestParseNodeFromMap(t *testing.T) {
	t.Parallel()
	t.Run("parse with uuid", func(t *testing.T) {
		data := map[string]interface{}{
			"uuid":     "test-uuid",
			"name":     "Test Node",
			"group_id": "group-1",
			"content":  "Test content",
			"summary":  "Test summary",
		}

		node, err := ParseNodeFromMap(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if node.Uuid != "test-uuid" {
			t.Errorf("expected uuid 'test-uuid', got '%s'", node.Uuid)
		}
		if node.Name != "Test Node" {
			t.Errorf("expected name 'Test Node', got '%s'", node.Name)
		}
		if node.GroupID != "group-1" {
			t.Errorf("expected group_id 'group-1', got '%s'", node.GroupID)
		}
		if node.Content != "Test content" {
			t.Errorf("expected content 'Test content', got '%s'", node.Content)
		}
	})

	t.Run("parse with id fallback", func(t *testing.T) {
		data := map[string]interface{}{
			"id":   "test-id",
			"name": "Test Node",
		}

		node, err := ParseNodeFromMap(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if node.Uuid != "test-id" {
			t.Errorf("expected uuid 'test-id', got '%s'", node.Uuid)
		}
	})

	t.Run("parse with timestamps", func(t *testing.T) {
		now := time.Now()
		data := map[string]interface{}{
			"uuid":       "test-uuid",
			"valid_at":   now,
			"created_at": now.Add(-time.Hour),
			"updated_at": now.Add(-time.Minute),
		}

		node, err := ParseNodeFromMap(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !node.ValidFrom.Equal(now) {
			t.Errorf("expected valid_from %v, got %v", now, node.ValidFrom)
		}
		if !node.CreatedAt.Equal(now.Add(-time.Hour)) {
			t.Errorf("expected created_at %v, got %v", now.Add(-time.Hour), node.CreatedAt)
		}
	})

	t.Run("parse with valid_from instead of valid_at", func(t *testing.T) {
		now := time.Now()
		data := map[string]interface{}{
			"uuid":       "test-uuid",
			"valid_from": now,
		}

		node, err := ParseNodeFromMap(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !node.ValidFrom.Equal(now) {
			t.Errorf("expected valid_from %v, got %v", now, node.ValidFrom)
		}
	})

	t.Run("parse with episode_type", func(t *testing.T) {
		data := map[string]interface{}{
			"uuid":         "test-uuid",
			"episode_type": "message",
		}

		node, err := ParseNodeFromMap(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if node.EpisodeType != EpisodeType("message") {
			t.Errorf("expected episode_type 'message', got '%s'", node.EpisodeType)
		}
	})

	t.Run("empty map", func(t *testing.T) {
		node, err := ParseNodeFromMap(map[string]interface{}{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if node.Metadata == nil {
			t.Error("expected Metadata to be initialized")
		}
		if node.Type != EpisodicNodeType {
			t.Errorf("expected default type '%s', got '%s'", EpisodicNodeType, node.Type)
		}
	})
}

func TestReverseNodes(t *testing.T) {
	t.Parallel()
	t.Run("empty slice", func(t *testing.T) {
		nodes := []*Node{}
		ReverseNodes(nodes)
		if len(nodes) != 0 {
			t.Error("expected empty slice to remain empty")
		}
	})

	t.Run("single element", func(t *testing.T) {
		nodes := []*Node{{Uuid: "1"}}
		ReverseNodes(nodes)
		if nodes[0].Uuid != "1" {
			t.Errorf("expected uuid '1', got '%s'", nodes[0].Uuid)
		}
	})

	t.Run("even number of elements", func(t *testing.T) {
		nodes := []*Node{
			{Uuid: "1"},
			{Uuid: "2"},
			{Uuid: "3"},
			{Uuid: "4"},
		}
		ReverseNodes(nodes)

		expected := []string{"4", "3", "2", "1"}
		for i, node := range nodes {
			if node.Uuid != expected[i] {
				t.Errorf("expected uuid '%s' at position %d, got '%s'", expected[i], i, node.Uuid)
			}
		}
	})

	t.Run("odd number of elements", func(t *testing.T) {
		nodes := []*Node{
			{Uuid: "a"},
			{Uuid: "b"},
			{Uuid: "c"},
		}
		ReverseNodes(nodes)

		expected := []string{"c", "b", "a"}
		for i, node := range nodes {
			if node.Uuid != expected[i] {
				t.Errorf("expected uuid '%s' at position %d, got '%s'", expected[i], i, node.Uuid)
			}
		}
	})
}

func TestConvertToFloat32Array(t *testing.T) {
	t.Parallel()
	t.Run("nil input", func(t *testing.T) {
		result := convertToFloat32Array(nil)
		if result != nil {
			t.Errorf("expected nil for nil input, got %v", result)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		result := convertToFloat32Array([]interface{}{})
		if len(result) != 0 {
			t.Errorf("expected empty slice, got %v", result)
		}
	})

	t.Run("valid float64 values", func(t *testing.T) {
		input := []interface{}{1.0, 2.5, 3.7, -0.5}
		result := convertToFloat32Array(input)

		if len(result) != 4 {
			t.Fatalf("expected 4 elements, got %d", len(result))
		}

		expected := []float32{1.0, 2.5, 3.7, -0.5}
		for i, val := range result {
			if val != expected[i] {
				t.Errorf("expected %f at position %d, got %f", expected[i], i, val)
			}
		}
	})

	t.Run("mixed types with non-float64", func(t *testing.T) {
		// Non-float64 values should result in zero values
		input := []interface{}{1.0, "string", 3.0}
		result := convertToFloat32Array(input)

		if len(result) != 3 {
			t.Fatalf("expected 3 elements, got %d", len(result))
		}

		if result[0] != 1.0 {
			t.Errorf("expected 1.0 at position 0, got %f", result[0])
		}
		if result[1] != 0.0 {
			t.Errorf("expected 0.0 for non-float64 at position 1, got %f", result[1])
		}
		if result[2] != 3.0 {
			t.Errorf("expected 3.0 at position 2, got %f", result[2])
		}
	})

	t.Run("invalid type returns nil", func(t *testing.T) {
		result := convertToFloat32Array("not a slice")
		if result != nil {
			t.Errorf("expected nil for invalid type, got %v", result)
		}
	})

	t.Run("large embedding vector", func(t *testing.T) {
		// Simulate typical embedding size (e.g., 384 dimensions)
		input := make([]interface{}, 384)
		for i := range input {
			input[i] = float64(i) * 0.001
		}

		result := convertToFloat32Array(input)
		if len(result) != 384 {
			t.Fatalf("expected 384 elements, got %d", len(result))
		}

		// Spot check a few values
		if result[0] != 0.0 {
			t.Errorf("expected 0.0 at position 0, got %f", result[0])
		}
		if result[100] != float32(0.1) {
			t.Errorf("expected 0.1 at position 100, got %f", result[100])
		}
	})
}

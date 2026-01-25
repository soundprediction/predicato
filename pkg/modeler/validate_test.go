package modeler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/soundprediction/predicato/pkg/types"
)

func TestValidateModeler(t *testing.T) {
	ctx := context.Background()

	t.Run("nil modeler returns error", func(t *testing.T) {
		_, err := ValidateModeler(ctx, nil, nil)
		if err == nil {
			t.Error("expected error for nil modeler")
		}
	})

	t.Run("NoOpModeler passes validation", func(t *testing.T) {
		modeler := &NoOpModeler{}
		result, err := ValidateModeler(ctx, modeler, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Valid {
			t.Error("expected NoOpModeler to be valid")
		}
		if !result.EntityResolution.Passed {
			t.Error("expected EntityResolution to pass")
		}
		if !result.RelationshipResolution.Passed {
			t.Error("expected RelationshipResolution to pass")
		}
	})

	t.Run("custom options", func(t *testing.T) {
		modeler := &NoOpModeler{}
		opts := &ValidateModelerOptions{
			NodeCount:               3,
			EdgeCount:               2,
			GroupID:                 "custom-group",
			Timeout:                 5 * time.Second,
			SkipCommunityValidation: true,
		}
		result, err := ValidateModeler(ctx, modeler, opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Valid {
			t.Error("expected valid result")
		}
		if result.CommunityDetection != nil {
			t.Error("expected nil CommunityDetection when skipped")
		}
	})

	t.Run("modeler that returns error", func(t *testing.T) {
		modeler := &MockGraphModeler{
			resolveEntitiesFn: func(ctx context.Context, input *EntityResolutionInput) (*EntityResolutionOutput, error) {
				return nil, errors.New("entity resolution failed")
			},
		}
		result, err := ValidateModeler(ctx, modeler, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Valid {
			t.Error("expected invalid result when entity resolution fails")
		}
		if result.EntityResolution.Passed {
			t.Error("expected EntityResolution to fail")
		}
		if result.EntityResolution.Error == nil {
			t.Error("expected error to be recorded")
		}
	})

	t.Run("modeler that returns nil output", func(t *testing.T) {
		modeler := &MockGraphModeler{
			resolveEntitiesFn: func(ctx context.Context, input *EntityResolutionInput) (*EntityResolutionOutput, error) {
				return nil, nil // nil output is invalid
			},
		}
		result, err := ValidateModeler(ctx, modeler, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Valid {
			t.Error("expected invalid result when output is nil")
		}
	})

	t.Run("modeler with nil UUIDMap", func(t *testing.T) {
		modeler := &MockGraphModeler{
			resolveEntitiesFn: func(ctx context.Context, input *EntityResolutionInput) (*EntityResolutionOutput, error) {
				return &EntityResolutionOutput{
					ResolvedNodes: []*types.Node{},
					UUIDMap:       nil, // nil UUIDMap is invalid
				}, nil
			},
		}
		result, err := ValidateModeler(ctx, modeler, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Valid {
			t.Error("expected invalid result when UUIDMap is nil")
		}
	})

	t.Run("community detection failure is warning", func(t *testing.T) {
		modeler := &MockGraphModeler{
			buildCommunitiesFn: func(ctx context.Context, input *CommunityInput) (*CommunityOutput, error) {
				return nil, errors.New("community detection failed")
			},
		}
		result, err := ValidateModeler(ctx, modeler, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Community failure should result in warning but still valid
		if !result.Valid {
			t.Error("expected valid result (community failure is warning)")
		}
		if len(result.Warnings) == 0 {
			t.Error("expected warning for community detection failure")
		}
	})
}

func TestGenerateSampleData(t *testing.T) {
	t.Run("generateSampleEpisode", func(t *testing.T) {
		episode := generateSampleEpisode("test-group")
		if episode == nil {
			t.Fatal("expected non-nil episode")
		}
		if episode.Uuid == "" {
			t.Error("expected non-empty UUID")
		}
		if episode.GroupID != "test-group" {
			t.Errorf("expected GroupID='test-group', got '%s'", episode.GroupID)
		}
		if episode.Type != types.EpisodicNodeType {
			t.Errorf("expected Type=EpisodicNodeType, got '%s'", episode.Type)
		}
	})

	t.Run("generateSampleNodes", func(t *testing.T) {
		nodes := generateSampleNodes(5, "test-group")
		if len(nodes) != 5 {
			t.Errorf("expected 5 nodes, got %d", len(nodes))
		}
		for _, n := range nodes {
			if n.Uuid == "" {
				t.Error("expected non-empty UUID")
			}
			if n.Name == "" {
				t.Error("expected non-empty Name")
			}
			if n.GroupID != "test-group" {
				t.Errorf("expected GroupID='test-group', got '%s'", n.GroupID)
			}
			if n.Type != types.EntityNodeType {
				t.Errorf("expected Type=EntityNodeType, got '%s'", n.Type)
			}
		}
	})

	t.Run("generateSampleNodes respects count", func(t *testing.T) {
		nodes := generateSampleNodes(15, "test") // More than available names
		if len(nodes) != 10 {
			t.Errorf("expected 10 nodes (max available), got %d", len(nodes))
		}
	})

	t.Run("generateSampleEdges", func(t *testing.T) {
		nodes := generateSampleNodes(5, "test-group")
		edges := generateSampleEdges(3, nodes, "test-group")
		if len(edges) != 3 {
			t.Errorf("expected 3 edges, got %d", len(edges))
		}
		for _, e := range edges {
			if e.Uuid == "" {
				t.Error("expected non-empty UUID")
			}
			if e.Name == "" {
				t.Error("expected non-empty Name")
			}
			if e.SourceNodeID == "" || e.TargetNodeID == "" {
				t.Error("expected non-empty source/target node IDs")
			}
		}
	})

	t.Run("generateSampleEdges with insufficient nodes", func(t *testing.T) {
		edges := generateSampleEdges(3, []*types.Node{}, "test")
		if len(edges) != 0 {
			t.Errorf("expected 0 edges for empty nodes, got %d", len(edges))
		}

		edges = generateSampleEdges(3, []*types.Node{{Uuid: "single"}}, "test")
		if len(edges) != 0 {
			t.Errorf("expected 0 edges for single node, got %d", len(edges))
		}
	})
}

func TestNoOpModeler(t *testing.T) {
	ctx := context.Background()
	modeler := &NoOpModeler{}

	t.Run("ResolveEntities identity mapping", func(t *testing.T) {
		nodes := []*types.Node{
			{Uuid: "uuid-1", Name: "Node1"},
			{Uuid: "uuid-2", Name: "Node2"},
		}
		input := &EntityResolutionInput{
			ExtractedNodes: nodes,
			Episode:        &types.Node{Uuid: "ep-1"},
		}

		output, err := modeler.ResolveEntities(ctx, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should return same nodes
		if len(output.ResolvedNodes) != 2 {
			t.Errorf("expected 2 nodes, got %d", len(output.ResolvedNodes))
		}

		// UUID map should be identity
		for _, n := range nodes {
			if mapped, ok := output.UUIDMap[n.Uuid]; !ok || mapped != n.Uuid {
				t.Errorf("expected identity mapping for %s", n.Uuid)
			}
		}

		// All should be "new"
		if output.NewCount != 2 {
			t.Errorf("expected NewCount=2, got %d", output.NewCount)
		}
		if output.MergedCount != 0 {
			t.Errorf("expected MergedCount=0, got %d", output.MergedCount)
		}
	})

	t.Run("ResolveEntities with nil input", func(t *testing.T) {
		output, err := modeler.ResolveEntities(ctx, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(output.ResolvedNodes) != 0 {
			t.Error("expected empty ResolvedNodes")
		}
		if len(output.UUIDMap) != 0 {
			t.Error("expected empty UUIDMap")
		}
	})

	t.Run("ResolveRelationships creates episodic edges", func(t *testing.T) {
		nodes := []*types.Node{
			{Uuid: "node-1"},
			{Uuid: "node-2"},
		}
		edges := []*types.Edge{
			types.NewEntityEdge("edge-1", "node-1", "node-2", "test", "KNOWS", types.EntityEdgeType),
		}
		input := &RelationshipResolutionInput{
			ExtractedEdges: edges,
			ResolvedNodes:  nodes,
			UUIDMap:        map[string]string{"node-1": "node-1", "node-2": "node-2"},
			Episode:        &types.Node{Uuid: "ep-1"},
		}

		output, err := modeler.ResolveRelationships(ctx, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should return same edges
		if len(output.ResolvedEdges) != 1 {
			t.Errorf("expected 1 edge, got %d", len(output.ResolvedEdges))
		}

		// Should create episodic edges for each node
		if len(output.EpisodicEdges) != 2 {
			t.Errorf("expected 2 episodic edges, got %d", len(output.EpisodicEdges))
		}
	})

	t.Run("ResolveRelationships with nil input", func(t *testing.T) {
		output, err := modeler.ResolveRelationships(ctx, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(output.ResolvedEdges) != 0 {
			t.Error("expected empty ResolvedEdges")
		}
		if len(output.EpisodicEdges) != 0 {
			t.Error("expected empty EpisodicEdges")
		}
	})

	t.Run("BuildCommunities always returns nil", func(t *testing.T) {
		input := &CommunityInput{
			Nodes:   []*types.Node{{Uuid: "node-1"}},
			GroupID: "test",
		}

		output, err := modeler.BuildCommunities(ctx, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output != nil {
			t.Error("expected nil output from NoOpModeler")
		}
	})
}

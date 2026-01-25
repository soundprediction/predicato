package modeler

import (
	"context"
	"testing"

	"github.com/soundprediction/predicato/pkg/types"
)

// MockGraphModeler is a test implementation of GraphModeler
type MockGraphModeler struct {
	resolveEntitiesFn      func(ctx context.Context, input *EntityResolutionInput) (*EntityResolutionOutput, error)
	resolveRelationshipsFn func(ctx context.Context, input *RelationshipResolutionInput) (*RelationshipResolutionOutput, error)
	buildCommunitiesFn     func(ctx context.Context, input *CommunityInput) (*CommunityOutput, error)
}

func (m *MockGraphModeler) ResolveEntities(ctx context.Context, input *EntityResolutionInput) (*EntityResolutionOutput, error) {
	if m.resolveEntitiesFn != nil {
		return m.resolveEntitiesFn(ctx, input)
	}
	// Default: return identity mapping
	uuidMap := make(map[string]string)
	for _, n := range input.ExtractedNodes {
		uuidMap[n.Uuid] = n.Uuid
	}
	return &EntityResolutionOutput{
		ResolvedNodes: input.ExtractedNodes,
		UUIDMap:       uuidMap,
		NewCount:      len(input.ExtractedNodes),
	}, nil
}

func (m *MockGraphModeler) ResolveRelationships(ctx context.Context, input *RelationshipResolutionInput) (*RelationshipResolutionOutput, error) {
	if m.resolveRelationshipsFn != nil {
		return m.resolveRelationshipsFn(ctx, input)
	}
	return &RelationshipResolutionOutput{
		ResolvedEdges: input.ExtractedEdges,
		EpisodicEdges: []*types.Edge{},
		NewCount:      len(input.ExtractedEdges),
	}, nil
}

func (m *MockGraphModeler) BuildCommunities(ctx context.Context, input *CommunityInput) (*CommunityOutput, error) {
	if m.buildCommunitiesFn != nil {
		return m.buildCommunitiesFn(ctx, input)
	}
	return nil, nil
}

// Ensure MockGraphModeler implements GraphModeler
var _ GraphModeler = (*MockGraphModeler)(nil)

func TestGraphModelerInterface(t *testing.T) {
	ctx := context.Background()
	mock := &MockGraphModeler{}

	t.Run("ResolveEntities with nil input", func(t *testing.T) {
		mock.resolveEntitiesFn = func(ctx context.Context, input *EntityResolutionInput) (*EntityResolutionOutput, error) {
			if input == nil {
				return &EntityResolutionOutput{
					ResolvedNodes: []*types.Node{},
					UUIDMap:       make(map[string]string),
				}, nil
			}
			return nil, nil
		}

		output, err := mock.ResolveEntities(ctx, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output == nil {
			t.Fatal("expected non-nil output")
		}
		if output.UUIDMap == nil {
			t.Fatal("expected non-nil UUIDMap")
		}
	})

	t.Run("ResolveEntities with nodes", func(t *testing.T) {
		mock.resolveEntitiesFn = nil // Use default

		nodes := []*types.Node{
			{Uuid: "node-1", Name: "Alice"},
			{Uuid: "node-2", Name: "Bob"},
		}
		input := &EntityResolutionInput{
			ExtractedNodes: nodes,
			Episode:        &types.Node{Uuid: "ep-1"},
			GroupID:        "test",
		}

		output, err := mock.ResolveEntities(ctx, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(output.ResolvedNodes) != 2 {
			t.Errorf("expected 2 resolved nodes, got %d", len(output.ResolvedNodes))
		}
		if output.NewCount != 2 {
			t.Errorf("expected NewCount=2, got %d", output.NewCount)
		}
		if len(output.UUIDMap) != 2 {
			t.Errorf("expected UUIDMap with 2 entries, got %d", len(output.UUIDMap))
		}
	})

	t.Run("ResolveRelationships with edges", func(t *testing.T) {
		edges := []*types.Edge{
			types.NewEntityEdge("edge-1", "node-1", "node-2", "test", "KNOWS", types.EntityEdgeType),
		}
		input := &RelationshipResolutionInput{
			ExtractedEdges: edges,
			ResolvedNodes:  []*types.Node{{Uuid: "node-1"}, {Uuid: "node-2"}},
			UUIDMap:        map[string]string{"node-1": "node-1", "node-2": "node-2"},
			Episode:        &types.Node{Uuid: "ep-1"},
			GroupID:        "test",
		}

		output, err := mock.ResolveRelationships(ctx, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(output.ResolvedEdges) != 1 {
			t.Errorf("expected 1 resolved edge, got %d", len(output.ResolvedEdges))
		}
		if output.NewCount != 1 {
			t.Errorf("expected NewCount=1, got %d", output.NewCount)
		}
	})

	t.Run("BuildCommunities returns nil", func(t *testing.T) {
		input := &CommunityInput{
			Nodes:   []*types.Node{{Uuid: "node-1"}},
			GroupID: "test",
		}

		output, err := mock.BuildCommunities(ctx, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output != nil {
			t.Error("expected nil output for skipped community detection")
		}
	})
}

func TestEntityResolutionInput(t *testing.T) {
	t.Run("empty options", func(t *testing.T) {
		input := &EntityResolutionInput{
			ExtractedNodes: []*types.Node{},
			GroupID:        "test",
		}
		if input.Options != nil {
			t.Error("expected nil options")
		}
	})

	t.Run("with options", func(t *testing.T) {
		input := &EntityResolutionInput{
			Options: &EntityResolutionOptions{
				SkipResolution:      true,
				SkipReflexion:       true,
				SimilarityThreshold: 0.9,
			},
		}
		if !input.Options.SkipResolution {
			t.Error("expected SkipResolution=true")
		}
		if input.Options.SimilarityThreshold != 0.9 {
			t.Errorf("expected SimilarityThreshold=0.9, got %f", input.Options.SimilarityThreshold)
		}
	})
}

func TestRelationshipResolutionInput(t *testing.T) {
	t.Run("with EdgeTypes", func(t *testing.T) {
		input := &RelationshipResolutionInput{
			ExtractedEdges: []*types.Edge{},
			EdgeTypes: map[string]interface{}{
				"KNOWS":    true,
				"WORKS_AT": true,
			},
		}
		if len(input.EdgeTypes) != 2 {
			t.Errorf("expected 2 edge types, got %d", len(input.EdgeTypes))
		}
	})
}

func TestModelerValidationResult(t *testing.T) {
	t.Run("valid result", func(t *testing.T) {
		result := &ModelerValidationResult{
			Valid: true,
			EntityResolution: &ValidationStepResult{
				Passed:  true,
				Latency: 100,
			},
			RelationshipResolution: &ValidationStepResult{
				Passed:  true,
				Latency: 50,
			},
			Warnings: []string{},
		}

		if !result.Valid {
			t.Error("expected Valid=true")
		}
		if !result.EntityResolution.Passed {
			t.Error("expected EntityResolution.Passed=true")
		}
	})

	t.Run("invalid result", func(t *testing.T) {
		result := &ModelerValidationResult{
			Valid: false,
			EntityResolution: &ValidationStepResult{
				Passed: false,
				Error:  context.DeadlineExceeded,
			},
			Warnings: []string{"EntityResolution timed out"},
		}

		if result.Valid {
			t.Error("expected Valid=false")
		}
		if result.EntityResolution.Error == nil {
			t.Error("expected error to be set")
		}
		if len(result.Warnings) != 1 {
			t.Errorf("expected 1 warning, got %d", len(result.Warnings))
		}
	})
}

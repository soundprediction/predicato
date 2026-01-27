package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestEntityEdgeSyncFields(t *testing.T) {
	now := time.Now()
	validAt := now.Add(-time.Hour)
	invalidAt := now.Add(time.Hour)

	edge := &EntityEdge{
		BaseEdge: BaseEdge{
			Uuid:         "test-uuid",
			GroupID:      "group-1",
			SourceNodeID: "source-uuid",
			TargetNodeID: "target-uuid",
			CreatedAt:    now,
		},
		Name:      "test-relation",
		Fact:      "test fact",
		ValidAt:   &validAt,
		InvalidAt: &invalidAt,
	}

	edge.syncFields()

	if edge.SourceID != edge.SourceNodeID {
		t.Errorf("SourceID = %s, want %s", edge.SourceID, edge.SourceNodeID)
	}
	if edge.TargetID != edge.TargetNodeID {
		t.Errorf("TargetID = %s, want %s", edge.TargetID, edge.TargetNodeID)
	}
	if edge.Summary != edge.Fact {
		t.Errorf("Summary = %s, want %s", edge.Summary, edge.Fact)
	}
	if !edge.ValidFrom.Equal(validAt) {
		t.Errorf("ValidFrom = %v, want %v", edge.ValidFrom, validAt)
	}
	if edge.ValidTo == nil || !edge.ValidTo.Equal(invalidAt) {
		t.Errorf("ValidTo = %v, want %v", edge.ValidTo, invalidAt)
	}
	if edge.Type != EntityEdgeType {
		t.Errorf("Type = %s, want %s", edge.Type, EntityEdgeType)
	}
}

func TestEntityEdgeUpdateFromCompat(t *testing.T) {
	now := time.Now()
	validFrom := now.Add(-time.Hour)
	validTo := now.Add(time.Hour)

	edge := &EntityEdge{
		SourceID:  "compat-source",
		TargetID:  "compat-target",
		Summary:   "compat summary",
		ValidFrom: validFrom,
		ValidTo:   &validTo,
	}

	edge.updateFromCompat()

	if edge.SourceNodeID != "compat-source" {
		t.Errorf("SourceNodeID = %s, want compat-source", edge.SourceNodeID)
	}
	if edge.TargetNodeID != "compat-target" {
		t.Errorf("TargetNodeID = %s, want compat-target", edge.TargetNodeID)
	}
	if edge.Fact != "compat summary" {
		t.Errorf("Fact = %s, want compat summary", edge.Fact)
	}
	if edge.ValidAt == nil || !edge.ValidAt.Equal(validFrom) {
		t.Errorf("ValidAt = %v, want %v", edge.ValidAt, validFrom)
	}
	if edge.InvalidAt == nil || !edge.InvalidAt.Equal(validTo) {
		t.Errorf("InvalidAt = %v, want %v", edge.InvalidAt, validTo)
	}
}

func TestEntityEdgeJSONRoundtrip(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	validAt := now.Add(-time.Hour)

	original := &EntityEdge{
		BaseEdge: BaseEdge{
			Uuid:         "test-uuid",
			GroupID:      "group-1",
			SourceNodeID: "source-uuid",
			TargetNodeID: "target-uuid",
			CreatedAt:    now,
		},
		Name:          "test-relation",
		Fact:          "test fact",
		ValidAt:       &validAt,
		Episodes:      []string{"ep-1", "ep-2"},
		FactEmbedding: []float32{0.1, 0.2, 0.3},
		Strength:      0.85,
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Unmarshal back
	var decoded EntityEdge
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Check key fields
	if decoded.Uuid != original.Uuid {
		t.Errorf("Uuid mismatch: got %s, want %s", decoded.Uuid, original.Uuid)
	}
	if decoded.Name != original.Name {
		t.Errorf("Name mismatch: got %s, want %s", decoded.Name, original.Name)
	}
	if decoded.Fact != original.Fact {
		t.Errorf("Fact mismatch: got %s, want %s", decoded.Fact, original.Fact)
	}
	if decoded.SourceNodeID != original.SourceNodeID {
		t.Errorf("SourceNodeID mismatch: got %s, want %s", decoded.SourceNodeID, original.SourceNodeID)
	}
	if decoded.TargetNodeID != original.TargetNodeID {
		t.Errorf("TargetNodeID mismatch: got %s, want %s", decoded.TargetNodeID, original.TargetNodeID)
	}
	if len(decoded.Episodes) != len(original.Episodes) {
		t.Errorf("Episodes length mismatch: got %d, want %d", len(decoded.Episodes), len(original.Episodes))
	}
	if len(decoded.FactEmbedding) != len(original.FactEmbedding) {
		t.Errorf("FactEmbedding length mismatch: got %d, want %d", len(decoded.FactEmbedding), len(original.FactEmbedding))
	}
}

func TestBaseEdgeGetters(t *testing.T) {
	now := time.Now()
	edge := &BaseEdge{
		Uuid:         "test-uuid",
		GroupID:      "group-1",
		SourceNodeID: "source-uuid",
		TargetNodeID: "target-uuid",
		CreatedAt:    now,
	}

	if edge.GetUUID() != "test-uuid" {
		t.Errorf("GetUUID() = %s, want test-uuid", edge.GetUUID())
	}
	if edge.GetGroupID() != "group-1" {
		t.Errorf("GetGroupID() = %s, want group-1", edge.GetGroupID())
	}
	if edge.GetSourceNodeUUID() != "source-uuid" {
		t.Errorf("GetSourceNodeUUID() = %s, want source-uuid", edge.GetSourceNodeUUID())
	}
	if edge.GetTargetNodeUUID() != "target-uuid" {
		t.Errorf("GetTargetNodeUUID() = %s, want target-uuid", edge.GetTargetNodeUUID())
	}
	if !edge.GetCreatedAt().Equal(now) {
		t.Errorf("GetCreatedAt() = %v, want %v", edge.GetCreatedAt(), now)
	}
}

func TestEdgeTypes(t *testing.T) {
	// Verify constant values haven't changed
	if EntityEdgeType != "entity" {
		t.Errorf("EntityEdgeType = %s, want entity", EntityEdgeType)
	}
	if EpisodicEdgeType != "episodic" {
		t.Errorf("EpisodicEdgeType = %s, want episodic", EpisodicEdgeType)
	}
	if CommunityEdgeType != "community" {
		t.Errorf("CommunityEdgeType = %s, want community", CommunityEdgeType)
	}
	if SourceEdgeType != "source" {
		t.Errorf("SourceEdgeType = %s, want source", SourceEdgeType)
	}
}

func TestGraphProviders(t *testing.T) {
	// Verify constant values haven't changed
	if GraphProviderNeo4j != "neo4j" {
		t.Errorf("GraphProviderNeo4j = %s, want neo4j", GraphProviderNeo4j)
	}
	if GraphProviderFalkorDB != "falkordb" {
		t.Errorf("GraphProviderFalkorDB = %s, want falkordb", GraphProviderFalkorDB)
	}
	if GraphProviderLadybug != "ladybug" {
		t.Errorf("GraphProviderLadybug = %s, want ladybug", GraphProviderLadybug)
	}
	if GraphProviderNeptune != "neptune" {
		t.Errorf("GraphProviderNeptune = %s, want neptune", GraphProviderNeptune)
	}
}

func TestCommunityEdge(t *testing.T) {
	now := time.Now()
	edge := &CommunityEdge{
		BaseEdge: BaseEdge{
			Uuid:         "community-edge-uuid",
			GroupID:      "group-1",
			SourceNodeID: "community-uuid",
			TargetNodeID: "entity-uuid",
			CreatedAt:    now,
		},
	}

	if edge.GetUUID() != "community-edge-uuid" {
		t.Errorf("GetUUID() = %s, want community-edge-uuid", edge.GetUUID())
	}
	if edge.GetSourceNodeUUID() != "community-uuid" {
		t.Errorf("GetSourceNodeUUID() = %s, want community-uuid", edge.GetSourceNodeUUID())
	}
}

package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNodeValidation(t *testing.T) {
	tests := []struct {
		name    string
		node    Node
		wantErr error
	}{
		{
			name: "valid node",
			node: Node{
				Name:    "test-node",
				GroupID: "group-1",
			},
			wantErr: nil,
		},
		{
			name: "empty name",
			node: Node{
				Name:    "",
				GroupID: "group-1",
			},
			wantErr: ErrEmptyName,
		},
		{
			name: "empty group_id",
			node: Node{
				Name:    "test-node",
				GroupID: "",
			},
			wantErr: ErrEmptyGroupID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.Validate()
			if err != tt.wantErr {
				t.Errorf("Node.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNodeValidateForCreate(t *testing.T) {
	tests := []struct {
		name    string
		node    Node
		wantErr error
	}{
		{
			name: "valid node for create",
			node: Node{
				Uuid:    "uuid-123",
				Name:    "test-node",
				GroupID: "group-1",
			},
			wantErr: nil,
		},
		{
			name: "empty uuid",
			node: Node{
				Uuid:    "",
				Name:    "test-node",
				GroupID: "group-1",
			},
			wantErr: ErrEmptyUUID,
		},
		{
			name: "empty name",
			node: Node{
				Uuid:    "uuid-123",
				Name:    "",
				GroupID: "group-1",
			},
			wantErr: ErrEmptyName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.ValidateForCreate()
			if err != tt.wantErr {
				t.Errorf("Node.ValidateForCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEpisodeValidation(t *testing.T) {
	tests := []struct {
		name    string
		episode Episode
		wantErr error
	}{
		{
			name: "valid episode",
			episode: Episode{
				Content: "test content",
				GroupID: "group-1",
			},
			wantErr: nil,
		},
		{
			name: "empty content",
			episode: Episode{
				Content: "",
				GroupID: "group-1",
			},
			wantErr: ErrEmptyContent,
		},
		{
			name: "empty group_id",
			episode: Episode{
				Content: "test content",
				GroupID: "",
			},
			wantErr: ErrEmptyGroupID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.episode.Validate()
			if err != tt.wantErr {
				t.Errorf("Episode.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEpisodeValidateForCreate(t *testing.T) {
	tests := []struct {
		name    string
		episode Episode
		wantErr error
	}{
		{
			name: "valid episode for create",
			episode: Episode{
				ID:      "ep-123",
				Content: "test content",
				GroupID: "group-1",
			},
			wantErr: nil,
		},
		{
			name: "empty id",
			episode: Episode{
				ID:      "",
				Content: "test content",
				GroupID: "group-1",
			},
			wantErr: ErrEmptyID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.episode.ValidateForCreate()
			if err != tt.wantErr {
				t.Errorf("Episode.ValidateForCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSearchConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  SearchConfig
		wantErr error
	}{
		{
			name: "valid config",
			config: SearchConfig{
				Limit: 10,
			},
			wantErr: nil,
		},
		{
			name: "zero limit (valid)",
			config: SearchConfig{
				Limit: 0,
			},
			wantErr: nil,
		},
		{
			name: "negative limit",
			config: SearchConfig{
				Limit: -1,
			},
			wantErr: ErrInvalidLimit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if err != tt.wantErr {
				t.Errorf("SearchConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSearchConfigWithDefaults(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		var config *SearchConfig
		result := config.WithDefaults()
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.Limit != 10 {
			t.Errorf("expected default Limit=10, got %d", result.Limit)
		}
		if result.CenterNodeDistance != 2 {
			t.Errorf("expected default CenterNodeDistance=2, got %d", result.CenterNodeDistance)
		}
		if !result.IncludeEdges {
			t.Error("expected default IncludeEdges=true")
		}
	})

	t.Run("zero values get defaults", func(t *testing.T) {
		config := &SearchConfig{
			Limit:              0,
			CenterNodeDistance: 0,
			MinScore:           0.5, // non-zero value should be preserved
		}
		result := config.WithDefaults()
		if result.Limit != 10 {
			t.Errorf("expected default Limit=10, got %d", result.Limit)
		}
		if result.CenterNodeDistance != 2 {
			t.Errorf("expected default CenterNodeDistance=2, got %d", result.CenterNodeDistance)
		}
		if result.MinScore != 0.5 {
			t.Errorf("expected MinScore=0.5 to be preserved, got %f", result.MinScore)
		}
	})

	t.Run("non-zero values preserved", func(t *testing.T) {
		config := &SearchConfig{
			Limit:              20,
			CenterNodeDistance: 3,
		}
		result := config.WithDefaults()
		if result.Limit != 20 {
			t.Errorf("expected Limit=20 to be preserved, got %d", result.Limit)
		}
		if result.CenterNodeDistance != 3 {
			t.Errorf("expected CenterNodeDistance=3 to be preserved, got %d", result.CenterNodeDistance)
		}
	})
}

func TestNodeJSONRoundtrip(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	original := &Node{
		Uuid:       "test-uuid",
		Name:       "test-node",
		Type:       EntityNodeType,
		GroupID:    "group-1",
		CreatedAt:  now,
		UpdatedAt:  now,
		EntityType: "person",
		Summary:    "A test person",
		Embedding:  []float32{0.1, 0.2, 0.3},
		Metadata: map[string]interface{}{
			"key": "value",
		},
		ValidFrom: now,
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Unmarshal back
	var decoded Node
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
	if decoded.Type != original.Type {
		t.Errorf("Type mismatch: got %s, want %s", decoded.Type, original.Type)
	}
	if decoded.GroupID != original.GroupID {
		t.Errorf("GroupID mismatch: got %s, want %s", decoded.GroupID, original.GroupID)
	}
	if len(decoded.Embedding) != len(original.Embedding) {
		t.Errorf("Embedding length mismatch: got %d, want %d", len(decoded.Embedding), len(original.Embedding))
	}
}

func TestEpisodeJSONRoundtrip(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	original := &Episode{
		ID:        "test-episode",
		Name:      "Test Episode",
		Content:   "This is test content",
		Source:    "test-source",
		Reference: now,
		CreatedAt: now,
		GroupID:   "group-1",
		Metadata: map[string]interface{}{
			"key": "value",
		},
		ContentEmbedding: []float32{0.1, 0.2, 0.3},
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Unmarshal back
	var decoded Episode
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Check key fields
	if decoded.ID != original.ID {
		t.Errorf("ID mismatch: got %s, want %s", decoded.ID, original.ID)
	}
	if decoded.Content != original.Content {
		t.Errorf("Content mismatch: got %s, want %s", decoded.Content, original.Content)
	}
	if decoded.GroupID != original.GroupID {
		t.Errorf("GroupID mismatch: got %s, want %s", decoded.GroupID, original.GroupID)
	}
}

func TestNodeTypes(t *testing.T) {
	// Verify constant values haven't changed
	if EntityNodeType != "entity" {
		t.Errorf("EntityNodeType = %s, want entity", EntityNodeType)
	}
	if EpisodicNodeType != "episodic" {
		t.Errorf("EpisodicNodeType = %s, want episodic", EpisodicNodeType)
	}
	if CommunityNodeType != "community" {
		t.Errorf("CommunityNodeType = %s, want community", CommunityNodeType)
	}
	if SourceNodeType != "source" {
		t.Errorf("SourceNodeType = %s, want source", SourceNodeType)
	}
}

func TestEpisodeTypes(t *testing.T) {
	// Verify constant values haven't changed
	if ConversationEpisodeType != "conversation" {
		t.Errorf("ConversationEpisodeType = %s, want conversation", ConversationEpisodeType)
	}
	if DocumentEpisodeType != "document" {
		t.Errorf("DocumentEpisodeType = %s, want document", DocumentEpisodeType)
	}
	if EventEpisodeType != "event" {
		t.Errorf("EventEpisodeType = %s, want event", EventEpisodeType)
	}
}

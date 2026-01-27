package dto

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Validation errors
var (
	ErrEmptyGroupID      = errors.New("group_id cannot be empty")
	ErrEmptyMessages     = errors.New("messages cannot be empty")
	ErrEmptyName         = errors.New("name cannot be empty")
	ErrGroupIDTooLong    = errors.New("group_id exceeds maximum length (256)")
	ErrNameTooLong       = errors.New("name exceeds maximum length (1024)")
	ErrContentTooLong    = errors.New("content exceeds maximum length (1MB)")
	ErrInvalidCharacters = errors.New("field contains invalid characters")
)

// MaxFieldLengths defines maximum lengths for fields to prevent abuse
const (
	MaxGroupIDLength  = 256
	MaxNameLength     = 1024
	MaxContentLength  = 1024 * 1024 // 1MB
	MaxMessagesCount  = 1000
	MaxEntityType     = 256
	MaxAttributeCount = 100
)

// AddMessagesRequest represents a request to add messages to the knowledge graph
type AddMessagesRequest struct {
	GroupID   string     `json:"group_id" binding:"required"`
	Messages  []Message  `json:"messages" binding:"required,dive"`
	Reference *time.Time `json:"reference,omitempty"`
}

// Validate performs validation on AddMessagesRequest
func (r *AddMessagesRequest) Validate() error {
	if strings.TrimSpace(r.GroupID) == "" {
		return ErrEmptyGroupID
	}
	if len(r.GroupID) > MaxGroupIDLength {
		return ErrGroupIDTooLong
	}
	if len(r.Messages) == 0 {
		return ErrEmptyMessages
	}
	if len(r.Messages) > MaxMessagesCount {
		return errors.New("messages count exceeds maximum (1000)")
	}
	for i, msg := range r.Messages {
		if err := msg.Validate(); err != nil {
			return fmt.Errorf("message %d: %w", i, err)
		}
	}
	return nil
}

// AddEntityNodeRequest represents a request to add an entity node
type AddEntityNodeRequest struct {
	GroupID    string                 `json:"group_id" binding:"required"`
	Name       string                 `json:"name" binding:"required"`
	EntityType string                 `json:"entity_type,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// Validate performs validation on AddEntityNodeRequest
func (r *AddEntityNodeRequest) Validate() error {
	if strings.TrimSpace(r.GroupID) == "" {
		return ErrEmptyGroupID
	}
	if len(r.GroupID) > MaxGroupIDLength {
		return ErrGroupIDTooLong
	}
	if strings.TrimSpace(r.Name) == "" {
		return ErrEmptyName
	}
	if len(r.Name) > MaxNameLength {
		return ErrNameTooLong
	}
	if len(r.EntityType) > MaxEntityType {
		return errors.New("entity_type exceeds maximum length (256)")
	}
	if len(r.Attributes) > MaxAttributeCount {
		return errors.New("attributes count exceeds maximum (100)")
	}
	return nil
}

// ClearDataRequest represents a request to clear graph data
type ClearDataRequest struct {
	GroupIDs []string `json:"group_ids,omitempty"`
}

// Validate performs validation on ClearDataRequest
func (r *ClearDataRequest) Validate() error {
	for _, groupID := range r.GroupIDs {
		if len(groupID) > MaxGroupIDLength {
			return ErrGroupIDTooLong
		}
	}
	return nil
}

// IngestResponse represents a response from ingest operations
type IngestResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	ProcessID string `json:"process_id,omitempty"`
}

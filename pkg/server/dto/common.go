package dto

import (
	"errors"
	"strings"
	"time"
)

// Message represents a chat message
type Message struct {
	Role      string     `json:"role" binding:"required"`
	Content   string     `json:"content" binding:"required"`
	Timestamp *time.Time `json:"timestamp,omitempty"`
}

// ValidRoles defines acceptable message roles
var ValidRoles = map[string]bool{
	"user":      true,
	"assistant": true,
	"system":    true,
}

// Validate performs validation on Message
func (m *Message) Validate() error {
	if strings.TrimSpace(m.Role) == "" {
		return errors.New("role cannot be empty")
	}
	if !ValidRoles[strings.ToLower(m.Role)] {
		return errors.New("invalid role: must be user, assistant, or system")
	}
	if strings.TrimSpace(m.Content) == "" {
		return errors.New("content cannot be empty")
	}
	if len(m.Content) > MaxContentLength {
		return ErrContentTooLong
	}
	return nil
}

// Result represents a generic API result
type Result struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// FactResult represents a fact result from the knowledge graph
type FactResult struct {
	UUID         string     `json:"uuid"`
	Fact         string     `json:"fact"`
	SourceName   string     `json:"source_name"`
	TargetName   string     `json:"target_name"`
	RelationType string     `json:"relation_type"`
	ValidAt      *time.Time `json:"valid_at,omitempty"`
	InvalidAt    *time.Time `json:"invalid_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	Score        *float64   `json:"score,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

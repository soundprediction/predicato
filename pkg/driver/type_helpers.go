// Package driver provides safe type conversion helpers for Neo4j/Memgraph database types.
package driver

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

// TypeConversionError represents an error during type conversion from database types.
type TypeConversionError struct {
	Expected string
	Actual   string
	Field    string
}

func (e *TypeConversionError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("type conversion error for field %q: expected %s, got %s", e.Field, e.Expected, e.Actual)
	}
	return fmt.Sprintf("type conversion error: expected %s, got %s", e.Expected, e.Actual)
}

// NewTypeConversionError creates a new TypeConversionError.
func NewTypeConversionError(expected, actual, field string) *TypeConversionError {
	return &TypeConversionError{
		Expected: expected,
		Actual:   actual,
		Field:    field,
	}
}

// AsRecord safely converts an interface{} to *db.Record.
// Returns the record and true if successful, nil and false otherwise.
func AsRecord(v any) (*db.Record, bool) {
	if v == nil {
		return nil, false
	}
	record, ok := v.(*db.Record)
	return record, ok
}

// AsRecordSlice safely converts an interface{} to []*db.Record.
// Returns the slice and true if successful, nil and false otherwise.
func AsRecordSlice(v any) ([]*db.Record, bool) {
	if v == nil {
		return nil, false
	}
	records, ok := v.([]*db.Record)
	return records, ok
}

// AsDBNode safely converts an interface{} to dbtype.Node.
// Returns the node and true if successful, zero value and false otherwise.
func AsDBNode(v any) (dbtype.Node, bool) {
	if v == nil {
		return dbtype.Node{}, false
	}
	node, ok := v.(dbtype.Node)
	return node, ok
}

// AsDBRelationship safely converts an interface{} to dbtype.Relationship.
// Returns the relationship and true if successful, zero value and false otherwise.
func AsDBRelationship(v any) (dbtype.Relationship, bool) {
	if v == nil {
		return dbtype.Relationship{}, false
	}
	rel, ok := v.(dbtype.Relationship)
	return rel, ok
}

// AsString safely converts an interface{} to string.
// Returns the string and true if successful, empty string and false otherwise.
func AsString(v any) (string, bool) {
	if v == nil {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// AsInt64 safely converts an interface{} to int64.
// Returns the int64 and true if successful, 0 and false otherwise.
func AsInt64(v any) (int64, bool) {
	if v == nil {
		return 0, false
	}
	i, ok := v.(int64)
	return i, ok
}

// AsFloat64 safely converts an interface{} to float64.
// Returns the float64 and true if successful, 0 and false otherwise.
func AsFloat64(v any) (float64, bool) {
	if v == nil {
		return 0, false
	}
	f, ok := v.(float64)
	return f, ok
}

// AsBool safely converts an interface{} to bool.
// Returns the bool and true if successful, false and false otherwise.
func AsBool(v any) (bool, bool) {
	if v == nil {
		return false, false
	}
	b, ok := v.(bool)
	return b, ok
}

// AsStringSlice safely converts an interface{} to []string.
// Returns the slice and true if successful, nil and false otherwise.
func AsStringSlice(v any) ([]string, bool) {
	if v == nil {
		return nil, false
	}
	s, ok := v.([]string)
	return s, ok
}

// AsAnySlice safely converts an interface{} to []any.
// Returns the slice and true if successful, nil and false otherwise.
func AsAnySlice(v any) ([]any, bool) {
	if v == nil {
		return nil, false
	}
	s, ok := v.([]any)
	return s, ok
}

// AsMap safely converts an interface{} to map[string]any.
// Returns the map and true if successful, nil and false otherwise.
func AsMap(v any) (map[string]any, bool) {
	if v == nil {
		return nil, false
	}
	m, ok := v.(map[string]any)
	return m, ok
}

// MustRecord converts an interface{} to *db.Record or returns an error.
func MustRecord(v any, field string) (*db.Record, error) {
	record, ok := AsRecord(v)
	if !ok {
		return nil, NewTypeConversionError("*db.Record", fmt.Sprintf("%T", v), field)
	}
	return record, nil
}

// MustRecordSlice converts an interface{} to []*db.Record or returns an error.
func MustRecordSlice(v any, field string) ([]*db.Record, error) {
	records, ok := AsRecordSlice(v)
	if !ok {
		return nil, NewTypeConversionError("[]*db.Record", fmt.Sprintf("%T", v), field)
	}
	return records, nil
}

// MustDBNode converts an interface{} to dbtype.Node or returns an error.
func MustDBNode(v any, field string) (dbtype.Node, error) {
	node, ok := AsDBNode(v)
	if !ok {
		return dbtype.Node{}, NewTypeConversionError("dbtype.Node", fmt.Sprintf("%T", v), field)
	}
	return node, nil
}

// MustDBRelationship converts an interface{} to dbtype.Relationship or returns an error.
func MustDBRelationship(v any, field string) (dbtype.Relationship, error) {
	rel, ok := AsDBRelationship(v)
	if !ok {
		return dbtype.Relationship{}, NewTypeConversionError("dbtype.Relationship", fmt.Sprintf("%T", v), field)
	}
	return rel, nil
}

// MustString converts an interface{} to string or returns an error.
func MustString(v any, field string) (string, error) {
	s, ok := AsString(v)
	if !ok {
		return "", NewTypeConversionError("string", fmt.Sprintf("%T", v), field)
	}
	return s, nil
}

// MustInt64 converts an interface{} to int64 or returns an error.
func MustInt64(v any, field string) (int64, error) {
	i, ok := AsInt64(v)
	if !ok {
		return 0, NewTypeConversionError("int64", fmt.Sprintf("%T", v), field)
	}
	return i, nil
}

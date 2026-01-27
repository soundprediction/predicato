// Package types defines the core data types for the predicato knowledge graph.
//
// This package contains the fundamental types used throughout predicato:
//   - Node: Represents entities, episodes, communities, and sources in the graph
//   - Edge/EntityEdge: Represents relationships between nodes
//   - Episode: Represents temporal data units to be processed
//   - SearchConfig: Configuration for search operations
//
// # Node Types
//
// Nodes can be of several types:
//   - EntityNodeType: Entities extracted from content
//   - EpisodicNodeType: Episodic memories or events
//   - CommunityNodeType: Communities of related entities
//   - SourceNodeType: Source nodes where content originates
//
// # Validation
//
// Types provide Validate() and ValidateForCreate() methods for input validation:
//
//	node := &types.Node{Name: "test", GroupID: "group-1"}
//	if err := node.Validate(); err != nil {
//	    // Handle validation error
//	}
//
// # JSON Serialization
//
// All types are designed to be JSON-serializable with appropriate struct tags.
// Sensitive fields are excluded from JSON serialization where appropriate.
package types

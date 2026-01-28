package gliner2

import (
	"time"
)

type Entity struct {
	Text       string  `json:"text"`
	Label      string  `json:"label"`
	Confidence float64 `json:"confidence,omitempty"`
	Start      int     `json:"start,omitempty"`
	End        int     `json:"end,omitempty"`
}

type Relation struct {
	Head       string  `json:"head"`
	Tail       string  `json:"tail"`
	Relation   string  `json:"relation"`
	Confidence float64 `json:"confidence,omitempty"`
	HeadSpan   *Span   `json:"head_span,omitempty"`
	TailSpan   *Span   `json:"tail_span,omitempty"`
}

type Span struct {
	Text  string `json:"text"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

type Fact struct {
	Source     string  `json:"source"`
	Target     string  `json:"target"`
	Type       string  `json:"type"`
	Confidence float64 `json:"confidence,omitempty"`
	SourceSpan *Span   `json:"source_span,omitempty"`
	TargetSpan *Span   `json:"target_span,omitempty"`
}

type Classification struct {
	Task       string   `json:"task"`
	Label      string   `json:"label"`
	Confidence float64  `json:"confidence,omitempty"`
	Labels     []string `json:"labels,omitempty"` // For multi-label
}

// GLInER2 API request/response types
type ExtractRequest struct {
	Task      string      `json:"task"`
	Text      string      `json:"text"`
	Schema    interface{} `json:"schema"`
	Threshold float64     `json:"threshold,omitempty"`
}

type ExtractResponse struct {
	Result interface{} `json:"result"`
}

type EntityResult struct {
	Entities map[string][]Entity `json:"entities"`
}

type RelationResult struct {
	RelationExtraction map[string][]RelationTuple `json:"relation_extraction"`
}

type RelationTuple struct {
	Head string `json:"head"`
	Tail string `json:"tail"`
}

type ClassificationResult struct {
	Classifications map[string]Classification `json:"classifications"`
}

type StructuredResult struct {
	Structured map[string][]map[string]interface{} `json:"structured"`
}

// Health check
type HealthResponse struct {
	Status    string    `json:"status"`
	Models    []string  `json:"models,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

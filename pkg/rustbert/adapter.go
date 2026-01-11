package rustbert

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/soundprediction/predicato/pkg/nlp"
	"github.com/soundprediction/predicato/pkg/types"
)

// LLMAdapter adapts the RustBert client to the llm.Client interface.
type LLMAdapter struct {
	client *Client
	task   string // "summarization", "text_generation", "ner", "qa"
}

// NewLLMAdapter creates a new LLM adapter for RustBert.
func NewLLMAdapter(client *Client, task string) *LLMAdapter {
	return &LLMAdapter{
		client: client,
		task:   task,
	}
}

// Chat implements the llm.Client interface.
func (a *LLMAdapter) Chat(ctx context.Context, messages []types.Message) (*types.Response, error) {
	// Simple concatenation of user messages for now
	var inputBuilder strings.Builder
	for _, msg := range messages {
		if msg.Role == llm.RoleUser {
			if inputBuilder.Len() > 0 {
				inputBuilder.WriteString("\n")
			}
			inputBuilder.WriteString(msg.Content)
		}
	}
	input := inputBuilder.String()

	var output string
	var err error

	switch a.task {
	case "summarization":
		summaries, err := a.client.Summarize(input)
		if err != nil {
			return nil, err
		}
		if len(summaries) > 0 {
			output = summaries[0]
		}
	case "text_generation", "generation":
		output, err = a.client.GenerateText(input)
		if err != nil {
			return nil, err
		}
	case "ner":
		entities, err := a.client.ExtractEntities(input)
		if err != nil {
			return nil, err
		}
		// Serialize to JSON for standard return
		b, err := json.Marshal(entities)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal NER entities: %w", err)
		}
		output = string(b)
	case "qa":
		// Naive protocol: Input format "Question: <q>\nContext: <ctx>"
		// Or try to parse last message?
		// For now, let's assume input is JUST the context, and we can't do QA easily without structured input.
		// Or assume json input?
		return nil, fmt.Errorf("QA task not fully implemented in chat interface yet")
	default:
		return nil, fmt.Errorf("unknown task: %s", a.task)
	}

	return &types.Response{
		Content: output,
	}, nil
}

// ChatWithStructuredOutput implements the llm.Client interface.
func (a *LLMAdapter) ChatWithStructuredOutput(ctx context.Context, messages []types.Message, schema any) (*types.Response, error) {
	// RustBert models don't support structured output schema enforcement natively yet.
	// We just call chat and hope for the best (or the task inherently returns structure like NER).
	return a.Chat(ctx, messages)
}

// Close implements the llm.Client interface.
func (a *LLMAdapter) Close() error {
	// The underlying client is shared, so we generally don't close it here
	// unless we own it. Application logic typically manages the rustbert.Client lifecycle.
	return nil
}

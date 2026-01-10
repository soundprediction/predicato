package rustbert

import (
	"fmt"
	"sync"

	"github.com/soundprediction/go-rust-bert/pkg/rustbert"
)

// Client wraps go-rust-bert models for use in Predicato.
type Client struct {
	nerModel           *rustbert.NERModel
	summarizationModel *rustbert.SummarizationModel
	qaModel            *rustbert.QAModel
	textGenModel       *rustbert.TextGenerationModel
	mu                 sync.Mutex
}

// Config holds configuration for RustBert models
type Config struct {
	NERModelID           string
	SummarizationModelID string
}

// NewClient creates a new RustBert client.
// It initializes models lazily or based on provided IDs if we wanted,
// but here we'll provide methods to load them.
func NewClient() *Client {
	return &Client{}
}

// LoadNERModel loads the NER model.
// If modelID is empty, it uses the default (BERT-based).
// If modelID is a local path, it attempts to load from files (requires specific structure).
// Note: The underlying binding might expects explicit file paths for custom models.
// For simplicity, we'll expose loading default for now, or custom from files if full paths provided.
// To keep it simple like GLiNER integration:
func (c *Client) LoadNERModel() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.nerModel != nil {
		return nil
	}

	// Using default BERT NER model
	m, err := rustbert.NewNERModel()
	if err != nil {
		return fmt.Errorf("failed to create NER model: %w", err)
	}
	c.nerModel = m
	return nil
}

// LoadSummarizationModel loads the Summarization model (BART/DistilBART default).
func (c *Client) LoadSummarizationModel() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.summarizationModel != nil {
		return nil
	}

	m, err := rustbert.NewSummarizationModel()
	if err != nil {
		return fmt.Errorf("failed to create Summarization model: %w", err)
	}
	c.summarizationModel = m
	return nil
}

// Close closes all loaded models.
func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.nerModel != nil {
		c.nerModel.Close()
		c.nerModel = nil
	}
	if c.summarizationModel != nil {
		c.summarizationModel.Close()
		c.summarizationModel = nil
	}
}

// Entity represents an extracted entity
type Entity struct {
	Text  string
	Label string
	Score float64
}

// ExtractEntities extracts named entities from text.
// Auto-loads model if not loaded? Better to require explicit load or load on first use.
// Let's load on first use if not loaded.
func (c *Client) ExtractEntities(text string) ([]Entity, error) {
	if c.nerModel == nil {
		if err := c.LoadNERModel(); err != nil {
			return nil, err
		}
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	results, err := c.nerModel.Predict(text)
	if err != nil {
		return nil, fmt.Errorf("NER prediction failed: %w", err)
	}

	entities := make([]Entity, len(results))
	for i, r := range results {
		entities[i] = Entity{
			Text:  r.Word,
			Label: r.Label,
			Score: r.Score,
		}
	}
	return entities, nil
}

// Summarize generates a summary of the text.
func (c *Client) Summarize(text string) ([]string, error) {
	if c.summarizationModel == nil {
		if err := c.LoadSummarizationModel(); err != nil {
			return nil, err
		}
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Summarize returns a list of summaries (usually just 1 for single input)
	results, err := c.summarizationModel.Summarize(text)
	if err != nil {
		return nil, fmt.Errorf("summarization failed: %w", err)
	}

	return results, nil
}

// LoadQAModel loads the Question Answering model.
func (c *Client) LoadQAModel() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.qaModel != nil {
		return nil
	}

	m, err := rustbert.NewQAModel()
	if err != nil {
		return fmt.Errorf("failed to create QA model: %w", err)
	}
	c.qaModel = m
	return nil
}

// AnswerQuestion answers a question based on the context.
func (c *Client) AnswerQuestion(question, context string) ([]rustbert.Answer, error) {
	if c.qaModel == nil {
		if err := c.LoadQAModel(); err != nil {
			return nil, err
		}
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	results, err := c.qaModel.Predict(question, context)
	if err != nil {
		return nil, fmt.Errorf("QA prediction failed: %w", err)
	}

	return results, nil
}

// LoadTextGenerationModel loads the Text Generation model.
func (c *Client) LoadTextGenerationModel() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.textGenModel != nil {
		return nil
	}

	m, err := rustbert.NewTextGenerationModel()
	if err != nil {
		return fmt.Errorf("failed to create Text Generation model: %w", err)
	}
	c.textGenModel = m
	return nil
}

// GenerateText generates text from a prompt.
func (c *Client) GenerateText(prompt string) (string, error) {
	if c.textGenModel == nil {
		if err := c.LoadTextGenerationModel(); err != nil {
			return "", err
		}
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	result, err := c.textGenModel.Generate(prompt, "")
	if err != nil {
		return "", fmt.Errorf("text generation failed: %w", err)
	}

	return result, nil
}

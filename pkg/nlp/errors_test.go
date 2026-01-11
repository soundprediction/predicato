package nlp_test

import (
	"testing"

	"github.com/soundprediction/predicato/pkg/nlp"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitError(t *testing.T) {
	t.Run("default message", func(t *testing.T) {
		err := nlp.NewRateLimitError()
		assert.Equal(t, "rate limit exceeded. Please try again later", err.Error())
	})

	t.Run("custom message", func(t *testing.T) {
		customMessage := "Custom rate limit message"
		err := nlp.NewRateLimitError(customMessage)
		assert.Equal(t, customMessage, err.Error())
	})
}

func TestRefusalError(t *testing.T) {
	t.Run("message assignment", func(t *testing.T) {
		message := "The LLM refused to respond to this prompt."
		err := nlp.NewRefusalError(message)
		assert.Equal(t, message, err.Error())
	})
}

func TestEmptyResponseError(t *testing.T) {
	t.Run("message assignment", func(t *testing.T) {
		message := "The LLM returned an empty response."
		err := nlp.NewEmptyResponseError(message)
		assert.Equal(t, message, err.Error())
	})
}

func TestCommonErrors(t *testing.T) {
	t.Run("error constants", func(t *testing.T) {
		assert.NotNil(t, nlp.ErrRateLimit)
		assert.NotNil(t, nlp.ErrRefusal)
		assert.NotNil(t, nlp.ErrEmptyResponse)
		assert.NotNil(t, nlp.ErrInvalidModel)

		assert.Contains(t, nlp.ErrRateLimit.Error(), "rate limit")
		assert.Contains(t, nlp.ErrRefusal.Error(), "refused")
		assert.Contains(t, nlp.ErrEmptyResponse.Error(), "empty")
		assert.Contains(t, nlp.ErrInvalidModel.Error(), "invalid model")
	})
}

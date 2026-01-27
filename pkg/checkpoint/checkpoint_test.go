package checkpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/soundprediction/predicato/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckpointManager(t *testing.T) {
	// Create temporary directory for tests
	tmpDir, err := os.MkdirTemp("", "predicato-checkpoint-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	ctx := context.Background()

	t.Run("Create manager with custom directory", func(t *testing.T) {
		manager, err := NewCheckpointManager(tmpDir)
		require.NoError(t, err)
		assert.Equal(t, tmpDir, manager.GetCheckpointDir())
	})

	t.Run("Create manager with default directory", func(t *testing.T) {
		manager, err := NewCheckpointManager("")
		require.NoError(t, err)
		expectedDir := filepath.Join(os.TempDir(), "predicato-checkpoints")
		assert.Equal(t, expectedDir, manager.GetCheckpointDir())
	})

	t.Run("Save and load checkpoint", func(t *testing.T) {
		manager, err := NewCheckpointManager(tmpDir)
		require.NoError(t, err)

		checkpoint := &EpisodeCheckpoint{
			EpisodeID:     "episode-123",
			GroupID:       "group-456",
			Step:          StepPrepared,
			CreatedAt:     time.Now(),
			LastUpdatedAt: time.Now(),
			Episode: types.Episode{
				ID:      "episode-123",
				Name:    "Test Episode",
				Content: "Test content",
				GroupID: "group-456",
			},
			Chunks: []string{"chunk1", "chunk2"},
		}

		// Save checkpoint
		err = manager.Save(ctx, checkpoint)
		require.NoError(t, err)

		// Load checkpoint
		loaded, err := manager.Load(ctx, "episode-123")
		require.NoError(t, err)
		require.NotNil(t, loaded)

		assert.Equal(t, checkpoint.EpisodeID, loaded.EpisodeID)
		assert.Equal(t, checkpoint.GroupID, loaded.GroupID)
		assert.Equal(t, checkpoint.Step, loaded.Step)
		assert.Equal(t, checkpoint.Episode.ID, loaded.Episode.ID)
		assert.Equal(t, len(checkpoint.Chunks), len(loaded.Chunks))
	})

	t.Run("Load non-existent checkpoint", func(t *testing.T) {
		manager, err := NewCheckpointManager(tmpDir)
		require.NoError(t, err)

		loaded, err := manager.Load(ctx, "non-existent")
		require.NoError(t, err)
		assert.Nil(t, loaded)
	})

	t.Run("Delete checkpoint", func(t *testing.T) {
		manager, err := NewCheckpointManager(tmpDir)
		require.NoError(t, err)

		checkpoint := &EpisodeCheckpoint{
			EpisodeID: "episode-delete",
			GroupID:   "group-456",
			Step:      StepPrepared,
			CreatedAt: time.Now(),
		}

		// Save and verify exists
		err = manager.Save(ctx, checkpoint)
		require.NoError(t, err)

		exists, err := manager.Exists(ctx, "episode-delete")
		require.NoError(t, err)
		assert.True(t, exists)

		// Delete and verify doesn't exist
		err = manager.Delete(ctx, "episode-delete")
		require.NoError(t, err)

		exists, err = manager.Exists(ctx, "episode-delete")
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Update step", func(t *testing.T) {
		manager, err := NewCheckpointManager(tmpDir)
		require.NoError(t, err)

		checkpoint := &EpisodeCheckpoint{
			EpisodeID: "episode-step",
			GroupID:   "group-456",
			Step:      StepPrepared,
			CreatedAt: time.Now(),
		}

		err = manager.Save(ctx, checkpoint)
		require.NoError(t, err)

		// Update step
		err = manager.UpdateStep(ctx, "episode-step", StepExtractedEntities)
		require.NoError(t, err)

		// Verify updated
		loaded, err := manager.Load(ctx, "episode-step")
		require.NoError(t, err)
		assert.Equal(t, StepExtractedEntities, loaded.Step)
	})

	t.Run("Record error", func(t *testing.T) {
		manager, err := NewCheckpointManager(tmpDir)
		require.NoError(t, err)

		checkpoint := &EpisodeCheckpoint{
			EpisodeID:    "episode-error",
			GroupID:      "group-456",
			Step:         StepExtractedEntities,
			CreatedAt:    time.Now(),
			AttemptCount: 0,
		}

		err = manager.Save(ctx, checkpoint)
		require.NoError(t, err)

		// Record error
		testErr := assert.AnError
		err = manager.RecordError(ctx, "episode-error", testErr, "stack trace here")
		require.NoError(t, err)

		// Verify error recorded
		loaded, err := manager.Load(ctx, "episode-error")
		require.NoError(t, err)
		assert.Equal(t, 1, loaded.AttemptCount)
		assert.Contains(t, loaded.LastError, "assert.AnError")
		assert.Equal(t, "stack trace here", loaded.LastErrorStack)
	})

	t.Run("List checkpoints", func(t *testing.T) {
		manager, err := NewCheckpointManager(tmpDir)
		require.NoError(t, err)

		// Create multiple checkpoints
		for i := 0; i < 3; i++ {
			checkpoint := &EpisodeCheckpoint{
				EpisodeID: fmt.Sprintf("episode-list-%d", i),
				GroupID:   "group-456",
				Step:      StepPrepared,
				CreatedAt: time.Now(),
			}
			err = manager.Save(ctx, checkpoint)
			require.NoError(t, err)
		}

		// List all
		checkpoints, err := manager.List(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(checkpoints), 3)
	})

	t.Run("Clean old checkpoints", func(t *testing.T) {
		manager, err := NewCheckpointManager(tmpDir)
		require.NoError(t, err)

		// Create old checkpoint - manually write with old timestamp
		oldTime := time.Now().Add(-48 * time.Hour)
		oldCheckpoint := &EpisodeCheckpoint{
			EpisodeID:     "episode-old",
			GroupID:       "group-456",
			Step:          StepPrepared,
			CreatedAt:     oldTime,
			LastUpdatedAt: oldTime,
		}
		// Manually write to preserve old timestamp
		data, err := json.MarshalIndent(oldCheckpoint, "", "  ")
		require.NoError(t, err)
		oldPath, err := manager.GetCheckpointPath("episode-old")
		require.NoError(t, err)
		err = os.WriteFile(oldPath, data, 0644)
		require.NoError(t, err)

		// Create recent checkpoint
		recentCheckpoint := &EpisodeCheckpoint{
			EpisodeID:     "episode-recent",
			GroupID:       "group-456",
			Step:          StepPrepared,
			CreatedAt:     time.Now(),
			LastUpdatedAt: time.Now(),
		}
		err = manager.Save(ctx, recentCheckpoint)
		require.NoError(t, err)

		// Clean old (older than 24 hours)
		removed, err := manager.CleanOld(ctx, 24*time.Hour)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, removed, 1)

		// Verify old checkpoint is gone but recent remains
		exists, err := manager.Exists(ctx, "episode-old")
		require.NoError(t, err)
		assert.False(t, exists)

		exists, err = manager.Exists(ctx, "episode-recent")
		require.NoError(t, err)
		assert.True(t, exists)
	})
}

func TestPathTraversalPrevention(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "predicato-checkpoint-security-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	ctx := context.Background()
	manager, err := NewCheckpointManager(tmpDir)
	require.NoError(t, err)

	// Create a sensitive file outside the checkpoint directory to verify it can't be accessed
	sensitiveFile := filepath.Join(tmpDir, "..", "sensitive.txt")
	err = os.WriteFile(sensitiveFile, []byte("sensitive data"), 0644)
	require.NoError(t, err)
	defer os.Remove(sensitiveFile)

	pathTraversalAttempts := []struct {
		name      string
		episodeID string
	}{
		{"simple path traversal", "../../../etc/passwd"},
		{"path traversal with dots", ".."},
		{"double traversal", "foo/../.."},
		{"forward slash", "foo/bar"},
		{"backslash", `foo\bar`},
		{"null byte", "episode\x00.json"},
		{"hidden file traversal", "../.hidden"},
		{"absolute path attempt", "/etc/passwd"},
		{"windows path", `C:\Windows\System32`},
		{"empty ID", ""},
	}

	for _, tc := range pathTraversalAttempts {
		t.Run("GetCheckpointPath_"+tc.name, func(t *testing.T) {
			_, err := manager.GetCheckpointPath(tc.episodeID)
			assert.ErrorIs(t, err, ErrInvalidEpisodeID, "Episode ID %q should be rejected", tc.episodeID)
		})

		t.Run("Load_"+tc.name, func(t *testing.T) {
			_, err := manager.Load(ctx, tc.episodeID)
			assert.Error(t, err, "Load should reject episode ID %q", tc.episodeID)
		})

		t.Run("Delete_"+tc.name, func(t *testing.T) {
			err := manager.Delete(ctx, tc.episodeID)
			assert.Error(t, err, "Delete should reject episode ID %q", tc.episodeID)
		})

		t.Run("Exists_"+tc.name, func(t *testing.T) {
			_, err := manager.Exists(ctx, tc.episodeID)
			assert.Error(t, err, "Exists should reject episode ID %q", tc.episodeID)
		})

		t.Run("Save_"+tc.name, func(t *testing.T) {
			checkpoint := &EpisodeCheckpoint{
				EpisodeID: tc.episodeID,
				GroupID:   "test-group",
				Step:      StepInitial,
			}
			err := manager.Save(ctx, checkpoint)
			assert.Error(t, err, "Save should reject episode ID %q", tc.episodeID)
		})
	}

	// Test that valid episode IDs still work
	validIDs := []string{
		"episode-123",
		"my_episode",
		"Episode.With.Dots",
		"episode-2024-01-15T10:30:00Z",
		"abc123def456",
		"a",
	}

	for _, id := range validIDs {
		t.Run("valid_ID_"+id, func(t *testing.T) {
			path, err := manager.GetCheckpointPath(id)
			require.NoError(t, err, "Valid episode ID %q should be accepted", id)
			assert.Contains(t, path, id, "Path should contain the episode ID")
			assert.True(t, filepath.IsAbs(path) || filepath.HasPrefix(path, tmpDir),
				"Path should be within checkpoint directory")
		})
	}
}

func TestProcessingSteps(t *testing.T) {
	steps := []ProcessingStep{
		StepInitial,
		StepPrepared,
		StepGotPreviousEpisodes,
		StepCreatedChunks,
		StepExtractedEntities,
		StepDeduplicatedEntities,
		StepExtractedEdges,
		StepResolvedEdges,
		StepExtractedAttributes,
		StepBuiltEpisodicEdges,
		StepPerformedGraphUpdate,
		StepUpdatedCommunities,
		StepCompleted,
	}

	// Verify all steps are unique
	stepMap := make(map[ProcessingStep]bool)
	for _, step := range steps {
		assert.False(t, stepMap[step], "Duplicate step: %s", step)
		stepMap[step] = true
	}
}

package theme

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"
)

const stateFileName = ".agentrc.state.json"

// RateLimitState is the schema of .agentrc.state.json.
type RateLimitState struct {
	UpdatedAt            time.Time `json:"updated_at"`
	RequestsLimit        int       `json:"requests_limit"`
	RequestsRemaining    int       `json:"requests_remaining"`
	RequestsResetAt      string    `json:"requests_reset_at"`
	TokensLimit          int       `json:"tokens_limit"`
	TokensRemaining      int       `json:"tokens_remaining"`
	TokensResetAt        string    `json:"tokens_reset_at"`
}

// RateLimitWidget shows Anthropic API rate limit usage.
// Priority: ClaudeInput stdin → .agentrc.state.json fallback.
type RateLimitWidget struct {
	Dir   string
	Input *ClaudeInput
}

func (w *RateLimitWidget) Name() string { return "rate limit" }

func (w *RateLimitWidget) Fetch(ctx context.Context) Result {
	// Priority 1: data from Claude Code stdin
	if w.Input != nil {
		fiveHour := w.Input.RateLimits.FiveHour.UsedPercentage
		if fiveHour > 0 {
			return Result{Value: fmt.Sprintf("%.0f%%", math.Round(fiveHour))}
		}
	}

	// Priority 2: read from state file
	state, err := loadRateLimitState(w.Dir)
	if err != nil || state == nil {
		return Result{Value: "(no data)"}
	}

	if state.RequestsLimit == 0 {
		return Result{Value: "(no data)"}
	}

	remaining := math.Round(float64(state.RequestsRemaining) / float64(state.RequestsLimit) * 100)
	used := math.Round(100 - remaining)
	stale := ""
	if time.Since(state.UpdatedAt) > 10*time.Minute {
		stale = " (stale)"
	}
	return Result{Value: fmt.Sprintf("%.0f%%%s", used, stale)}
}

func loadRateLimitState(dir string) (*RateLimitState, error) {
	path := filepath.Join(dir, stateFileName)
	data, err := os.ReadFile(path) // #nosec G304
	if err != nil {
		return nil, err
	}
	var state RateLimitState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

// SaveRateLimitState writes state to .agentrc.state.json in dir.
// Called by external API-calling code to persist rate limit headers.
func SaveRateLimitState(dir string, state *RateLimitState) error {
	state.UpdatedAt = time.Now()
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, stateFileName), data, 0600)
}

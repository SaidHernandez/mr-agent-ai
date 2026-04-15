package theme

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const claudeInputCacheFile = ".agentrc.input.json"

// ClaudeInput holds the JSON context that Claude Code passes via stdin
// when running a statusLine command.
type ClaudeInput struct {
	CWD       string `json:"cwd"`
	Workspace struct {
		ProjectDir string `json:"project_dir"`
	} `json:"workspace"`
	Cost struct {
		TotalCostUSD float64 `json:"total_cost_usd"`
	} `json:"cost"`
	ContextWindow struct {
		UsedPercentage   float64 `json:"used_percentage"`
		TotalInputTokens int     `json:"total_input_tokens"`
	} `json:"context_window"`
	RateLimits struct {
		FiveHour struct {
			UsedPercentage float64 `json:"used_percentage"`
		} `json:"five_hour"`
		SevenDay struct {
			UsedPercentage float64 `json:"used_percentage"`
		} `json:"seven_day"`
	} `json:"rate_limits"`
}

// ReadClaudeInput reads JSON from stdin if it is a pipe (not a TTY).
// Returns nil without error when running in a regular terminal.
func ReadClaudeInput() (*ClaudeInput, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return nil, nil
	}
	// If stdin is a character device (TTY), there is no piped input.
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return nil, nil
	}
	var input ClaudeInput
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		return nil, nil // silently ignore malformed input
	}
	return &input, nil
}

// persistClaudeInput saves ClaudeInput to disk so manual invocations can reuse it.
func persistClaudeInput(dir string, input *ClaudeInput) {
	data, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(filepath.Join(dir, claudeInputCacheFile), data, 0600)
}

// loadCachedClaudeInput reads the last persisted ClaudeInput from disk.
func loadCachedClaudeInput(dir string) (*ClaudeInput, error) {
	data, err := os.ReadFile(filepath.Join(dir, claudeInputCacheFile)) // #nosec G304
	if err != nil {
		return nil, err
	}
	var input ClaudeInput
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, err
	}
	return &input, nil
}

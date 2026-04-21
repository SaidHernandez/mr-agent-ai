package theme

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const tokensCacheFile = ".agentrc.tokens.json" // #nosec G101

type tokensCache struct {
	CachedAt     time.Time `json:"cached_at"`
	InputTokens  int       `json:"input_tokens"`
	OutputTokens int       `json:"output_tokens"`
	TotalCostUSD float64   `json:"total_cost_usd"`
}

// TokensWidget shows token usage and cost.
// Priority: ClaudeInput stdin → .agentrc.tokens.json cache.
type TokensWidget struct {
	Dir   string
	Input *ClaudeInput
}

func (w *TokensWidget) Name() string { return "token" }

func (w *TokensWidget) Fetch(ctx context.Context) Result {
	// Priority 1: data from Claude Code stdin
	if w.Input != nil {
		cost := w.Input.Cost.TotalCostUSD
		total := w.Input.ContextWindow.TotalInputTokens
		if cost > 0 || total > 0 {
			return Result{Value: fmt.Sprintf("%s (~$%.2f)", formatTokens(total), cost)}
		}
	}

	// Priority 2: read from local cache
	cache, err := loadTokensCache(w.Dir)
	if err != nil || cache == nil {
		return Result{Value: "(no data)"}
	}

	suffix := ""
	if time.Since(cache.CachedAt) > time.Hour {
		suffix = " (stale)"
	}

	total := cache.InputTokens + cache.OutputTokens
	return Result{Value: fmt.Sprintf("%s (~$%.2f)%s", formatTokens(total), cache.TotalCostUSD, suffix)}
}

// ContextWidget shows the context window usage percentage.
type ContextWidget struct {
	Input *ClaudeInput
}

func (w *ContextWidget) Name() string { return "ctx" }

func (w *ContextWidget) Fetch(ctx context.Context) Result {
	if w.Input != nil {
		return Result{Value: fmt.Sprintf("%.0f%%", w.Input.ContextWindow.UsedPercentage)}
	}
	return Result{Value: "(no data)"}
}

func loadTokensCache(dir string) (*tokensCache, error) {
	data, err := os.ReadFile(filepath.Join(dir, tokensCacheFile)) // #nosec G304
	if err != nil {
		return nil, err
	}
	var c tokensCache
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func formatTokens(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	result := ""
	for n > 0 {
		group := n % 1000
		n /= 1000
		if n > 0 {
			result = fmt.Sprintf(",%03d", group) + result
		} else {
			result = fmt.Sprintf("%d", group) + result
		}
	}
	return result
}

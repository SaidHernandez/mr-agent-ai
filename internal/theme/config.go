package theme

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const configFileName = ".agentrc.yml"

// Config holds all theme configuration from .agentrc.yml.
type Config struct {
	Widgets WidgetsConfig `yaml:"widgets"`
}

// WidgetsConfig holds per-widget configuration.
type WidgetsConfig struct {
	Branch    BranchWidgetConfig    `yaml:"branch"`
	PR        PRWidgetConfig        `yaml:"pr"`
	Linear    LinearWidgetConfig    `yaml:"linear"`
	Tokens    TokensWidgetConfig    `yaml:"tokens"`
	RateLimit RateLimitWidgetConfig `yaml:"rate-limit"`
}

// BranchWidgetConfig configures the branch widget.
type BranchWidgetConfig struct {
	Enabled bool `yaml:"enabled"`
}

// PRWidgetConfig configures the GitHub PR widget.
type PRWidgetConfig struct {
	Enabled bool   `yaml:"enabled"`
	Repo    string `yaml:"repo"` // "owner/repo", auto-detected if empty
}

// LinearWidgetConfig configures the Linear ticket widget.
type LinearWidgetConfig struct {
	Enabled bool   `yaml:"enabled"`
	Prefix  string `yaml:"prefix"` // e.g. "ENG"
}

// TokensWidgetConfig configures the token usage widget.
type TokensWidgetConfig struct {
	Enabled bool `yaml:"enabled"`
}

// RateLimitWidgetConfig configures the rate limit widget.
type RateLimitWidgetConfig struct {
	Enabled bool `yaml:"enabled"`
}

// DefaultConfig returns a Config with all widgets enabled and sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Widgets: WidgetsConfig{
			Branch:    BranchWidgetConfig{Enabled: true},
			PR:        PRWidgetConfig{Enabled: true},
			Linear:    LinearWidgetConfig{Enabled: true},
			Tokens:    TokensWidgetConfig{Enabled: true},
			RateLimit: RateLimitWidgetConfig{Enabled: true},
		},
	}
}

// ConfigPath returns the absolute path to .agentrc.yml in dir.
func ConfigPath(dir string) string {
	return filepath.Join(dir, configFileName)
}

// LoadConfig reads .agentrc.yml from dir. Returns DefaultConfig if not found.
func LoadConfig(dir string) (*Config, error) {
	path := ConfigPath(dir)
	data, err := os.ReadFile(path) // #nosec G304
	if os.IsNotExist(err) {
		return DefaultConfig(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", configFileName, err)
	}
	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", configFileName, err)
	}
	return cfg, nil
}

// SaveConfig writes cfg to .agentrc.yml in dir.
func SaveConfig(dir string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	path := ConfigPath(dir)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing %s: %w", configFileName, err)
	}
	return nil
}

// buildWidgets constructs the active widget list based on config and ClaudeInput.
func buildWidgets(dir string, cfg *Config, input *ClaudeInput) []Widget {
	var widgets []Widget
	if cfg.Widgets.Branch.Enabled {
		widgets = append(widgets, &BranchWidget{Dir: dir})
	}
	if cfg.Widgets.PR.Enabled {
		widgets = append(widgets, &PRWidget{Dir: dir, Config: cfg.Widgets.PR})
	}
	if cfg.Widgets.Linear.Enabled {
		widgets = append(widgets, &LinearWidget{Dir: dir, Config: cfg.Widgets.Linear})
	}
	if cfg.Widgets.RateLimit.Enabled {
		widgets = append(widgets, &RateLimitWidget{Dir: dir, Input: input})
	}
	if cfg.Widgets.Tokens.Enabled {
		widgets = append(widgets, &TokensWidget{Dir: dir, Input: input})
		widgets = append(widgets, &ContextWidget{Input: input})
	}
	return widgets
}

// Show renders the context panel. If statusLine is true, outputs a compact
// single line suitable for Claude Code's status line.
func Show(dir string, statusLine bool) error {
	cfg, err := LoadConfig(dir)
	if err != nil {
		return err
	}
	input, _ := ReadClaudeInput()
	if input != nil {
		persistClaudeInput(dir, input)
	} else {
		input, _ = loadCachedClaudeInput(dir)
	}
	widgets := buildWidgets(dir, cfg, input)
	results := RunAll(widgets)
	if statusLine {
		fmt.Print(RenderStatusLine(results, termWidth()))
	} else {
		fmt.Print(RenderPanel(results))
	}
	return nil
}

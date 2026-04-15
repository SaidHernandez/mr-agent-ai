package theme

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// configureClaude adds the statusLine field to ~/.claude/settings.json.
// Preserves all existing fields in the file.
func configureClaude(_ string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot determine home directory: %w", err)
	}

	claudeDir := filepath.Join(home, ".claude")
	settingsPath := filepath.Join(claudeDir, "settings.json")

	// Read existing settings or start fresh
	settings := map[string]any{}
	if data, err := os.ReadFile(settingsPath); err == nil { // #nosec G304
		if err := json.Unmarshal(data, &settings); err != nil {
			return fmt.Errorf("parsing ~/.claude/settings.json: %w", err)
		}
	}

	// Set or replace the statusLine field
	settings["statusLine"] = map[string]any{
		"type":            "command",
		"command":         "mr-agent-ai theme show --format=statusline",
		"refreshInterval": 30,
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling settings: %w", err)
	}

	if err := os.MkdirAll(claudeDir, 0750); err != nil {
		return fmt.Errorf("creating ~/.claude: %w", err)
	}

	if err := os.WriteFile(settingsPath, append(data, '\n'), 0600); err != nil {
		return fmt.Errorf("writing ~/.claude/settings.json: %w", err)
	}

	return nil
}

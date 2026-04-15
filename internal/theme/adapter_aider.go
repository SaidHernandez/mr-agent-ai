package theme

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const aiderScriptContent = `#!/bin/bash
# mr-agent-ai context panel — runs at Aider session start
mr-agent-ai theme show
`

const aiderConfEntry = "\n# mr-agent-ai context widget\nscript_file: ~/.config/mr-agent-ai/aider-context.sh\n"

// configureAider sets up an Aider startup script and adds it to ~/.aider.conf.yml.
func configureAider(_ string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot determine home directory: %w", err)
	}

	// Create the script directory and file
	scriptDir := filepath.Join(home, ".config", "mr-agent-ai")
	if err := os.MkdirAll(scriptDir, 0750); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	scriptPath := filepath.Join(scriptDir, "aider-context.sh")
	if err := os.WriteFile(scriptPath, []byte(aiderScriptContent), 0750); err != nil { // #nosec G306 -- script needs execute bit
		return fmt.Errorf("writing aider-context.sh: %w", err)
	}

	// Update or create ~/.aider.conf.yml
	confPath := filepath.Join(home, ".aider.conf.yml")
	existing := ""
	if data, err := os.ReadFile(confPath); err == nil { // #nosec G304
		existing = string(data)
	}

	// Avoid duplicate entries
	if strings.Contains(existing, "mr-agent-ai context widget") {
		return nil
	}

	content := existing + aiderConfEntry
	if err := os.WriteFile(confPath, []byte(content), 0600); err != nil { // #nosec G703
		return fmt.Errorf("writing ~/.aider.conf.yml: %w", err)
	}

	return nil
}

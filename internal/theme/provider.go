package theme

import (
	"os"
	"os/exec"
	"path/filepath"
)

// Provider represents an AI agent that can be configured with context widgets.
type Provider struct {
	Name        string
	Description string
	Detect      func() bool
	Configure   func(dir string) error
}

// Providers is the list of all supported AI agent integrations.
var Providers = []Provider{
	{
		Name:        "Claude Code",
		Description: "Anthropic CLI (~/.claude)",
		Detect:      detectClaude,
		Configure:   configureClaude,
	},
	{
		Name:        "Aider",
		Description: "AI pair programmer (~/.aider.conf.yml)",
		Detect:      detectAider,
		Configure:   configureAider,
	},
	{
		Name:        "VS Code / Copilot",
		Description: "GitHub Copilot in VS Code (.vscode/tasks.json)",
		Detect:      detectVSCode,
		Configure:   configureVSCode,
	},
}

// DetectInstalled returns only the providers found on the current system.
func DetectInstalled() []Provider {
	var found []Provider
	for _, p := range Providers {
		if p.Detect() {
			found = append(found, p)
		}
	}
	return found
}

func detectClaude() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	_, err = os.Stat(filepath.Join(home, ".claude"))
	return err == nil
}

func detectAider() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	if _, err := os.Stat(filepath.Join(home, ".aider.conf.yml")); err == nil {
		return true
	}
	_, err = exec.LookPath("aider")
	return err == nil
}

func detectVSCode() bool {
	if _, err := exec.LookPath("code"); err == nil {
		return true
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	_, err = os.Stat(filepath.Join(home, ".vscode"))
	return err == nil
}

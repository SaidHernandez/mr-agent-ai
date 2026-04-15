package theme

import (
	"fmt"
	"os"

	"mr-agent-ai/internal/ui"
)

// InitConfig runs the interactive setup flow:
// 1. Creates .agentrc.yml in dir (if not already present)
// 2. Detects installed AI agents
// 3. Prompts the user to select which ones to configure
// 4. Applies the corresponding adapter for each selected provider
func InitConfig(dir string) error {
	fmt.Println()
	fmt.Println("  Detecting installed agents...")
	fmt.Println()

	installed := DetectInstalled()
	all := Providers

	for _, p := range all {
		found := false
		for _, i := range installed {
			if i.Name == p.Name {
				found = true
				break
			}
		}
		if found {
			fmt.Printf("  \033[32m✓\033[0m %-18s found\n", p.Name)
		} else {
			fmt.Printf("  \033[90m✗\033[0m %-18s not found\n", p.Name)
		}
	}
	fmt.Println()

	if len(installed) == 0 {
		fmt.Println("  No supported agents detected.")
		fmt.Println("  You can still run:  mr-agent-ai theme show")
		fmt.Println()
		return createConfigFile(dir)
	}

	items := make([]string, len(installed))
	for i, p := range installed {
		items[i] = fmt.Sprintf("%-18s— %s", p.Name, p.Description)
	}

	selected := ui.Select("  Which agents to configure?", items)
	if len(selected) == 0 {
		fmt.Println("\n  No agents selected. Skipping configuration.")
		return createConfigFile(dir)
	}

	fmt.Println()

	if err := createConfigFile(dir); err != nil {
		return err
	}

	for _, idx := range selected {
		p := installed[idx]
		if err := p.Configure(dir); err != nil {
			fmt.Printf("  \033[31m✗\033[0m %-18s %v\n", p.Name, err)
		} else {
			fmt.Printf("  \033[32m✓\033[0m %-18s configured\n", p.Name)
		}
	}

	fmt.Println()
	fmt.Println("  Done! Run \033[1mmr-agent-ai theme show\033[0m to preview your panel.")
	fmt.Println()
	return nil
}

// createConfigFile writes .agentrc.yml if it does not already exist.
func createConfigFile(dir string) error {
	path := ConfigPath(dir)
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("  \033[33m~\033[0m %-18s already exists (skipped)\n", configFileName)
		return nil
	}
	if err := SaveConfig(dir, DefaultConfig()); err != nil {
		return err
	}
	fmt.Printf("  \033[32m✓\033[0m %-18s created\n", configFileName)
	return nil
}

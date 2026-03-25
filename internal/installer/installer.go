package installer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// Install runs the interactive installer: selects tools, selects skills, writes files.
func Install(targetDir string) error {
	printBanner()

	reader := bufio.NewReader(os.Stdin)

	// Step 1: select AI coding tools
	toolLabels := make([]string, len(agentTools))
	for i, t := range agentTools {
		toolLabels[i] = fmt.Sprintf("%-14s — %s", t.Name, t.Description)
	}
	toolIndices := promptSelect(reader, "Select your AI coding tools:", toolLabels)
	if len(toolIndices) == 0 {
		fmt.Println("\n[mr-agent-ai] No tools selected. Nothing installed.")
		return nil
	}
	var selectedTools []AgentTool
	for _, i := range toolIndices {
		selectedTools = append(selectedTools, agentTools[i])
	}

	// Step 2: select skills
	agentLabels := make([]string, len(agents))
	for i, a := range agents {
		agentLabels[i] = fmt.Sprintf("%-18s — %s", a.Name, a.Description)
	}
	skillIndices := promptSelect(reader, "\nSelect which skills to install:", agentLabels)
	if len(skillIndices) == 0 {
		fmt.Println("\n[mr-agent-ai] No skills selected. Nothing installed.")
		return nil
	}
	var selected []Agent
	for _, i := range skillIndices {
		selected = append(selected, agents[i])
	}

	fmt.Printf("\n[mr-agent-ai] Installing %d skill(s) for %d tool(s) into: %s\n\n",
		len(selected), len(selectedTools), targetDir)

	// Step 3: write skills/ directory — source of truth for all tools
	if err := writeSkills(targetDir, selected); err != nil {
		return err
	}

	// Step 4: generate tool-specific config files
	fmt.Println()
	for _, t := range selectedTools {
		if err := t.Generate(targetDir, selected); err != nil {
			return fmt.Errorf("failed to generate config for %s: %w", t.Name, err)
		}
	}

	fmt.Printf("\n[mr-agent-ai] Done. %d skill(s) installed for %d tool(s).\n",
		len(selected), len(selectedTools))
	fmt.Println("\nNext: tell your AI agent to read the generated config file before writing any code.")
	return nil
}

// writeSkills writes each selected skill's SKILL.md and assets to skills/<name>/.
func writeSkills(targetDir string, selected []Agent) error {
	for _, a := range selected {
		skillDir := filepath.Join(targetDir, "skills", a.Dir)
		if err := os.MkdirAll(skillDir, 0755); err != nil {
			return fmt.Errorf("failed to create skills/%s: %w", a.Dir, err)
		}
		skillPath := filepath.Join(skillDir, "SKILL.md")
		if err := os.WriteFile(skillPath, []byte(a.Content()), 0644); err != nil {
			return fmt.Errorf("failed to write skills/%s/SKILL.md: %w", a.Dir, err)
		}
		fmt.Printf("  [ok] skills/%s/SKILL.md\n", a.Dir)

		for _, asset := range a.Assets {
			assetDir := filepath.Join(skillDir, "assets")
			if err := os.MkdirAll(assetDir, 0755); err != nil {
				return fmt.Errorf("failed to create assets dir for %s: %w", a.Dir, err)
			}
			assetPath := filepath.Join(assetDir, asset.Path)
			if err := os.WriteFile(assetPath, []byte(asset.Content), 0644); err != nil {
				return fmt.Errorf("failed to write asset %s: %w", asset.Path, err)
			}
			fmt.Printf("  [ok] skills/%s/assets/%s\n", a.Dir, asset.Path)
		}
	}
	return nil
}

func printBanner() {
	fmt.Print(`
  __  __ ____        _                    _         _    ___
 |  \/  |  _ \      / \   __ _  ___ _ __ | |_      / \  |_ _|
 | |\/| | |_) |    / _ \ / _` + "`" + `\|/ _ \ '_ \| __|    / _ \  | |
 | |  | |  _ <    / ___ \ (_| |  __/ | | | |_    / ___ \ | |
 |_|  |_|_| \_\  /_/   \_\__, |\___|_| |_|\__|  /_/   \_\___|
                          |___/
  Multi-Agent Skill Installer
` + "\n")
}

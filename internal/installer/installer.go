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
	toolLabels := make([]string, len(tools))
	for i, t := range tools {
		toolLabels[i] = fmt.Sprintf("%-14s — %s", t.Name, t.Description)
	}
	toolIndices := promptSelect(reader, "Select your AI coding tools:", toolLabels)
	if len(toolIndices) == 0 {
		fmt.Println("\n[mr-agent-ai] No tools selected. Nothing installed.")
		return nil
	}
	var selectedTools []Tool
	for _, i := range toolIndices {
		selectedTools = append(selectedTools, tools[i])
	}

	// Step 2: select architecture skills
	archLabels := make([]string, len(archSkills))
	for i, s := range archSkills {
		archLabels[i] = fmt.Sprintf("%-18s — %s", s.Name, s.Description)
	}
	archIndices := promptSelect(reader, "\nSelect which architecture skills to install:", archLabels)
	if len(archIndices) == 0 {
		fmt.Println("\n[mr-agent-ai] No skills selected. Nothing installed.")
		return nil
	}
	var selected []Skill
	for _, i := range archIndices {
		selected = append(selected, archSkills[i])
	}

	// Step 3: select programming language skill (single choice, optional)
	langLabels := make([]string, len(langSkills))
	for i, l := range langSkills {
		langLabels[i] = fmt.Sprintf("%-14s — %s", l.Name, l.Description)
	}
	langIdx := promptSelectOne(reader, "\nSelect a programming language skill (optional):", langLabels)
	if langIdx >= 0 {
		selected = append(selected, langSkills[langIdx])
	}

	fmt.Printf("\n[mr-agent-ai] Installing %d skill(s) for %d tool(s) into: %s\n\n",
		len(selected), len(selectedTools), targetDir)

	// Step 4: write skills/ directory — source of truth for all tools
	if err := writeSkills(targetDir, selected); err != nil {
		return err
	}

	// Step 5: generate tool-specific config files
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

// writeSkills writes each selected skill's SKILL.md and assets to skills/<dir>/.
func writeSkills(targetDir string, selected []Skill) error {
	for _, s := range selected {
		skillDir := filepath.Join(targetDir, "skills", s.Dir)
		if err := os.MkdirAll(skillDir, 0o750); err != nil {
			return fmt.Errorf("failed to create skills/%s: %w", s.Dir, err)
		}
		skillPath := filepath.Join(skillDir, "SKILL.md")
		if err := os.WriteFile(skillPath, []byte(s.Content()), 0644); err != nil { // #nosec G306 -- project files are intentionally world-readable
			return fmt.Errorf("failed to write skills/%s/SKILL.md: %w", s.Dir, err)
		}
		fmt.Printf("  [ok] skills/%s/SKILL.md\n", s.Dir)

		for _, asset := range s.Assets {
			assetDir := filepath.Join(skillDir, "assets")
			if err := os.MkdirAll(assetDir, 0o750); err != nil {
				return fmt.Errorf("failed to create assets dir for %s: %w", s.Dir, err)
			}
			assetPath := filepath.Join(assetDir, asset.Path)
			if err := os.WriteFile(assetPath, []byte(asset.Content), 0644); err != nil { // #nosec G306 -- project files are intentionally world-readable
				return fmt.Errorf("failed to write asset %s: %w", asset.Path, err)
			}
			fmt.Printf("  [ok] skills/%s/assets/%s\n", s.Dir, asset.Path)
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

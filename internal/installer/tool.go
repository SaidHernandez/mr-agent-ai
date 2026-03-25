package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AgentTool represents a supported AI coding tool that can consume the installed skills.
type AgentTool struct {
	Name        string
	Description string
	Generate    func(targetDir string, selected []Agent) error
}

var agentTools = []AgentTool{
	{
		Name:        "Claude Code",
		Description: "Anthropic CLI — reads AGENTS.md + skills/ automatically",
		Generate:    generateClaudeCode,
	},
	{
		Name:        "Cursor",
		Description: "AI editor — .cursor/rules/<skill>.mdc per skill",
		Generate:    generateCursor,
	},
	{
		Name:        "GitHub Copilot",
		Description: "VS Code / JetBrains — .github/copilot-instructions.md",
		Generate:    generateCopilot,
	},
	{
		Name:        "Windsurf",
		Description: "Codeium editor — .windsurfrules",
		Generate:    generateWindsurf,
	},
	{
		Name:        "Cline",
		Description: "VS Code extension — .clinerules",
		Generate:    generateCline,
	},
	{
		Name:        "Aider",
		Description: "Terminal pair programmer — CONVENTIONS.md",
		Generate:    generateAider,
	},
	{
		Name:        "Continue",
		Description: "Open-source VS Code / JetBrains — .continuerules",
		Generate:    generateContinue,
	},
	{
		Name:        "OpenCode",
		Description: "SST terminal AI coder — AGENTS.md compatible",
		Generate:    generateOpenCode,
	},
}

// ─── Generators ───────────────────────────────────────────────────────────────

func generateClaudeCode(targetDir string, selected []Agent) error {
	path := filepath.Join(targetDir, "AGENTS.md")
	if err := os.WriteFile(path, []byte(agentsIndex(selected)), 0644); err != nil {
		return fmt.Errorf("failed to write AGENTS.md: %w", err)
	}
	fmt.Println("  [ok] AGENTS.md  (Claude Code)")
	return nil
}

func generateCursor(targetDir string, selected []Agent) error {
	rulesDir := filepath.Join(targetDir, ".cursor", "rules")
	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		return err
	}
	for _, a := range selected {
		content := "---\ndescription: " + a.Trigger + "\nglobs: \nalwaysApply: false\n---\n\n" + a.Content()
		path := filepath.Join(rulesDir, a.Name+".mdc")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
		fmt.Printf("  [ok] .cursor/rules/%s.mdc\n", a.Name)
	}
	return nil
}

func generateCopilot(targetDir string, selected []Agent) error {
	dir := filepath.Join(targetDir, ".github")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(dir, "copilot-instructions.md")
	if err := os.WriteFile(path, []byte(combinedContent(selected)), 0644); err != nil {
		return err
	}
	fmt.Println("  [ok] .github/copilot-instructions.md  (GitHub Copilot)")
	return nil
}

func generateWindsurf(targetDir string, selected []Agent) error {
	path := filepath.Join(targetDir, ".windsurfrules")
	if err := os.WriteFile(path, []byte(combinedContent(selected)), 0644); err != nil {
		return err
	}
	fmt.Println("  [ok] .windsurfrules  (Windsurf)")
	return nil
}

func generateCline(targetDir string, selected []Agent) error {
	path := filepath.Join(targetDir, ".clinerules")
	if err := os.WriteFile(path, []byte(combinedContent(selected)), 0644); err != nil {
		return err
	}
	fmt.Println("  [ok] .clinerules  (Cline)")
	return nil
}

func generateAider(targetDir string, selected []Agent) error {
	path := filepath.Join(targetDir, "CONVENTIONS.md")
	if err := os.WriteFile(path, []byte(combinedContent(selected)), 0644); err != nil {
		return err
	}
	fmt.Println("  [ok] CONVENTIONS.md  (Aider)")
	return nil
}

func generateContinue(targetDir string, selected []Agent) error {
	dir := filepath.Join(targetDir, ".continue")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(dir, "rules.md")
	if err := os.WriteFile(path, []byte(combinedContent(selected)), 0644); err != nil {
		return err
	}
	fmt.Println("  [ok] .continue/rules.md  (Continue)")
	return nil
}

func generateOpenCode(targetDir string, selected []Agent) error {
	path := filepath.Join(targetDir, "AGENTS.md")
	if _, err := os.Stat(path); err == nil {
		fmt.Println("  [ok] AGENTS.md  (OpenCode — already written by Claude Code)")
		return nil
	}
	if err := os.WriteFile(path, []byte(agentsIndex(selected)), 0644); err != nil {
		return fmt.Errorf("failed to write AGENTS.md: %w", err)
	}
	fmt.Println("  [ok] AGENTS.md  (OpenCode)")
	return nil
}

// ─── Content builders ─────────────────────────────────────────────────────────

// agentsIndex generates the AGENTS.md index file used by Claude Code and OpenCode.
func agentsIndex(selected []Agent) string {
	var sb strings.Builder
	sb.WriteString("# Agent Skills Index\n\n")
	sb.WriteString("Load the relevant skill **before** writing any code for that layer.\n\n")
	sb.WriteString("## How to Use\n\n")
	sb.WriteString("1. Find the skill whose trigger matches your current task.\n")
	sb.WriteString("2. Read the `SKILL.md` at the listed path.\n")
	sb.WriteString("3. Follow every rule in the loaded skill without exception.\n")
	sb.WriteString("4. Multiple skills can be active simultaneously.\n\n")
	sb.WriteString("## Installed Skills\n\n")
	sb.WriteString("| Skill | Trigger | Path |\n")
	sb.WriteString("|-------|---------|------|\n")
	for _, a := range selected {
		path := fmt.Sprintf("skills/%s/SKILL.md", a.Dir)
		sb.WriteString(fmt.Sprintf("| `mr-agent-%s` | %s | [`%s`](%s) |\n",
			a.Name, a.Trigger, path, path))
	}
	return sb.String()
}

// combinedContent merges all selected skills into a single markdown document.
// Used by tools that read a single config file (Copilot, Windsurf, Cline, Aider, Continue).
func combinedContent(selected []Agent) string {
	var sb strings.Builder
	sb.WriteString("# Agent Skills\n\n")
	sb.WriteString("> Generated by mr-agent-ai. Read the relevant section before writing any code for that layer.\n\n")
	sb.WriteString("## How to Use\n\n")
	sb.WriteString("1. Find the skill whose trigger matches your current task.\n")
	sb.WriteString("2. Follow every rule in that section without exception.\n")
	sb.WriteString("3. Multiple skills can be active simultaneously.\n\n")
	sb.WriteString("---\n\n")
	for _, a := range selected {
		sb.WriteString("## " + a.Name + "\n\n")
		sb.WriteString("**Trigger:** " + a.Trigger + "\n\n")
		sb.WriteString(a.Content())
		sb.WriteString("\n\n---\n\n")
	}
	return sb.String()
}

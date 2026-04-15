package main

import (
	"fmt"
	"os"

	"mr-agent-ai/internal/audit"
	"mr-agent-ai/internal/skills"
	"mr-agent-ai/internal/theme"

	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:     "mr-agent-ai",
		Short:   "Multi-Agent AI Skill Installer",
		Long:    `Installs agent skills (SKILL.md files) and AGENTS.md into the current project.`,
		Version: version,
	}

	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install agent skills into the current project",
		Long: `Creates skills/<agent>/SKILL.md for each agent and AGENTS.md in the current directory.

Run this command from the root of your project:

  cd my-project
  mr-agent-ai install`,
		RunE: func(cmd *cobra.Command, args []string) error {
			targetDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("cannot determine current directory: %w", err)
			}
			return skills.Install(targetDir)
		},
	}

	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Scan the current project for supply chain security issues",
		Long: `Scans the project for missing security configurations and offers to apply fixes.

Run this command from the root of your project:

  cd my-project
  mr-agent-ai audit`,
		RunE: func(cmd *cobra.Command, args []string) error {
			targetDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("cannot determine current directory: %w", err)
			}
			return audit.RunAudit(targetDir)
		},
	}

	themeCmd := &cobra.Command{
		Use:   "theme",
		Short: "Context widget panel for AI agents",
		Long: `Shows a context panel with git branch, open PR, Linear ticket,
token usage, and rate limit info.

Run theme init to configure your AI agent integrations:

  cd my-project
  mr-agent-ai theme init`,
	}

	themeInitCmd := &cobra.Command{
		Use:   "init",
		Short: "Detect installed AI agents and configure context widgets",
		Long: `Detects which AI agents are installed (Claude Code, Aider, VS Code)
and configures the context widget panel for each selected agent.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			targetDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("cannot determine current directory: %w", err)
			}
			return theme.InitConfig(targetDir)
		},
	}

	themeShowCmd := &cobra.Command{
		Use:   "show",
		Short: "Render the context panel",
		Long: `Renders the context panel with all active widgets.

Use --format=statusline for Claude Code's status line integration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			targetDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("cannot determine current directory: %w", err)
			}
			format, _ := cmd.Flags().GetString("format")
			return theme.Show(targetDir, format == "statusline")
		},
	}
	themeShowCmd.Flags().String("format", "panel", "Output format: panel | statusline")

	themeCmd.AddCommand(themeInitCmd, themeShowCmd)

	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(auditCmd)
	rootCmd.AddCommand(themeCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

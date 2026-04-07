package main

import (
	"fmt"
	"os"

	"mr-agent-ai/internal/installer"

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
			return installer.Install(targetDir)
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
			return installer.RunAudit(targetDir)
		},
	}

	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(auditCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

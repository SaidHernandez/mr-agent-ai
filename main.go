package main

import (
	"fmt"
	"os"

	"mr-agent-ai/installer/internal/installer"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "mr-agent-ai",
		Short: "Multi-Agent AI Skill Installer",
		Long:  `Installs agent skills (SKILL.md files) and AGENTS.md into the current project.`,
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

	rootCmd.AddCommand(installCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

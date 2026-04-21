package theme

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// configureVSCode creates or updates .vscode/tasks.json in the project directory
// with a task that runs mr-agent-ai theme show on folder open.
func configureVSCode(dir string) error {
	vscodeDir := filepath.Join(dir, ".vscode")
	if err := os.MkdirAll(vscodeDir, 0750); err != nil {
		return fmt.Errorf("creating .vscode: %w", err)
	}

	tasksPath := filepath.Join(vscodeDir, "tasks.json")

	// Read existing tasks.json or create a fresh structure
	tasksDoc := map[string]any{
		"version": "2.0.0",
		"tasks":   []any{},
	}

	if data, err := os.ReadFile(tasksPath); err == nil { // #nosec G304
		if err := json.Unmarshal(data, &tasksDoc); err != nil {
			return fmt.Errorf("parsing .vscode/tasks.json: %w", err)
		}
	}

	// Check if the mr-agent context task already exists
	tasks, _ := tasksDoc["tasks"].([]any)
	for _, t := range tasks {
		if task, ok := t.(map[string]any); ok {
			if task["label"] == "mr-agent context" {
				return nil // already configured
			}
		}
	}

	// Append the new task
	newTask := map[string]any{
		"label":   "mr-agent context",
		"type":    "shell",
		"command": "mr-agent-ai theme show",
		"presentation": map[string]any{
			"reveal": "always",
			"panel":  "dedicated",
		},
		"runOptions": map[string]any{
			"runOn": "folderOpen",
		},
	}
	tasksDoc["tasks"] = append(tasks, newTask)

	data, err := json.MarshalIndent(tasksDoc, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling tasks.json: %w", err)
	}

	if err := os.WriteFile(tasksPath, append(data, '\n'), 0600); err != nil {
		return fmt.Errorf("writing .vscode/tasks.json: %w", err)
	}

	return nil
}

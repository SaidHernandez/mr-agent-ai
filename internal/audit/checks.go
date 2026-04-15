package audit

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func detectEcosystems(dir string) []string {
	type entry struct{ file, name string }
	candidates := []entry{
		{"package.json", "Node.js"},
		{"pyproject.toml", "Python"},
		{"requirements.txt", "Python"},
		{"go.mod", "Go"},
		{"pom.xml", "Java"},
		{"build.gradle", "Java"},
	}
	seen := map[string]bool{}
	var result []string
	for _, c := range candidates {
		if fileExists(dir, c.file) && !seen[c.name] {
			result = append(result, c.name)
			seen[c.name] = true
		}
	}
	return result
}

func hasEco(ecosystems []string, name string) bool {
	for _, e := range ecosystems {
		if e == name {
			return true
		}
	}
	return false
}

func missingGitignorePatterns(dir string) []string {
	required := []string{".env", "*.key", "*.pem", "*.p12", ".agentrc.yml", ".vscode/tasks.json"}
	b, err := os.ReadFile(filepath.Join(dir, ".gitignore")) // #nosec G304 -- path constructed from os.Getwd(), not external input
	if err != nil {
		return required
	}
	content := string(b)
	var missing []string
	for _, p := range required {
		if !strings.Contains(content, p) {
			missing = append(missing, p)
		}
	}
	return missing
}

func isEnvTracked(dir string) bool {
	cmd := exec.Command("git", "ls-files", "--error-unmatch", ".env")
	cmd.Dir = dir
	return cmd.Run() == nil
}

func fileExists(dir, rel string) bool {
	_, err := os.Stat(filepath.Join(dir, rel))
	return err == nil
}

func globExists(dir, pattern string) bool {
	matches, _ := filepath.Glob(filepath.Join(dir, pattern))
	return len(matches) > 0
}

func fileContains(dir, rel, substr string) bool {
	b, err := os.ReadFile(filepath.Join(dir, rel)) // #nosec G304 -- path constructed from os.Getwd(), not external input
	if err != nil {
		return false
	}
	return strings.Contains(string(b), substr)
}

func anyWorkflowContains(dir, substr string) bool {
	pattern := filepath.Join(dir, ".github", "workflows", "*.yml")
	files, _ := filepath.Glob(pattern)
	for _, f := range files {
		b, err := os.ReadFile(f) // #nosec G304 -- path comes from filepath.Glob on a trusted directory
		if err != nil {
			continue
		}
		if strings.Contains(string(b), substr) {
			return true
		}
	}
	return false
}

func workflowsMissingPermissions(dir string) []string {
	pattern := filepath.Join(dir, ".github", "workflows", "*.yml")
	files, _ := filepath.Glob(pattern)
	var missing []string
	for _, f := range files {
		b, err := os.ReadFile(f) // #nosec G304 -- path comes from filepath.Glob on a trusted directory
		if err != nil {
			continue
		}
		if !strings.Contains(string(b), "permissions:") {
			missing = append(missing, filepath.Base(f))
		}
	}
	return missing
}

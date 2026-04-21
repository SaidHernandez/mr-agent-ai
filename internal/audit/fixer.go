package audit

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type auditFix struct {
	display   string
	destPath  string
	template  string
	append    bool
	generator func(dir string) string
}

func applyFix(dir string, f auditFix) error {
	var content string
	if f.generator != nil {
		content = f.generator(dir)
	} else {
		content = mustReadContent(f.template)
	}
	dest := filepath.Join(dir, f.destPath)

	if err := os.MkdirAll(filepath.Dir(dest), 0o750); err != nil {
		return err
	}

	if f.append {
		file, err := os.OpenFile(dest, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644) // #nosec G302 G304 -- config files for user projects are intentionally world-readable
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = file.WriteString(content)
		return err
	}

	if err := os.WriteFile(dest, []byte(content), 0o644); err != nil { // #nosec G306 -- config files for user projects are intentionally world-readable
		return err
	}
	fmt.Printf("  ✓  %s\n", f.destPath)
	return nil
}

func dedupFixes(fixes []auditFix) []auditFix {
	seen := map[string]bool{}
	var result []auditFix
	for _, f := range fixes {
		if !seen[f.destPath] {
			result = append(result, f)
			seen[f.destPath] = true
		}
	}
	return result
}

func parseAuditSelection(input string, count int) []int {
	if strings.ToLower(input) == "all" {
		result := make([]int, count)
		for i := range result {
			result[i] = i
		}
		return result
	}
	var selected []int
	for _, part := range strings.Split(input, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || n < 1 || n > count {
			continue
		}
		selected = append(selected, n-1)
	}
	return selected
}

func generateDependabotYML(dir string) string {
	ecosystems := detectEcosystems(dir)

	type ecoEntry struct {
		pkgEco    string
		groupName string
	}
	ecoMap := map[string]ecoEntry{
		"Node.js": {"npm", "npm-deps"},
		"Python":  {"pip", "python-deps"},
		"Go":      {"gomod", "go-deps"},
		"Java":    {"maven", "java-deps"},
	}

	var sb strings.Builder
	sb.WriteString("version: 2\nupdates:\n")

	for _, eco := range ecosystems {
		entry, ok := ecoMap[eco]
		if !ok {
			continue
		}
		sb.WriteString(fmt.Sprintf(
			"  - package-ecosystem: %s\n    directory: \"/\"\n    schedule:\n      interval: weekly\n    labels:\n      - dependencies\n    groups:\n      %s:\n        patterns:\n          - \"*\"\n\n",
			entry.pkgEco, entry.groupName,
		))
	}

	sb.WriteString("  - package-ecosystem: github-actions\n    directory: \"/\"\n    schedule:\n      interval: weekly\n    labels:\n      - dependencies\n      - ci\n    groups:\n      actions:\n        patterns:\n          - \"*\"\n")

	return sb.String()
}

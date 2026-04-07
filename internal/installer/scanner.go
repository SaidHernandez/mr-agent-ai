package installer

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// RunAudit is the entry point for the audit command.
func RunAudit(dir string) error {
	ecosystems := detectEcosystems(dir)

	fmt.Println()
	fmt.Println("  Supply Chain Audit")
	fmt.Println("  " + strings.Repeat("─", 38))
	fmt.Println()

	if len(ecosystems) == 0 {
		fmt.Println("  Ecosystems   (none detected)")
	} else {
		fmt.Printf("  Ecosystems   %s\n", strings.Join(ecosystems, " · "))
	}
	fmt.Println()

	var fixes []auditFix
	issues := 0
	warnings := 0

	hasNode := hasEco(ecosystems, "Node.js")
	hasPython := hasEco(ecosystems, "Python")
	hasGo := hasEco(ecosystems, "Go")
	hasDocker := fileExists(dir, "Dockerfile")
	hasIaC := fileExists(dir, "docker-compose.yml") || globExists(dir, "*.tf") || fileExists(dir, "k8s")

	// ── Supply Chain ──────────────────────────────────────────────────────────
	printSection("Supply Chain")

	// Dependabot
	if ok := fileExists(dir, ".github/dependabot.yml"); ok {
		printCheck("Dependabot", true, false, "")
	} else {
		printCheck("Dependabot", false, false, "not configured")
		issues++
		fixes = append(fixes, auditFix{
			display:  ".github/dependabot.yml",
			destPath: ".github/dependabot.yml",
			template: "audit/dependabot.yml",
		})
	}

	// Shelf time — npm
	if hasNode {
		if ok := fileContains(dir, ".npmrc", "min-release-age"); ok {
			printCheck("npm shelf-time", true, false, "")
		} else {
			printCheck("npm shelf-time", false, false, ".npmrc → min-release-age=7d")
			issues++
			fixes = append(fixes, auditFix{
				display:  ".npmrc",
				destPath: ".npmrc",
				template: "audit/npmrc-shelf-time.txt",
				append:   true,
			})
		}
	}

	// Shelf time — Python
	if hasPython {
		if ok := fileContains(dir, "pyproject.toml", "exclude-newer"); ok {
			printCheck("Python shelf-time", true, false, "")
		} else {
			printCheck("Python shelf-time", false, false, "pyproject.toml → [tool.uv] exclude-newer = \"1 week\"")
			issues++
		}
	}

	// Lock files
	if hasNode {
		hasLock := fileExists(dir, "package-lock.json") || fileExists(dir, "yarn.lock") || fileExists(dir, "pnpm-lock.yaml")
		printCheck("npm lock file", hasLock, false, "")
		if !hasLock {
			issues++
		}
	}

	if hasPython {
		hasLock := fileExists(dir, "uv.lock") || fileExists(dir, "poetry.lock")
		printCheck("Python lock file", hasLock, false, "")
		if !hasLock {
			issues++
		}
	}

	if hasGo {
		if ok := fileExists(dir, "go.sum"); ok {
			printCheck("go.sum present", true, false, "")
		} else {
			printCheck("go.sum present", false, false, "run: go mod tidy")
			issues++
		}
	}

	// ── Secrets ───────────────────────────────────────────────────────────────
	fmt.Println()
	printSection("Secrets")

	// .gitignore patterns
	missingPatterns := missingGitignorePatterns(dir)
	if len(missingPatterns) == 0 {
		printCheck(".gitignore", true, false, "")
	} else {
		printCheck(".gitignore", false, false, "missing: "+strings.Join(missingPatterns, " · "))
		issues++
		fixes = append(fixes, auditFix{
			display:  ".gitignore",
			destPath: ".gitignore",
			template: "audit/gitignore-secrets.txt",
			append:   true,
		})
	}

	// .env tracked in git
	if envTracked := isEnvTracked(dir); envTracked {
		printCheck(".env not tracked", false, false, ".env is committed — remove it from git")
		issues++
	} else {
		printCheck(".env not tracked", true, false, "")
	}

	// Pre-commit hooks
	hasPreCommit := fileExists(dir, ".pre-commit-config.yaml") || fileExists(dir, ".husky")
	if hasPreCommit {
		printCheck("pre-commit hooks", true, false, "")
	} else {
		printCheck("pre-commit hooks", false, false, "not configured")
		issues++
		fixes = append(fixes, auditFix{
			display:  ".pre-commit-config.yaml",
			destPath: ".pre-commit-config.yaml",
			template: "audit/pre-commit-config.yaml",
		})
	}

	// Secret scanner
	hasGitleaks := fileExists(dir, ".gitleaks.toml") || anyWorkflowContains(dir, "gitleaks")
	hasTrufflehog := anyWorkflowContains(dir, "trufflehog")
	if hasGitleaks || hasTrufflehog {
		printCheck("secret scanner", true, false, "")
	} else {
		printCheck("secret scanner", false, false, "gitleaks or trufflehog not found")
		issues++
		fixes = append(fixes, auditFix{
			display:  ".gitleaks.toml",
			destPath: ".gitleaks.toml",
			template: "audit/gitleaks.toml",
		})
	}

	// ── SAST ──────────────────────────────────────────────────────────────────
	fmt.Println()
	printSection("SAST")

	// Security linter per ecosystem
	if hasGo {
		if ok := anyWorkflowContains(dir, "gosec"); ok {
			printCheck("gosec", true, false, "")
		} else {
			printCheck("gosec", false, false, "not in CI")
			issues++
			fixes = append(fixes, auditFix{
				display:  ".github/workflows/sast.yml",
				destPath: ".github/workflows/sast.yml",
				template: "audit/sast-go.yml",
			})
		}
	}

	if hasPython {
		if ok := anyWorkflowContains(dir, "bandit"); ok {
			printCheck("bandit", true, false, "")
		} else {
			printCheck("bandit", false, false, "not in CI")
			issues++
			if !hasGo { // avoid duplicate fix entry
				fixes = append(fixes, auditFix{
					display:  ".github/workflows/sast.yml",
					destPath: ".github/workflows/sast.yml",
					template: "audit/sast-python.yml",
				})
			}
		}
	}

	if hasNode {
		hasSast := anyWorkflowContains(dir, "eslint-plugin-security") ||
			anyWorkflowContains(dir, "eslint") && fileContains(dir, "package.json", "eslint-plugin-security")
		if hasSast {
			printCheck("eslint-security", true, false, "")
		} else {
			printCheck("eslint-security", false, false, "not in CI")
			issues++
			if !hasGo && !hasPython {
				fixes = append(fixes, auditFix{
					display:  ".github/workflows/sast.yml",
					destPath: ".github/workflows/sast.yml",
					template: "audit/sast-node.yml",
				})
			}
		}
	}

	// CodeQL
	if ok := anyWorkflowContains(dir, "codeql"); ok {
		printCheck("CodeQL", true, false, "")
	} else {
		printCheck("CodeQL", false, false, "not configured")
		issues++
		fixes = append(fixes, auditFix{
			display:  ".github/workflows/codeql.yml",
			destPath: ".github/workflows/codeql.yml",
			template: "audit/codeql.yml",
		})
	}

	// SonarCloud
	hasSonar := fileExists(dir, "sonar-project.properties") || anyWorkflowContains(dir, "sonar")
	if hasSonar {
		printCheck("SonarCloud", true, false, "")
	} else {
		printCheck("SonarCloud", false, false, "not configured")
		issues++
	}

	// ── Containers ────────────────────────────────────────────────────────────
	if hasDocker {
		fmt.Println()
		printSection("Containers")

		if ok := !fileContains(dir, "Dockerfile", ":latest"); ok {
			printCheck("no :latest tag", true, false, "")
		} else {
			printCheck("no :latest tag", false, false, "pin base image to a specific digest")
			issues++
		}

		if ok := fileContains(dir, "Dockerfile", "USER "); ok {
			printCheck("non-root USER", true, false, "")
		} else {
			printCheck("non-root USER", false, false, "add USER nonroot before ENTRYPOINT")
			issues++
		}

		if ok := anyWorkflowContains(dir, "trivy"); ok {
			printCheck("Trivy in CI", true, false, "")
		} else {
			printCheck("Trivy in CI", false, false, "not configured")
			issues++
			fixes = append(fixes, auditFix{
				display:  ".github/workflows/trivy.yml",
				destPath: ".github/workflows/trivy.yml",
				template: "audit/trivy.yml",
			})
		}
	}

	// ── IaC ───────────────────────────────────────────────────────────────────
	if hasIaC {
		fmt.Println()
		printSection("IaC")

		hasCheckov := anyWorkflowContains(dir, "checkov") || anyWorkflowContains(dir, "terrascan")
		if hasCheckov {
			printCheck("Checkov / Terrascan", true, false, "")
		} else {
			printCheck("Checkov / Terrascan", false, false, "not in CI")
			issues++
			fixes = append(fixes, auditFix{
				display:  ".github/workflows/iac-scan.yml",
				destPath: ".github/workflows/iac-scan.yml",
				template: "audit/iac-scan.yml",
			})
		}
	}

	// ── GitHub Actions ────────────────────────────────────────────────────────
	missingPerms := workflowsMissingPermissions(dir)
	if len(missingPerms) > 0 {
		fmt.Println()
		printSection("GitHub Actions")
		printCheck("permissions", false, true, "missing in: "+strings.Join(missingPerms, ", "))
		warnings++
	}

	// ── Summary ───────────────────────────────────────────────────────────────
	fmt.Println()
	fmt.Println("  " + strings.Repeat("─", 38))

	if issues == 0 && warnings == 0 {
		fmt.Println("  All checks passed.")
		fmt.Println()
		return nil
	}

	parts := []string{}
	if issues > 0 {
		parts = append(parts, fmt.Sprintf("%d issue", issues))
		if issues > 1 {
			parts[len(parts)-1] += "s"
		}
	}
	if warnings > 0 {
		parts = append(parts, fmt.Sprintf("%d warning", warnings))
	}
	fmt.Printf("  %s\n", strings.Join(parts, "  ·  "))

	// ── Available fixes ───────────────────────────────────────────────────────
	fixes = dedupFixes(fixes)
	if len(fixes) == 0 {
		fmt.Println()
		return nil
	}

	fmt.Println()
	fmt.Println("  Available fixes")
	fmt.Println("  " + strings.Repeat("─", 15))
	for i, f := range fixes {
		fmt.Printf("  %d  %s\n", i+1, f.display)
	}

	fmt.Println()
	fmt.Print("  Apply (1,3 or all) › ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return nil
	}

	selected := parseAuditSelection(strings.TrimSpace(scanner.Text()), len(fixes))
	if len(selected) == 0 {
		return nil
	}

	fmt.Println()
	for _, idx := range selected {
		if err := applyFix(dir, fixes[idx]); err != nil {
			fmt.Printf("  ✗  %s  %s\n", fixes[idx].display, err)
		}
	}
	fmt.Println()

	return nil
}

// ── Types ─────────────────────────────────────────────────────────────────────

type auditFix struct {
	display  string
	destPath string
	template string
	append   bool
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func printSection(name string) {
	fmt.Printf("  %s\n", name)
}

func printCheck(label string, ok bool, warn bool, hint string) {
	sym := "✗"
	if ok {
		sym = "✓"
	} else if warn {
		sym = "⚠"
	}
	if hint != "" && !ok {
		fmt.Printf("  %s  %-24s %s\n", sym, label, hint)
	} else {
		fmt.Printf("  %s  %s\n", sym, label)
	}
}

func hasEco(ecosystems []string, name string) bool {
	for _, e := range ecosystems {
		if e == name {
			return true
		}
	}
	return false
}

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

func missingGitignorePatterns(dir string) []string {
	required := []string{".env", "*.key", "*.pem", "*.p12"}
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

func applyFix(dir string, f auditFix) error {
	content := mustReadDoc(f.template)
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
		if err != nil {
			return err
		}
	} else {
		if err := os.WriteFile(dest, []byte(content), 0o644); err != nil { // #nosec G306 -- config files for user projects are intentionally world-readable
			return err
		}
	}

	fmt.Printf("  ✓  %s\n", f.destPath)
	return nil
}

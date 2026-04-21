package skills

import "time"

// Asset is an optional file placed under the skill's assets/ subdirectory.
type Asset struct {
	Path    string
	Content string
}

// Skill defines a single installable skill.
type Skill struct {
	Name        string
	Dir         string // e.g. "arch/orchestrator" or "lang/typescript"
	Description string
	Trigger     string
	Content     func() string
	Assets      []Asset
}

func skillFrontmatter(name, description, trigger, allowedTools string) string {
	return "---\n" +
		"name: mr-agent-" + name + "\n" +
		"description: >\n" +
		"  " + description + "\n" +
		"  Trigger: " + trigger + "\n" +
		"license: MIT\n" +
		"allowed-tools: " + allowedTools + "\n" +
		"metadata:\n" +
		"  version: \"1.0\"\n" +
		"  generated: " + time.Now().Format("2006-01-02") + "\n" +
		"---\n\n"
}

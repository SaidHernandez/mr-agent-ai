package installer

import (
	"strings"
	"testing"
	"time"
)

// ─── skillFrontmatter ─────────────────────────────────────────────────────────

func TestSkillFrontmatter_ContainsRequiredFields(t *testing.T) {
	result := skillFrontmatter("my-skill", "A description", "When triggered.", "Read, Write")

	checks := map[string]string{
		"name":          "name: mr-agent-my-skill",
		"description":   "description:",
		"trigger":       "Trigger: When triggered.",
		"license":       "license: MIT",
		"allowed-tools": "allowed-tools: Read, Write",
		"version":       `version: "1.0"`,
		"generated":     "generated: " + time.Now().Format("2006-01-02"),
	}

	for field, expected := range checks {
		if !strings.Contains(result, expected) {
			t.Errorf("skillFrontmatter: missing %q, expected to find %q", field, expected)
		}
	}
}

func TestSkillFrontmatter_StartsAndEndsWithDelimiter(t *testing.T) {
	result := skillFrontmatter("s", "d", "t", "tools")
	if !strings.HasPrefix(result, "---\n") {
		t.Error("skillFrontmatter: should start with ---")
	}
	if !strings.Contains(result, "\n---\n") {
		t.Error("skillFrontmatter: should contain closing ---")
	}
}

// ─── archSkills catalog ───────────────────────────────────────────────────────

func TestArchSkills_NoDuplicateNames(t *testing.T) {
	seen := map[string]bool{}
	for _, s := range archSkills {
		if seen[s.Name] {
			t.Errorf("archSkills: duplicate skill name %q", s.Name)
		}
		seen[s.Name] = true
	}
}

func TestArchSkills_NoDuplicateDirs(t *testing.T) {
	seen := map[string]bool{}
	for _, s := range archSkills {
		if seen[s.Dir] {
			t.Errorf("archSkills: duplicate skill dir %q", s.Dir)
		}
		seen[s.Dir] = true
	}
}

func TestArchSkills_AllFieldsPopulated(t *testing.T) {
	for _, s := range archSkills {
		if s.Name == "" {
			t.Error("archSkills: skill with empty Name")
		}
		if s.Dir == "" {
			t.Errorf("archSkills[%s]: empty Dir", s.Name)
		}
		if s.Description == "" {
			t.Errorf("archSkills[%s]: empty Description", s.Name)
		}
		if s.Trigger == "" {
			t.Errorf("archSkills[%s]: empty Trigger", s.Name)
		}
		if s.Content == nil {
			t.Errorf("archSkills[%s]: nil Content func", s.Name)
		}
	}
}

func TestArchSkills_ContentNotEmpty(t *testing.T) {
	for _, s := range archSkills {
		content := s.Content()
		if strings.TrimSpace(content) == "" {
			t.Errorf("archSkills[%s]: Content() returned empty string", s.Name)
		}
		if !strings.Contains(content, "---") {
			t.Errorf("archSkills[%s]: Content() missing frontmatter delimiter", s.Name)
		}
	}
}

// ─── langSkills catalog ───────────────────────────────────────────────────────

func TestLangSkills_NoDuplicateNames(t *testing.T) {
	seen := map[string]bool{}
	for _, s := range langSkills {
		if seen[s.Name] {
			t.Errorf("langSkills: duplicate skill name %q", s.Name)
		}
		seen[s.Name] = true
	}
}

func TestLangSkills_AllFieldsPopulated(t *testing.T) {
	for _, s := range langSkills {
		if s.Name == "" {
			t.Error("langSkills: skill with empty Name")
		}
		if s.Dir == "" {
			t.Errorf("langSkills[%s]: empty Dir", s.Name)
		}
		if s.Content == nil {
			t.Errorf("langSkills[%s]: nil Content func", s.Name)
		}
	}
}

func TestLangSkills_ContentNotEmpty(t *testing.T) {
	for _, s := range langSkills {
		content := s.Content()
		if strings.TrimSpace(content) == "" {
			t.Errorf("langSkills[%s]: Content() returned empty string", s.Name)
		}
	}
}

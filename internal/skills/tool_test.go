package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ─── Helpers ──────────────────────────────────────────────────────────────────

func stubSkills() []Skill {
	return []Skill{
		{
			Name:        "test-skill",
			Dir:         "arch/test-skill",
			Description: "A test skill",
			Trigger:     "When testing.",
			Content:     func() string { return "# Test Skill\nsome content\n" },
		},
		{
			Name:        "another-skill",
			Dir:         "arch/another-skill",
			Description: "Another test skill",
			Trigger:     "When testing again.",
			Content:     func() string { return "# Another Skill\nmore content\n" },
		},
	}
}

func stubSkillWithAsset() Skill {
	return Skill{
		Name:        "skill-with-asset",
		Dir:         "arch/skill-with-asset",
		Description: "Skill with asset",
		Trigger:     "When testing assets.",
		Content:     func() string { return "# Skill With Asset\n" },
		Assets: []Asset{
			{Path: "template.sql", Content: "SELECT 1;"},
		},
	}
}

// ─── agentsIndex ──────────────────────────────────────────────────────────────

func TestAgentsIndex_ContainsAllSkills(t *testing.T) {
	skills := stubSkills()
	result := agentsIndex(skills)

	for _, s := range skills {
		if !strings.Contains(result, s.Name) {
			t.Errorf("agentsIndex: missing skill name %q", s.Name)
		}
		if !strings.Contains(result, s.Trigger) {
			t.Errorf("agentsIndex: missing trigger for %q", s.Name)
		}
		expectedPath := "skills/" + s.Dir + "/SKILL.md"
		if !strings.Contains(result, expectedPath) {
			t.Errorf("agentsIndex: missing path %q", expectedPath)
		}
	}
}

func TestAgentsIndex_HasHeader(t *testing.T) {
	result := agentsIndex(stubSkills())
	if !strings.HasPrefix(result, "# Agent Skills Index") {
		t.Error("agentsIndex: missing expected header")
	}
}

// ─── combinedContent ──────────────────────────────────────────────────────────

func TestCombinedContent_ContainsAllSkills(t *testing.T) {
	skills := stubSkills()
	result := combinedContent(skills)

	for _, s := range skills {
		if !strings.Contains(result, s.Name) {
			t.Errorf("combinedContent: missing skill name %q", s.Name)
		}
		if !strings.Contains(result, s.Content()) {
			t.Errorf("combinedContent: missing content for %q", s.Name)
		}
		if !strings.Contains(result, s.Trigger) {
			t.Errorf("combinedContent: missing trigger for %q", s.Name)
		}
	}
}

func TestCombinedContent_HasHeader(t *testing.T) {
	result := combinedContent(stubSkills())
	if !strings.HasPrefix(result, "# Agent Skills") {
		t.Error("combinedContent: missing expected header")
	}
}

// ─── writeSkills ──────────────────────────────────────────────────────────────

func TestWriteSkills_CreatesSkillFiles(t *testing.T) {
	dir := t.TempDir()
	skills := stubSkills()

	if err := writeSkills(dir, skills); err != nil {
		t.Fatalf("writeSkills: unexpected error: %v", err)
	}

	for _, s := range skills {
		path := filepath.Join(dir, "skills", s.Dir, "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("writeSkills: expected file %q not found: %v", path, err)
			continue
		}
		if string(data) != s.Content() {
			t.Errorf("writeSkills: content mismatch for %q", s.Name)
		}
	}
}

func TestWriteSkills_CreatesAssets(t *testing.T) {
	dir := t.TempDir()
	skill := stubSkillWithAsset()

	if err := writeSkills(dir, []Skill{skill}); err != nil {
		t.Fatalf("writeSkills: unexpected error: %v", err)
	}

	assetPath := filepath.Join(dir, "skills", skill.Dir, "assets", "template.sql")
	data, err := os.ReadFile(assetPath)
	if err != nil {
		t.Fatalf("writeSkills: asset not found at %q: %v", assetPath, err)
	}
	if string(data) != "SELECT 1;" {
		t.Errorf("writeSkills: asset content mismatch, got %q", string(data))
	}
}

// ─── Generators ───────────────────────────────────────────────────────────────

func TestGenerateClaudeCode_CreatesAgentsMd(t *testing.T) {
	dir := t.TempDir()
	skills := stubSkills()

	if err := generateClaudeCode(dir, skills); err != nil {
		t.Fatalf("generateClaudeCode: unexpected error: %v", err)
	}

	path := filepath.Join(dir, "AGENTS.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("generateClaudeCode: AGENTS.md not found: %v", err)
	}
	if !strings.Contains(string(data), "test-skill") {
		t.Error("generateClaudeCode: AGENTS.md does not contain skill name")
	}
}

func TestGenerateCursor_CreatesOneMdcPerSkill(t *testing.T) {
	dir := t.TempDir()
	skills := stubSkills()

	if err := generateCursor(dir, skills); err != nil {
		t.Fatalf("generateCursor: unexpected error: %v", err)
	}

	for _, s := range skills {
		path := filepath.Join(dir, ".cursor", "rules", s.Name+".mdc")
		if _, err := os.Stat(path); err != nil {
			t.Errorf("generateCursor: expected file %q not found", path)
		}
	}
}

func TestGenerateCursor_MdcContainsTrigger(t *testing.T) {
	dir := t.TempDir()
	skills := stubSkills()

	if err := generateCursor(dir, skills); err != nil {
		t.Fatalf("generateCursor: unexpected error: %v", err)
	}

	path := filepath.Join(dir, ".cursor", "rules", skills[0].Name+".mdc")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("generateCursor: file not found: %v", err)
	}
	if !strings.Contains(string(data), skills[0].Trigger) {
		t.Error("generateCursor: .mdc file missing trigger")
	}
}

func TestGenerateCopilot_CreatesCopilotInstructions(t *testing.T) {
	dir := t.TempDir()

	if err := generateCopilot(dir, stubSkills()); err != nil {
		t.Fatalf("generateCopilot: unexpected error: %v", err)
	}

	path := filepath.Join(dir, ".github", "copilot-instructions.md")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("generateCopilot: file not found: %v", err)
	}
}

func TestGenerateWindsurf_CreatesWindsurfrules(t *testing.T) {
	dir := t.TempDir()

	if err := generateWindsurf(dir, stubSkills()); err != nil {
		t.Fatalf("generateWindsurf: unexpected error: %v", err)
	}

	path := filepath.Join(dir, ".windsurfrules")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("generateWindsurf: .windsurfrules not found: %v", err)
	}
}

func TestGenerateCline_CreatesClinerules(t *testing.T) {
	dir := t.TempDir()

	if err := generateCline(dir, stubSkills()); err != nil {
		t.Fatalf("generateCline: unexpected error: %v", err)
	}

	path := filepath.Join(dir, ".clinerules")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("generateCline: .clinerules not found: %v", err)
	}
}

func TestGenerateAider_CreatesConventionsMd(t *testing.T) {
	dir := t.TempDir()

	if err := generateAider(dir, stubSkills()); err != nil {
		t.Fatalf("generateAider: unexpected error: %v", err)
	}

	path := filepath.Join(dir, "CONVENTIONS.md")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("generateAider: CONVENTIONS.md not found: %v", err)
	}
}

func TestGenerateContinue_CreatesContinueRules(t *testing.T) {
	dir := t.TempDir()

	if err := generateContinue(dir, stubSkills()); err != nil {
		t.Fatalf("generateContinue: unexpected error: %v", err)
	}

	path := filepath.Join(dir, ".continue", "rules.md")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("generateContinue: .continue/rules.md not found: %v", err)
	}
}

func TestGenerateOpenCode_CreatesAgentsMd(t *testing.T) {
	dir := t.TempDir()

	if err := generateOpenCode(dir, stubSkills()); err != nil {
		t.Fatalf("generateOpenCode: unexpected error: %v", err)
	}

	path := filepath.Join(dir, "AGENTS.md")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("generateOpenCode: AGENTS.md not found: %v", err)
	}
}

func TestGenerateOpenCode_SkipsIfClaudeCodeAlreadyWrote(t *testing.T) {
	dir := t.TempDir()
	agentsPath := filepath.Join(dir, "AGENTS.md")

	original := "# Original content"
	if err := os.WriteFile(agentsPath, []byte(original), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := generateOpenCode(dir, stubSkills()); err != nil {
		t.Fatalf("generateOpenCode: unexpected error: %v", err)
	}

	data, _ := os.ReadFile(agentsPath)
	if string(data) != original {
		t.Error("generateOpenCode: should not overwrite existing AGENTS.md")
	}
}

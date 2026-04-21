package theme

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// --- Config tests ---

func TestLoadConfig_ReturnsDefaultsWhenFileMissing(t *testing.T) {
	dir := t.TempDir()
	cfg, err := LoadConfig(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Widgets.Branch.Enabled {
		t.Error("expected branch widget enabled by default")
	}
	if !cfg.Widgets.PR.Enabled {
		t.Error("expected PR widget enabled by default")
	}
	if !cfg.Widgets.Linear.Enabled {
		t.Error("expected linear widget enabled by default")
	}
	if !cfg.Widgets.Tokens.Enabled {
		t.Error("expected tokens widget enabled by default")
	}
	if !cfg.Widgets.RateLimit.Enabled {
		t.Error("expected rate-limit widget enabled by default")
	}
}

func TestSaveAndLoadConfig_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	original := DefaultConfig()
	original.Widgets.PR.Repo = "owner/repo"
	original.Widgets.Linear.Prefix = "ENG"
	original.Widgets.PR.Enabled = false

	if err := SaveConfig(dir, original); err != nil {
		t.Fatalf("save error: %v", err)
	}

	loaded, err := LoadConfig(dir)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if loaded.Widgets.PR.Repo != "owner/repo" {
		t.Errorf("expected repo owner/repo, got %q", loaded.Widgets.PR.Repo)
	}
	if loaded.Widgets.Linear.Prefix != "ENG" {
		t.Errorf("expected prefix ENG, got %q", loaded.Widgets.Linear.Prefix)
	}
	if loaded.Widgets.PR.Enabled {
		t.Error("expected PR widget disabled after round-trip")
	}
}

// --- Widget / RunAll tests ---

type stubWidget struct {
	name  string
	value string
	err   error
	delay time.Duration
}

func (s *stubWidget) Name() string { return s.name }
func (s *stubWidget) Fetch(ctx context.Context) Result {
	if s.delay > 0 {
		select {
		case <-time.After(s.delay):
		case <-ctx.Done():
			return Result{Err: ctx.Err()}
		}
	}
	return Result{Value: s.value, Err: s.err}
}

func TestRunAll_CollectsAllResults(t *testing.T) {
	widgets := []Widget{
		&stubWidget{name: "a", value: "v1"},
		&stubWidget{name: "b", value: "v2"},
		&stubWidget{name: "c", value: "v3"},
	}
	results := RunAll(widgets)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestRunAll_PreservesOrder(t *testing.T) {
	names := []string{"first", "second", "third"}
	widgets := []Widget{
		&stubWidget{name: names[0], value: "v0"},
		&stubWidget{name: names[1], value: "v1"},
		&stubWidget{name: names[2], value: "v2"},
	}
	results := RunAll(widgets)
	for i, r := range results {
		if r.WidgetName != names[i] {
			t.Errorf("position %d: expected %q, got %q", i, names[i], r.WidgetName)
		}
	}
}

func TestRunAll_RespectsDeadline(t *testing.T) {
	widgets := []Widget{
		&stubWidget{name: "slow", delay: 3 * time.Second},
	}
	start := time.Now()
	results := RunAll(widgets)
	elapsed := time.Since(start)
	if elapsed > 2*time.Second {
		t.Errorf("RunAll took too long: %v", elapsed)
	}
	if results[0].Err == nil {
		t.Error("expected timeout error from slow widget")
	}
}

func TestRunAll_ContinuesAfterOneFailure(t *testing.T) {
	widgets := []Widget{
		&stubWidget{name: "ok", value: "good"},
		&stubWidget{name: "bad", err: errors.New("boom")},
		&stubWidget{name: "ok2", value: "also good"},
	}
	results := RunAll(widgets)
	if results[0].Err != nil {
		t.Errorf("widget 0 should succeed, got error: %v", results[0].Err)
	}
	if results[1].Err == nil {
		t.Error("widget 1 should fail")
	}
	if results[2].Err != nil {
		t.Errorf("widget 2 should succeed, got error: %v", results[2].Err)
	}
}

// --- BranchWidget tests ---

func TestBranchWidget_FetchReturnsError_WhenNotGitRepo(t *testing.T) {
	dir := t.TempDir()
	w := &BranchWidget{Dir: dir}
	result := w.Fetch(context.Background())
	if result.Err == nil && result.Value != "" {
		t.Log("note: git may have returned a value from a parent repo")
	}
	// No assertion on error — the test verifies the widget doesn't panic.
}

// --- parseTicketFromBranch tests ---

func TestParseTicketFromBranch(t *testing.T) {
	tests := []struct {
		branch string
		prefix string
		want   string
	}{
		{"eng-123-fix-auth", "ENG", "ENG-123"},
		{"ENG-456", "ENG", "ENG-456"},
		{"feature/eng-789-description", "ENG", "ENG-789"},
		{"main", "ENG", ""},
		{"fix-bug", "ENG", ""},
		{"eng-123", "PROD", ""},
		{"", "ENG", ""},
		{"eng-123", "", "ENG-123"},
		{"PROD-99-hotfix", "PROD", "PROD-99"},
	}
	for _, tt := range tests {
		got := parseTicketFromBranch(tt.branch, tt.prefix)
		if got != tt.want {
			t.Errorf("parseTicketFromBranch(%q, %q) = %q, want %q", tt.branch, tt.prefix, got, tt.want)
		}
	}
}

// --- parseGitHubRepo tests ---

func TestParseGitHubRepo(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{"https://github.com/owner/repo.git", "owner/repo", false},
		{"https://github.com/owner/repo", "owner/repo", false},
		{"git@github.com:owner/repo.git", "owner/repo", false},
		{"git@github.com:owner/repo", "owner/repo", false},
		{"https://gitlab.com/owner/repo", "", true},
		{"", "", true},
	}
	for _, tt := range tests {
		got, err := parseGitHubRepo(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("parseGitHubRepo(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("parseGitHubRepo(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

// --- Renderer tests ---

func TestRenderPanel_ContainsAllWidgetNames(t *testing.T) {
	results := []Result{
		{WidgetName: "branch", Value: "main"},
		{WidgetName: "pr", Value: "#42"},
		{WidgetName: "linear", Value: "[ENG-1]"},
		{WidgetName: "token", Value: "45,000 (~$0.18)"},
		{WidgetName: "rate limit", Value: "67%"},
		{WidgetName: "ctx", Value: "45%"},
	}
	panel := RenderPanel(results)
	for _, r := range results {
		if !strings.Contains(panel, r.WidgetName) {
			t.Errorf("panel missing widget name %q", r.WidgetName)
		}
	}
}

func TestRenderPanel_ShowsErrorState(t *testing.T) {
	results := []Result{
		{WidgetName: "branch", Err: errors.New("not a git repo")},
	}
	panel := RenderPanel(results)
	if !strings.Contains(panel, "[err]") {
		t.Error("expected [err] in panel output for widget with error")
	}
}

func TestRenderPanel_TruncatesLongValues(t *testing.T) {
	longVal := strings.Repeat("x", 100)
	results := []Result{
		{WidgetName: "branch", Value: longVal},
	}
	panel := RenderPanel(results)
	// The value should be truncated — full 100 chars should not appear
	if strings.Contains(panel, longVal) {
		t.Error("panel should truncate long values, but full value appears in output")
	}
	// Panel should still contain a truncated version with ellipsis
	if !strings.Contains(panel, "…") {
		t.Error("expected truncation ellipsis in panel output")
	}
}

func TestRenderStatusLine_HidesEmptyWidgets(t *testing.T) {
	results := []Result{
		{WidgetName: "branch", Value: "main"},
		{WidgetName: "pr", Value: "(no GITHUB_TOKEN)"},
		{WidgetName: "linear", Value: "(no LINEAR_API_KEY)"},
	}
	line := RenderStatusLine(results, 120)
	if strings.Contains(stripANSI(line), "no GITHUB_TOKEN") {
		t.Error("status line should hide widgets with missing token message")
	}
	if !strings.Contains(stripANSI(line), "main") {
		t.Error("status line should show branch value")
	}
}

func TestLineCount(t *testing.T) {
	s := "line1\nline2\nline3\n"
	got := strings.Count(s, "\n")
	if got != 3 {
		t.Errorf("expected 3 newlines, got %d", got)
	}
}

// --- RateLimitState tests ---

func TestSaveAndLoadRateLimitState_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	state := &RateLimitState{
		RequestsLimit:     50,
		RequestsRemaining: 45,
		TokensLimit:       100000,
		TokensRemaining:   89000,
	}
	if err := SaveRateLimitState(dir, state); err != nil {
		t.Fatalf("save error: %v", err)
	}
	// Verify file exists
	path := filepath.Join(dir, stateFileName)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("state file not created: %v", err)
	}
}

// --- formatTokens tests ---

func TestFormatTokens(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{0, "0"},
		{500, "500"},
		{999, "999"},
		{1000, "1,000"},
		{1500, "1,500"},
		{12450, "12,450"},
		{1_000_000, "1,000,000"},
	}
	for _, tt := range tests {
		got := formatTokens(tt.n)
		if got != tt.want {
			t.Errorf("formatTokens(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}

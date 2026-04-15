package theme

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var ghClient = &http.Client{Timeout: 1200 * time.Millisecond}

// PRWidget fetches the open PR for the current branch from GitHub.
type PRWidget struct {
	Dir    string
	Config PRWidgetConfig
}

func (w *PRWidget) Name() string { return "pr" }

func (w *PRWidget) Fetch(ctx context.Context) Result {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return Result{Value: "(no GITHUB_TOKEN)"}
	}

	repo := w.Config.Repo
	if repo == "" {
		var err error
		repo, err = detectRepo(ctx, w.Dir)
		if err != nil {
			return Result{Value: "(no repo)"}
		}
	}

	branch, err := currentBranch(ctx, w.Dir)
	if err != nil {
		return Result{Value: "(no branch)"}
	}

	// Build the owner prefix for head filter
	parts := strings.SplitN(repo, "/", 2)
	if len(parts) != 2 {
		return Result{Value: "(invalid repo)"}
	}
	owner := parts[0]

	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls?state=open&head=%s:%s", repo, owner, branch)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return Result{Err: err}
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := ghClient.Do(req)
	if err != nil {
		return Result{Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{Value: fmt.Sprintf("(GitHub %d)", resp.StatusCode)}
	}

	var prs []struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&prs); err != nil {
		return Result{Err: err}
	}

	if len(prs) == 0 {
		return Result{Value: "(no open PR)"}
	}

	title := prs[0].Title
	if len(title) > 28 {
		title = title[:27] + "…"
	}
	return Result{Value: fmt.Sprintf("#%d — %s", prs[0].Number, title)}
}

func detectRepo(ctx context.Context, dir string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "remote", "get-url", "origin")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return parseGitHubRepo(strings.TrimSpace(string(out)))
}

func parseGitHubRepo(rawURL string) (string, error) {
	rawURL = strings.TrimSuffix(rawURL, ".git")
	if strings.HasPrefix(rawURL, "https://github.com/") {
		return strings.TrimPrefix(rawURL, "https://github.com/"), nil
	}
	if strings.HasPrefix(rawURL, "git@github.com:") {
		return strings.TrimPrefix(rawURL, "git@github.com:"), nil
	}
	return "", fmt.Errorf("unrecognized GitHub remote URL: %s", rawURL)
}

func currentBranch(ctx context.Context, dir string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

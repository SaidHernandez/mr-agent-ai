package theme

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var linearClient = &http.Client{Timeout: 1200 * time.Millisecond}
var ticketRe = regexp.MustCompile(`(?i)([A-Za-z]{1,10}-\d+)`)

// LinearWidget fetches a Linear issue based on the ticket ID in the branch name.
type LinearWidget struct {
	Dir    string
	Config LinearWidgetConfig
}

func (w *LinearWidget) Name() string { return "linear" }

func (w *LinearWidget) Fetch(ctx context.Context) Result {
	apiKey := os.Getenv("LINEAR_API_KEY")
	if apiKey == "" {
		return Result{Value: "(no LINEAR_API_KEY)"}
	}

	branch, err := currentBranch(ctx, w.Dir)
	if err != nil {
		// currentBranch defined in pr.go (same package)
		cmd := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
		cmd.Dir = w.Dir
		out, e := cmd.Output()
		if e != nil {
			return Result{Value: "(no branch)"}
		}
		branch = strings.TrimSpace(string(out))
		_ = err
	}

	ticketID := parseTicketFromBranch(branch, w.Config.Prefix)
	if ticketID == "" {
		return Result{Value: "(no ticket)"}
	}

	query := fmt.Sprintf(`{"query": "{ issue(id: \"%s\") { title state { name } } }"}`, ticketID)
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.linear.app/graphql", bytes.NewBufferString(query))
	if err != nil {
		return Result{Err: err}
	}
	req.Header.Set("Authorization", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := linearClient.Do(req)
	if err != nil {
		return Result{Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{Value: fmt.Sprintf("[%s] (Linear %d)", ticketID, resp.StatusCode)}
	}

	var result struct {
		Data struct {
			Issue struct {
				Title string `json:"title"`
				State struct {
					Name string `json:"name"`
				} `json:"state"`
			} `json:"issue"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Result{Value: fmt.Sprintf("[%s]", ticketID)}
	}

	title := result.Data.Issue.Title
	if len(title) > 22 {
		title = title[:21] + "…"
	}
	state := result.Data.Issue.State.Name
	if title == "" {
		return Result{Value: fmt.Sprintf("[%s]", ticketID)}
	}
	return Result{Value: fmt.Sprintf("[%s] %s (%s)", ticketID, title, state)}
}

// parseTicketFromBranch extracts a Linear ticket ID from a branch name.
// If prefix is non-empty, only returns tickets matching that prefix.
func parseTicketFromBranch(branch, prefix string) string {
	matches := ticketRe.FindAllString(branch, -1)
	if len(matches) == 0 {
		return ""
	}
	if prefix == "" {
		return strings.ToUpper(matches[0])
	}
	pfxLower := strings.ToLower(prefix) + "-"
	for _, m := range matches {
		if strings.HasPrefix(strings.ToLower(m), pfxLower) {
			return strings.ToUpper(m)
		}
	}
	return ""
}

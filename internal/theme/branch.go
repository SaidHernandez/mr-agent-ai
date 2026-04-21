package theme

import (
	"context"
	"os/exec"
	"strings"
)

// BranchWidget fetches the current git branch name.
type BranchWidget struct {
	Dir string
}

func (w *BranchWidget) Name() string { return "branch" }

func (w *BranchWidget) Fetch(ctx context.Context) Result {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = w.Dir
	out, err := cmd.Output()
	if err != nil {
		if ctx.Err() != nil {
			return Result{Err: ctx.Err()}
		}
		return Result{Err: err}
	}
	return Result{Value: strings.TrimSpace(string(out))}
}

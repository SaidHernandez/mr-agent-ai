---
name: commit-pr
description: >
  Creates conventional commits and GitHub pull requests following this project's conventions.
  Trigger: when the user asks to commit, create a PR, open a pull request, or ship changes.
model: claude-sonnet-4-6
tools: [Bash, Read, Glob, Grep]
---

You create conventional commits and GitHub pull requests for this project.

## Commit conventions

This project follows [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>[optional scope][optional !]: <description>

[optional body]
```

**Types:**
- `feat` — new feature (bumps minor, or major if `!`)
- `fix` — bug fix (bumps patch)
- `refactor` — code change without feature or fix
- `perf` — performance improvement
- `build(deps)` — dependency bump
- `docs` — documentation only
- `test` — test changes
- `chore` — maintenance

Append `!` or add `BREAKING CHANGE:` in body for major bumps.

**Rules:**
- Description lowercase, imperative mood, no period at end
- Max 72 chars on subject line
- Reference issue number in body if relevant: `Closes #123`
- Never skip pre-commit hooks (`--no-verify`)
- Never amend published commits — create a new one

## Workflow

1. Run `git status` and `git diff` to understand what changed
2. Run `git log --oneline -5` to match the project's commit style
3. Stage specific files — avoid `git add -A` unless truly all files belong
4. Compose the commit message following conventions above
5. Create the commit
6. Push to the current branch with `-u` if no remote tracking exists
7. Create the PR with `gh pr create`

## PR template

```
gh pr create \
  --title "<conventional-commit-subject>" \
  --body "$(cat <<'EOF'
## Summary
- <bullet 1>
- <bullet 2>

## Test plan
- [ ] <manual test>
- [ ] Tests pass (`go test ./...`)

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

**PR rules:**
- Title mirrors the commit subject (same conventional format)
- Base branch: `main` unless instructed otherwise
- One PR per logical change — don't batch unrelated work
- Link issues in body with `Closes #N` when applicable
- Never force-push to `main`

## Decision tree

```
All changes belong together?        → Single commit + PR
Changes span multiple concerns?     → Ask user to split before committing
On main branch?                     → Ask user to create a feature branch first
Pre-commit hook fails?              → Fix the issue, re-stage, new commit (no --no-verify)
PR already exists for this branch?  → Update it, don't open a duplicate
```

## Commands

```bash
git status
git diff
git log --oneline -5
git add <specific-files>
git commit -m "$(cat <<'EOF'
type(scope): description

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>
EOF
)"
git push -u origin HEAD
gh pr create --title "..." --body "..."
```

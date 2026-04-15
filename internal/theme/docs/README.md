# theme

## Context

When working inside an AI agent session it's easy to lose track of where you are — which branch, whether there's an open PR, how much context window is left, or whether you're hitting rate limits. This package adds a real-time context panel that surfaces that information directly in the agent's UI. It integrates with Claude Code's status bar, Aider's script hook, and VS Code tasks so the panel stays visible without leaving the editor.

## When to use

### `theme init` — configure once per machine

```bash
cd my-project
mr-agent-ai theme init
```

Detects which agents are installed and wires the panel into each one:

```
  Detecting installed agents...

  ✓ Claude Code       found
  ✗ Aider             not found
  ✓ VS Code / Copilot found

   1.  Claude Code      — Anthropic CLI (~/.claude)
   2.  VS Code / Copilot — GitHub Copilot in VS Code (.vscode/tasks.json)

  Enter numbers (e.g. 1,3 or 'all'): all

  ✓ .agentrc.yml       created
  ✓ Claude Code        configured
  ✓ VS Code / Copilot  configured

  Done! Run mr-agent-ai theme show to preview your panel.
```

### `theme show` — preview the panel

```bash
mr-agent-ai theme show
```

```
┌─────────────────────────────────────────────┐
│  mr-agent-ai  context panel                 │
├─────────────────────────────────────────────┤
│  branch      feat/ENG-123-user-auth         │
│  pr          #42 — Add user registration…   │
│  linear      [ENG-123] User auth (In Pro…   │
│  tokens      ctx: 23%  $0.042               │
│  rate-limit  req: 45%                       │
├─────────────────────────────────────────────┤
│  updated 17:03:23                           │
└─────────────────────────────────────────────┘
```

Once configured, the panel renders automatically inside the agent — no need to run it manually.

## What it generates

| Agent | File modified | Behavior |
|-------|--------------|----------|
| Claude Code | `~/.claude/settings.json` | Adds `statusLine` command, refreshes every 30s |
| Aider | `~/.aider.conf.yml` + `~/.config/mr-agent-ai/aider-context.sh` | Runs panel script on session start |
| VS Code | `.vscode/tasks.json` | Adds task that runs on folder open |

### Local config

`.agentrc.yml` is created in the project root (already in `.gitignore`):

```yaml
widgets:
  pr:
    enabled: true
    repo: owner/repo      # optional, auto-detected from git remote
  linear:
    enabled: true
    prefix: ENG           # optional, filters by ticket prefix
```

### Environment variables

| Variable | Widget | How to get one |
|----------|--------|----------------|
| `GITHUB_TOKEN` | pr | [github.com/settings/tokens](https://github.com/settings/tokens) — scope: `repo` |
| `LINEAR_API_KEY` | linear | Linear → Settings → API → Personal API keys |

Add them to your shell profile (`~/.zshrc` or `~/.bashrc`):

```bash
export GITHUB_TOKEN="ghp_xxxxxxxxxxxxxxxxxxxx"
export LINEAR_API_KEY="lin_api_xxxxxxxxxxxxxxxxxxxx"
```

### Runtime cache files

| File | Content |
|------|---------|
| `.agentrc.input.json` | Last Claude Code stdin snapshot (rate limit, tokens, context) |
| `.agentrc.tokens.json` | Token usage cache (legacy fallback) |
| `.agentrc.state.json` | Rate limit state (written by external API callers via `SaveRateLimitState`) |

## Widget data sources

| Widget | Primary source | Fallback |
|--------|---------------|---------|
| branch | `git rev-parse` | — |
| pr | GitHub API (`GITHUB_TOKEN`) | `(no GITHUB_TOKEN)` |
| linear | Linear API (`LINEAR_API_KEY`) | `(no LINEAR_API_KEY)` |
| tokens | Claude Code stdin (JSON) | `.agentrc.input.json` |
| rate-limit | Claude Code stdin (JSON) | `.agentrc.input.json` → `.agentrc.state.json` |
| ctx | Claude Code stdin (JSON) | `.agentrc.input.json` |

## Adding a new widget

1. Create `<name>.go` implementing the `Widget` interface (`Name()` + `Fetch()`)
2. Add its config struct to `config.go` and register it in `buildWidgets()`

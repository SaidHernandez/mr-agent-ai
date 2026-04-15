```
  __  __ ____        _                    _         _    ___
 |  \/  |  _ \      / \   __ _  ___ _ __ | |_      / \  |_ _|
 | |\/| | |_) |    / _ \ / _`\|/ _ \ '_ \| __|    / _ \  | |
 | |  | |  _ <    / ___ \ (_| |  __/ | | | |_    / ___ \ | |
 |_|  |_|_| \_\  /_/   \_\__, |\___|_| |_|\__|  /_/   \_\___|
                          |___/
  Multi-Agent Skill Installer
```

<div align="center">

<p>
<a href="https://github.com/SaidHernandez/mr-agent-ai/releases"><img src="https://img.shields.io/github/v/release/SaidHernandez/mr-agent-ai" alt="Release"></a>
<a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT"></a>
<img src="https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white" alt="Go 1.21+">
<img src="https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey" alt="Platform">
</p>

</div>

---

## What it is

`mr-agent-ai` is a CLI that prepares any project for multi-agent AI development. It ships three independent features:

| Command | What it does | Docs |
|---------|-------------|------|
| `mr-agent-ai install` | Installs structured agent skills into your project | [skills →](internal/skills/docs/README.md) |
| `mr-agent-ai audit` | Scans for supply chain security issues and applies fixes | [audit →](internal/audit/docs/README.md) |
| `mr-agent-ai theme` | Sets up a real-time context panel for your AI agent | [theme →](internal/theme/docs/README.md) |

---

## Why

AI agents are powerful but context-blind by default. Without rules, they invent patterns, leak business logic into controllers, skip Trace-IDs, and write whatever framework style they last saw.

`mr-agent-ai install` gives your agent a precise set of rules — one per architectural layer — before it writes a single line of code. `mr-agent-ai audit` closes the supply chain gaps that most projects accumulate silently. `mr-agent-ai theme` keeps the agent aware of what is happening in the project at all times.

---

## Install

### macOS (Homebrew)

```bash
brew tap SaidHernandez/mr-agent-ai https://github.com/SaidHernandez/mr-agent-ai
brew install mr-agent-ai
```

### Linux (curl)

```bash
curl -fsSL https://raw.githubusercontent.com/SaidHernandez/mr-agent-ai/main/install.sh | bash
```

### Windows (PowerShell)

```powershell
iwr -useb https://raw.githubusercontent.com/SaidHernandez/mr-agent-ai/main/install.ps1 | iex
```

### From source

```bash
git clone https://github.com/SaidHernandez/mr-agent-ai
cd mr-agent-ai
go build -o mr-agent-ai .
```

---

## Commands

### `install` — agent skills

```bash
cd my-project
mr-agent-ai install
```

Interactive prompt to select your AI tools and skills. Generates the right config file for each tool and writes `skills/<layer>/SKILL.md` as the shared source of truth.

Supported tools: Claude Code, Cursor, GitHub Copilot, Windsurf, Cline, Aider, Continue, OpenCode.

→ [Full reference](internal/skills/docs/README.md)

---

### `audit` — supply chain security

```bash
cd my-project
mr-agent-ai audit
```

Scans for missing Dependabot config, insecure `.gitignore` patterns, absent pre-commit hooks, missing SAST workflows, Dockerfile issues, and more. Offers one-step fixes for every detected issue.

→ [Full reference](internal/audit/docs/README.md)

---

### `theme` — context panel

```bash
mr-agent-ai theme init   # configure once
mr-agent-ai theme show   # preview the panel
```

`theme init` detects installed agents (Claude Code, Aider, VS Code) and wires the panel into each one. Once configured, the panel renders automatically in the agent's status bar showing branch, open PR, Linear ticket, token usage, and rate limits.

→ [Full reference](internal/theme/docs/README.md)

---

## Releasing a New Version

Tags and releases are created automatically on every merge to `main` using [conventional commits](https://www.conventionalcommits.org):

| Commit prefix | Version bump |
|---------------|-------------|
| `BREAKING CHANGE` in body | major |
| `feat:` | minor |
| `fix:` / `perf:` | patch |
| `chore:` / `docs:` / `ci:` | no release |

GitHub Actions builds binaries for macOS, Linux, and Windows (amd64/arm64) and updates the Homebrew formula automatically.

---

<div align="center">
<a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT"></a>
</div>

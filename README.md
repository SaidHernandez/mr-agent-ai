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

<p><strong>One command. Any project. Instant multi-agent skills.</strong></p>

<p>
<a href="https://github.com/SaidHernandez/mr-agent-ai/releases"><img src="https://img.shields.io/github/v/release/SaidHernandez/mr-agent-ai" alt="Release"></a>
<a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT"></a>
<img src="https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white" alt="Go 1.21+">
<img src="https://img.shields.io/badge/platform-macOS%20%7C%20Linux-lightgrey" alt="Platform">
</p>

</div>

---

## What It Does

Mr. Agent AI installs a set of **agent skills** — structured instruction sets — directly into your project. Skills are organized in two categories:

- **Architecture skills** — teach your agent how to handle each layer of your stack: API security, controllers, services, repositories, migrations, observability, and QA.
- **Language & framework skills** — teach your agent the best practices for your chosen stack: React, Vue, TypeScript, Java, Python, or Go.

**Before**: Your AI agent writes generic code with no project context.

**After**: Your agent reads the relevant skill before touching any layer, follows your architecture, propagates Trace-IDs, applies the standard error shape, and knows when NOT to put logic in the wrong layer.

---

## Why

AI agents are powerful but context-blind by default. Skills give them the rules they need before writing a single line of code — so they stop inventing patterns and start following yours.

- No more `useMemo` in React 19 projects
- No more business logic leaking into controllers
- No more raw DB errors surfacing to the API
- No more logs without Trace-IDs

---

## Install

```bash
brew tap SaidHernandez/mr-agent-ai https://github.com/SaidHernandez/mr-agent-ai
brew install mr-agent-ai
```

---

## Usage

### Install skills

```bash
cd my-project
mr-agent-ai install
```

**Step 1** — pick your AI coding tools:

```
Select your AI coding tools:

   1.  Claude Code    — Anthropic CLI — reads AGENTS.md automatically
   2.  Cursor         — AI editor — .cursor/rules/<skill>.mdc per skill
   3.  GitHub Copilot — VS Code / JetBrains — .github/copilot-instructions.md
   4.  Windsurf       — Codeium editor — .windsurfrules
   5.  Cline          — VS Code extension — .clinerules
   6.  Aider          — Terminal pair programmer — CONVENTIONS.md
   7.  Continue       — Open-source VS Code / JetBrains — .continuerules
   8.  OpenCode       — SST terminal AI coder — AGENTS.md compatible

  Enter numbers (e.g. 1,3 or 'all'): 1
```

**Step 2** — pick your architecture skills:

```
Select which architecture skills to install:

   1.  orchestrator   — Coordinates agents, delegates tasks, manages state.
   2.  api-security   — CORS, rate-limiting, JWT authentication patterns.
   3.  controller     — API endpoints, DTOs, standard error responses.
   4.  service-layer  — Business logic, DI/IoC, domain exceptions.
   5.  repository     — SQL queries, transactions, DB error mapping.
   6.  migration      — PostgreSQL schema design, indexes, migrations.
   7.  observability  — Structured logging, Trace-ID propagation.
   8.  qa             — Unit and E2E tests, Playwright, coverage.
   9.  skill-creator  — Creates new skills following the standard spec.
  10.  supply-chain   — Detects supply chain risks and configures shelf-time protection.

  Enter numbers (e.g. 1,3 or 'all'): all
```

**Step 3** — pick a language or framework skill (optional, single choice):

```
Select a programming language skill (optional):

   1.  react       — React 19, Zustand, Tailwind, Server Components.
   2.  vue         — Vue 3, Composition API, Pinia, composables.
   3.  typescript  — Strict mode, type inference, type guards.
   4.  java        — Optional, null safety, clean API design.
   5.  python      — Type hints, Pythonic patterns, error handling.
   6.  golang      — Error handling, goroutines, idiomatic Go.

  Enter a number (or press Enter to skip): 1
```

### Audit supply chain security

Scans the project for missing security configurations and offers to apply fixes:

```bash
cd my-project
mr-agent-ai audit
```

```
  Supply Chain Audit
  ──────────────────────────────────────

  Ecosystems   Go · Node.js

  Supply Chain
  ✓  Dependabot
  ✗  npm shelf-time         .npmrc → min-release-age=7d

  Secrets
  ✗  .gitignore             missing: .env · *.pem · *.key
  ✗  secret scanner         gitleaks or trufflehog not found

  SAST
  ✗  gosec                  not in CI
  ✗  CodeQL                 not configured

  ──────────────────────────────────────
  4 issues

  Available fixes
  ───────────────
  1  .gitignore
  2  .npmrc
  3  .gitleaks.toml
  4  .github/workflows/sast.yml
  5  .github/workflows/codeql.yml

  Apply (1,3 or all) › _
```

---

Result in your project:

```
my-project/
├── AGENTS.md                          ← Claude Code / OpenCode index
├── .cursor/rules/                     ← Cursor (one .mdc per skill)
├── .github/copilot-instructions.md    ← GitHub Copilot
├── .windsurfrules                     ← Windsurf
├── .clinerules                        ← Cline
├── CONVENTIONS.md                     ← Aider
├── .continue/rules.md                 ← Continue
└── skills/
    ├── arch/                          ← architecture skills
    │   ├── orchestrator/SKILL.md
    │   ├── api-security/SKILL.md
    │   ├── controller/SKILL.md
    │   │   └── assets/ERROR_RESPONSE.ts
    │   ├── service-layer/SKILL.md
    │   ├── repository/SKILL.md
    │   ├── migration/SKILL.md
    │   │   └── assets/MIGRATION_TEMPLATE.sql
    │   ├── observability/SKILL.md
    │   ├── qa/SKILL.md
    │   └── skill-creator/SKILL.md
    │       └── assets/SKILL-TEMPLATE.md
    └── lang/                          ← language / framework skill
        └── react/SKILL.md
```

---

## Skills

### Architecture

| Skill | Trigger | Key Rules |
|-------|---------|-----------|
| `orchestrator` | Coordinating multiple agents | Trace-ID generation, no business logic |
| `api-security` | CORS, auth, rate limiting | Explicit allowlists, token-bucket, JWT algorithm validation |
| `controller` | API endpoints, DTOs | Zod v4, `as const` error codes, OpenAPI 3.0, standard error shape |
| `service-layer` | Business logic | DI/IoC via interfaces, `DomainException`, no SQL, no HTTP |
| `repository` | Database operations | No `SELECT *`, cursor pagination, transactions, PG error mapping |
| `migration` | Schema changes | `COMMENT ON COLUMN` required, FK indexes, reversible with `down` |
| `observability` | Logging, metrics | Trace-ID every log line, structured JSON, no PII |
| `qa` | Tests | `test.each` for edge cases, Playwright POM, ≥85% branch coverage |
| `skill-creator` | Creating new skills | Template, naming conventions, checklist |
| `supply-chain` | Adding/updating dependencies | Shelf-time configs, red flags, secret scanning, audit commands |

### Language & Framework

| Skill | Key Rules |
|-------|-----------|
| `react` | React 19 Compiler, no manual memoization, Zustand 5, Tailwind cn() |
| `vue` | Composition API, `<script setup>`, Pinia + storeToRefs, composables |
| `typescript` | Strict mode, no `any`, `as const`, explicit return types |
| `java` | Optional as return type only, orElseThrow, no nested Optionals |
| `python` | Type hints, no mutable defaults, context managers, specific exceptions |
| `golang` | Handle all errors, defer cleanup, small interfaces, context propagation |

---

## How It Works

Once installed, tell your agent to read `AGENTS.md` before working on any layer. The agent finds the matching skill trigger, reads that `SKILL.md`, and follows its rules for the duration of the task.

```
You: "create the user registration endpoint"
Agent: reads AGENTS.md → loads controller + service-layer + repository skills
Agent: uses Zod v4 DTO, delegates to service, service throws DomainException,
       repository maps PG errors, all logs carry Trace-ID
```

---

## Releasing a New Version

Tags and releases are created automatically on every merge to `main` using [conventional commits](https://www.conventionalcommits.org):

| Commit prefix | Version bump |
|---------------|-------------|
| `BREAKING CHANGE` in body | major |
| `feat:` | minor |
| `fix:` / `perf:` | patch |
| `chore:` / `docs:` / `ci:` | no release |

GitHub Actions builds binaries for macOS + Linux (amd64/arm64) and updates the Homebrew formula automatically.

---

## Credits

Skill structure and conventions inspired by [Gentleman-Skills](https://github.com/Gentleman-Programming/Gentleman-Skills).

---

<div align="center">
<a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT"></a>
</div>

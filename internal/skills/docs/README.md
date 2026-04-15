# skills

## Context

AI agents write generic code by default — they have no knowledge of your architecture, your error shapes, or your logging conventions. This package solves that by installing a set of **skill files** into your project: one per architectural layer and one for your language stack. Each skill is a structured markdown file the agent reads before touching that layer, so it follows your rules instead of inventing its own.

## When to use

Run once when setting up a new project, or when onboarding a new AI tool to an existing one:

```bash
cd my-project
mr-agent-ai install
```

### Step 1 — select your AI tools

```
Select your AI coding tools:

   1.  Claude Code    — Anthropic CLI — reads AGENTS.md + skills/ automatically
   2.  Cursor         — AI editor — .cursor/rules/<skill>.mdc per skill
   3.  GitHub Copilot — VS Code / JetBrains — .github/copilot-instructions.md
   4.  Windsurf       — Codeium editor — .windsurfrules
   5.  Cline          — VS Code extension — .clinerules
   6.  Aider          — Terminal pair programmer — CONVENTIONS.md
   7.  Continue       — Open-source VS Code / JetBrains — .continuerules
   8.  OpenCode       — SST terminal AI coder — AGENTS.md compatible

  Enter numbers (e.g. 1,3 or 'all'): 1
```

### Step 2 — select architecture skills

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

### Step 3 — select a language skill (optional)

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

### Result

```
[mr-agent-ai] Installing 11 skill(s) for 1 tool(s) into: /my-project

  [ok] skills/arch/orchestrator/SKILL.md
  [ok] skills/arch/api-security/SKILL.md
  [ok] skills/arch/controller/SKILL.md
  [ok] skills/arch/controller/assets/ERROR_RESPONSE.ts
  [ok] skills/arch/service-layer/SKILL.md
  [ok] skills/arch/repository/SKILL.md
  [ok] skills/arch/migration/SKILL.md
  [ok] skills/arch/migration/assets/MIGRATION_TEMPLATE.sql
  [ok] skills/arch/observability/SKILL.md
  [ok] skills/arch/qa/SKILL.md
  [ok] skills/arch/skill-creator/SKILL.md
  [ok] skills/arch/skill-creator/assets/SKILL-TEMPLATE.md
  [ok] skills/arch/supply-chain/SKILL.md
  [ok] skills/lang/react/SKILL.md

  [ok] AGENTS.md  (Claude Code)
```

## What it generates

### Shared source of truth

| File | Description |
|------|-------------|
| `skills/arch/<skill>/SKILL.md` | Instructions for each architecture layer |
| `skills/arch/controller/assets/ERROR_RESPONSE.ts` | Standard error response shape |
| `skills/arch/migration/assets/MIGRATION_TEMPLATE.sql` | Migration file template |
| `skills/arch/skill-creator/assets/SKILL-TEMPLATE.md` | Template for creating new skills |
| `skills/lang/<skill>/SKILL.md` | Language/framework best practices |

### Tool-specific config files

| Tool | File generated |
|------|---------------|
| Claude Code | `AGENTS.md` |
| OpenCode | `AGENTS.md` |
| Cursor | `.cursor/rules/<skill>.mdc` per skill |
| GitHub Copilot | `.github/copilot-instructions.md` |
| Windsurf | `.windsurfrules` |
| Cline | `.clinerules` |
| Aider | `CONVENTIONS.md` |
| Continue | `.continue/rules.md` |

## Adding a new skill

1. Create `content/arch/<name>.md` or `content/lang/<name>.md`
2. Add the entry to `skills_arch.go` or `skills_lang.go`

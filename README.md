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

Mr. Agent AI installs a set of **agent skills** — structured instruction sets — directly into your project. Each skill teaches your AI coding agent how to handle a specific layer of your stack: React frontend, API security, controllers, services, repositories, migrations, observability, and QA.

**Before**: Your AI agent writes generic code with no project context.

**After**: Your agent reads the relevant skill before touching any layer, follows your architecture, propagates Trace-IDs, uses Zod v4, applies the standard error shape, and knows when NOT to put logic in the wrong layer.

---

## Why

AI agents are powerful but context-blind by default. Skills give them the rules they need before writing a single line of code — so they stop inventing patterns and start following yours.

- No more `useMemo` in React 19 projects
- No more `z.string().email()` when you're on Zod v4
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

Run inside any project:

```bash
cd my-project
mr-agent-ai install
```

The CLI prompts you to select which skills to install:

```
  Install orchestrator       — Central coordination agent [Y/n]:
  Install frontend           — React 19, Tailwind, Zustand, XSS prevention [Y/n]:
  Install api-security       — CORS, rate limiting, JWT [Y/n]:
  Install controller         — Zod v4 DTOs, OpenAPI, error responses [Y/n]:
  Install service-layer      — Business logic, DI/IoC, domain exceptions [Y/n]:
  Install repository         — SQL optimization, transactions, idempotency [Y/n]:
  Install migration          — PostgreSQL schema design, indexing [Y/n]:
  Install observability      — Structured logging, Trace-ID propagation [Y/n]:
  Install qa                 — Unit tests, Playwright POM, minimal mocking [Y/n]:
  Install skill-creator      — Creates new agent skills [Y/n]:
```

Result in your project:

```
my-project/
├── AGENTS.md                          ← load this first — skill index
└── skills/
    ├── orchestrator/SKILL.md
    ├── frontend/SKILL.md
    ├── api-security/SKILL.md
    ├── controller/SKILL.md
    │   └── assets/ERROR_RESPONSE.ts
    ├── service-layer/SKILL.md
    ├── repository/SKILL.md
    ├── migration/SKILL.md
    │   └── assets/MIGRATION_TEMPLATE.sql
    ├── observability/SKILL.md
    ├── qa/SKILL.md
    └── skill-creator/SKILL.md
        └── assets/SKILL-TEMPLATE.md
```

---

## Skills

| Skill | Trigger | Key Rules |
|-------|---------|-----------|
| `orchestrator` | Coordinating multiple agents | Trace-ID generation, no business logic |
| `frontend` | React components, design systems | React 19 Compiler, no manual memoization, Tailwind cn(), Zustand 5 |
| `api-security` | CORS, auth, rate limiting | Explicit allowlists, token-bucket, JWT algorithm validation |
| `controller` | API endpoints, DTOs | Zod v4, `as const` error codes, OpenAPI 3.0, standard error shape |
| `service-layer` | Business logic | DI/IoC via interfaces, `DomainException`, no SQL, no HTTP |
| `repository` | Database operations | No `SELECT *`, cursor pagination, transactions, PG error mapping |
| `migration` | Schema changes | `COMMENT ON COLUMN` required, FK indexes, reversible with `down` |
| `observability` | Logging, metrics | Trace-ID every log line, structured JSON, no PII |
| `qa` | Tests | `test.each` for edge cases, Playwright POM, ≥85% branch coverage |
| `skill-creator` | Creating new skills | Template, naming conventions, checklist |

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

```bash
git tag v1.x.x
git push origin v1.x.x
```

GitHub Actions builds binaries for macOS + Linux (amd64/arm64) and updates the Homebrew formula automatically.

---

<div align="center">
<a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT"></a>
</div>

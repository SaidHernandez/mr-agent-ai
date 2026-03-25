# mr-agent-ai

Multi-Agent AI Skill Installer. Installs agent skills into any project with a single command.

## Install

```bash
brew tap SaidHernandez/mr-agent-ai https://github.com/SaidHernandez/mr-agent-ai
brew install mr-agent-ai
```

## Usage

Run inside any project:

```bash
cd my-project
mr-agent-ai install
```

The CLI asks which skills to install and creates:

```
my-project/
├── AGENTS.md                        ← skill index — load this first
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

## Skills

| Skill | Trigger |
|-------|---------|
| `orchestrator` | Coordinating multiple agents or managing workflow state |
| `frontend` | React 19 components, Tailwind, Zustand, XSS prevention |
| `api-security` | CORS, rate limiting, JWT |
| `controller` | API endpoints, Zod v4 DTOs, error responses |
| `service-layer` | Business logic, DI/IoC, domain exceptions |
| `repository` | SQL optimization, transactions, idempotency |
| `migration` | PostgreSQL schema design, indexing |
| `observability` | Structured logging, Trace-ID propagation |
| `qa` | Unit tests, Playwright E2E, minimal mocking |
| `skill-creator` | Creating new agent skills |

## Release a new version

```bash
git tag v1.0.0
git push origin v1.0.0
```

GitHub Actions runs GoReleaser, builds binaries for macOS + Linux (amd64/arm64), and updates the Homebrew tap automatically.

## Setup (first time)

1. Push a tag — GoReleaser builds the binaries and writes `Formula/mr-agent-ai.rb` directly into this repo.
2. No extra secrets or repos needed — uses the built-in `GITHUB_TOKEN`.

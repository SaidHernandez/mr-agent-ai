package skills

// archSkills contains all architecture-layer skills (orchestrator, frontend, API, etc.).
var archSkills = []Skill{
	{
		Name:        "orchestrator",
		Dir:         "arch/orchestrator",
		Description: "Coordinates agents, delegates tasks, manages state.",
		Trigger:     "When coordinating multiple agents or managing workflow state.",
		Content:     orchestratorSkill,
	},
	{
		Name:        "api-security",
		Dir:         "arch/api-security",
		Description: "CORS, rate-limiting, JWT authentication patterns.",
		Trigger:     "When configuring CORS, rate limits, JWT tokens, or API authentication.",
		Content:     apiSecuritySkill,
	},
	{
		Name:        "controller",
		Dir:         "arch/controller",
		Description: "API endpoints, DTOs, standard error responses.",
		Trigger:     "When creating API endpoints, validating DTOs, or standardizing error responses.",
		Content:     controllerSkill,
		Assets: []Asset{
			{Path: "ERROR_RESPONSE.ts", Content: mustReadAsset("ERROR_RESPONSE.ts")},
		},
	},
	{
		Name:        "service-layer",
		Dir:         "arch/service-layer",
		Description: "Business logic, DI/IoC, domain exceptions.",
		Trigger:     "When implementing business logic, applying DI/IoC, or handling domain exceptions.",
		Content:     serviceLayerSkill,
	},
	{
		Name:        "repository",
		Dir:         "arch/repository",
		Description: "SQL queries, transactions, DB error mapping.",
		Trigger:     "When performing database operations, managing transactions, or optimizing SQL queries.",
		Content:     repositorySkill,
	},
	{
		Name:        "migration",
		Dir:         "arch/migration",
		Description: "PostgreSQL schema design, indexes, migrations.",
		Trigger:     "When creating DB migrations, designing schemas, or managing indexes.",
		Content:     migrationSkill,
		Assets: []Asset{
			{Path: "MIGRATION_TEMPLATE.sql", Content: mustReadAssetWithDate("MIGRATION_TEMPLATE.sql")},
		},
	},
	{
		Name:        "observability",
		Dir:         "arch/observability",
		Description: "Structured logging, Trace-ID propagation.",
		Trigger:     "When adding logging, propagating Trace-IDs, or collecting metrics.",
		Content:     observabilitySkill,
	},
	{
		Name:        "qa",
		Dir:         "arch/qa",
		Description: "Unit and E2E tests, Playwright, coverage.",
		Trigger:     "When writing unit tests, E2E tests, covering edge cases, or reviewing test coverage.",
		Content:     qaSkill,
	},
	{
		Name:        "skill-creator",
		Dir:         "arch/skill-creator",
		Description: "Creates new skills following the standard spec.",
		Trigger:     "When asked to create a new skill, add agent instructions, or document patterns for AI.",
		Content:     skillCreatorSkill,
		Assets: []Asset{
			{Path: "SKILL-TEMPLATE.md", Content: mustReadAsset("SKILL-TEMPLATE.md")},
		},
	},
	{
		Name:        "supply-chain",
		Dir:         "arch/supply-chain",
		Description: "Detects supply chain risks and configures shelf-time protection.",
		Trigger:     "When adding or updating dependencies, reviewing dependency PRs, or auditing packages.",
		Content:     supplyChainSkill,
	},
}

// ─── Skill content functions ───────────────────────────────────────────────────

func orchestratorSkill() string {
	return skillFrontmatter(
		"orchestrator",
		"Central coordination agent. Routes tasks to sub-agents, tracks state, propagates Trace-ID.",
		"When coordinating multiple agents or managing workflow state.",
		"Read, Write, Bash, Task",
	) + mustReadDoc("arch/orchestrator.md")
}

func apiSecuritySkill() string {
	return skillFrontmatter(
		"api-security",
		"API Security agent. CORS allowlists, token-bucket rate limiting, JWT HS256 validation.",
		"When configuring CORS, rate limits, JWT tokens, or API authentication.",
		"Read, Edit, Write, Glob, Grep",
	) + mustReadDoc("arch/api-security.md")
}

func controllerSkill() string {
	return skillFrontmatter(
		"controller",
		"Controller agent. OpenAPI-compatible endpoints, Zod v4 DTOs, TypeScript as-const error codes.",
		"When creating API endpoints, validating DTOs, or standardizing error responses.",
		"Read, Edit, Write, Glob, Grep",
	) + mustReadDoc("arch/controller.md")
}

func serviceLayerSkill() string {
	return skillFrontmatter(
		"service-layer",
		"Service Layer agent. Business logic only, DI/IoC, DomainException, no SQL, no HTTP.",
		"When implementing business logic, applying DI/IoC, or handling domain exceptions.",
		"Read, Edit, Write, Glob, Grep",
	) + mustReadDoc("arch/service-layer.md")
}

func repositorySkill() string {
	return skillFrontmatter(
		"repository",
		"Repository agent. SQL optimization, idempotency keys, cursor pagination, DB error mapping.",
		"When performing database operations, managing transactions, or optimizing SQL queries.",
		"Read, Edit, Write, Glob, Grep",
	) + mustReadDoc("arch/repository.md")
}

func migrationSkill() string {
	return skillFrontmatter(
		"migration",
		"Migration agent. PostgreSQL schema design, FK indexes, COMMENT ON COLUMN, reversible migrations.",
		"When creating DB migrations, designing schemas, or managing indexes.",
		"Read, Edit, Write, Glob, Grep, Bash",
	) + mustReadDoc("arch/migration.md")
}

func observabilitySkill() string {
	return skillFrontmatter(
		"observability",
		"Observability agent. Structured JSON logging, Trace-ID propagation across all layers.",
		"When adding logging, propagating Trace-IDs, or collecting metrics.",
		"Read, Edit, Write, Glob, Grep",
	) + mustReadDoc("arch/observability.md")
}

func qaSkill() string {
	return skillFrontmatter(
		"qa",
		"QA agent. Unit + E2E tests, Playwright POM, table-driven parametrize, minimal mocking.",
		"When writing unit tests, E2E tests, covering edge cases, or reviewing test coverage.",
		"Read, Edit, Write, Glob, Grep, Bash",
	) + mustReadDoc("arch/qa.md")
}

func skillCreatorSkill() string {
	return skillFrontmatter(
		"skill-creator",
		"Meta-skill. Creates new agent skills following the standard spec and naming conventions.",
		"When asked to create a new skill, add agent instructions, or document patterns for AI.",
		"Read, Edit, Write, Glob, Grep",
	) + mustReadDoc("arch/skill-creator.md")
}

func supplyChainSkill() string {
	return skillFrontmatter(
		"supply-chain",
		"Supply chain security agent. Shelf-time configs, dependency review, secret scanning, audit commands.",
		"When adding dependencies, reviewing dependency update PRs, or auditing packages.",
		"Read, Glob, Grep, Bash",
	) + mustReadDoc("arch/supply-chain.md")
}

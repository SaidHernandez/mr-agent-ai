package installer

import (
	"strings"
	"time"
)

// bt replaces TICK placeholder with a real backtick.
// Used because Go raw strings can't contain backtick characters.
func bt(s string) string {
	return strings.NewReplacer("TICK", "`").Replace(s)
}

var fence = "```"

// Asset is an optional file placed under skills/<agent>/assets/.
type Asset struct {
	Path    string // relative to skills/<agent>/assets/
	Content string
}

// Agent defines a single agent skill to be installed.
type Agent struct {
	Name        string
	Dir         string
	Description string
	Trigger     string
	Content     func() string
	Assets      []Asset
}

var agents = []Agent{
	{
		Name:        "orchestrator",
		Dir:         "orchestrator",
		Description: "Central coordination agent. Delegates tasks to sub-agents and manages workflow state.",
		Trigger:     "When coordinating multiple agents or managing workflow state.",
		Content:     orchestratorSkill,
	},
	{
		Name:        "frontend",
		Dir:         "frontend",
		Description: "React 19 agent. Component architecture, design systems, Tailwind cn(), no manual memoization, XSS prevention.",
		Trigger:     "When building React components, implementing design systems, managing state, or handling frontend security.",
		Content:     frontendSkill,
	},
	{
		Name:        "api-security",
		Dir:         "api-security",
		Description: "API Security agent. CORS allowlists, rate-limiting, JWT strategies.",
		Trigger:     "When configuring CORS, rate limits, JWT tokens, or API authentication.",
		Content:     apiSecuritySkill,
	},
	{
		Name:        "controller",
		Dir:         "controller",
		Description: "Controller agent. OpenAI-compatible endpoints, Zod v4 DTOs, TypeScript const types, standard error shape.",
		Trigger:     "When creating API endpoints, validating DTOs, or standardizing error responses.",
		Content:     controllerSkill,
		Assets: []Asset{
			{Path: "ERROR_RESPONSE.ts", Content: errorResponseAsset()},
		},
	},
	{
		Name:        "service-layer",
		Dir:         "service-layer",
		Description: "Service Layer agent. Business logic, DI/IoC, TypeScript const types, domain exception handling.",
		Trigger:     "When implementing business logic, applying DI/IoC, or handling domain exceptions.",
		Content:     serviceLayerSkill,
	},
	{
		Name:        "repository",
		Dir:         "repository",
		Description: "Repository agent. SQL optimization, idempotency, transactions, DB exception mapping.",
		Trigger:     "When performing database operations, managing transactions, or optimizing SQL queries.",
		Content:     repositorySkill,
	},
	{
		Name:        "migration",
		Dir:         "migration",
		Description: "Migration agent. PostgreSQL schema design, indexing strategies, cardinality validation.",
		Trigger:     "When creating DB migrations, designing schemas, or managing indexes.",
		Content:     migrationSkill,
		Assets: []Asset{
			{Path: "MIGRATION_TEMPLATE.sql", Content: migrationTemplateAsset()},
		},
	},
	{
		Name:        "observability",
		Dir:         "observability",
		Description: "Observability agent. Structured logging with Trace-ID propagation across all layers.",
		Trigger:     "When adding logging, propagating Trace-IDs, or collecting metrics.",
		Content:     observabilitySkill,
	},
	{
		Name:        "qa",
		Dir:         "qa",
		Description: "QA agent. Unit + E2E tests, Playwright POM, parametrize patterns, minimal mocking.",
		Trigger:     "When writing unit tests, E2E tests, covering edge cases, or reviewing test coverage.",
		Content:     qaSkill,
	},
	{
		Name:        "skill-creator",
		Dir:         "skill-creator",
		Description: "Meta-skill. Creates new agent skills following the standard spec.",
		Trigger:     "When asked to create a new skill, add agent instructions, or document patterns for AI.",
		Content:     skillCreatorSkill,
		Assets: []Asset{
			{Path: "SKILL-TEMPLATE.md", Content: skillTemplateAsset()},
		},
	},
}

func skillFrontmatter(name, description, trigger, allowedTools string) string {
	return "---\n" +
		"name: mr-agent-" + name + "\n" +
		"description: >\n" +
		"  " + description + "\n" +
		"  Trigger: " + trigger + "\n" +
		"license: MIT\n" +
		"allowed-tools: " + allowedTools + "\n" +
		"metadata:\n" +
		"  version: \"1.0\"\n" +
		"  generated: " + time.Now().Format("2006-01-02") + "\n" +
		"---\n\n"
}

// ─── Assets ───────────────────────────────────────────────────────────────────

func errorResponseAsset() string {
	return `// Standard error response shape — use this in every controller error path.
// See: skills/controller/SKILL.md

export interface ErrorResponse {
  success: false;
  code: string;      // SCREAMING_SNAKE_CASE — machine-readable
  message: string;   // Human-readable, safe to display in UI
  detail: string;    // Technical context for debugging — never PII
}

export interface SuccessResponse<T = unknown> {
  success: true;
  data: T;
}

export type ApiResponse<T = unknown> = SuccessResponse<T> | ErrorResponse;

// ── Error code registry ───────────────────────────────────────────────────────
export const ERROR_CODES = {
  VALIDATION_ERROR:      'VALIDATION_ERROR',
  UNAUTHORIZED:          'UNAUTHORIZED',
  FORBIDDEN:             'FORBIDDEN',
  NOT_FOUND:             'NOT_FOUND',
  CONFLICT:              'CONFLICT',
  UNPROCESSABLE:         'UNPROCESSABLE',
  RATE_LIMIT_EXCEEDED:   'RATE_LIMIT_EXCEEDED',
  INTERNAL_ERROR:        'INTERNAL_ERROR',
} as const;

export type ErrorCode = (typeof ERROR_CODES)[keyof typeof ERROR_CODES];

// ── HTTP status map ───────────────────────────────────────────────────────────
export const HTTP_STATUS: Record<ErrorCode, number> = {
  VALIDATION_ERROR:    400,
  UNAUTHORIZED:        401,
  FORBIDDEN:           403,
  NOT_FOUND:           404,
  CONFLICT:            409,
  UNPROCESSABLE:       422,
  RATE_LIMIT_EXCEEDED: 429,
  INTERNAL_ERROR:      500,
};
`
}

func migrationTemplateAsset() string {
	return `-- Migration: <migration_name>
-- Created:   ` + time.Now().Format("2006-01-02") + `
-- Description: <what this migration does>

-- ── up.sql ───────────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS <table_name> (
  id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  external_id UUID        NOT NULL DEFAULT gen_random_uuid() UNIQUE,
  -- Add your columns here
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Column descriptions (REQUIRED for every field)
COMMENT ON COLUMN <table_name>.id          IS 'Internal surrogate key — never exposed externally';
COMMENT ON COLUMN <table_name>.external_id IS 'UUID exposed in API responses — immutable after creation';
COMMENT ON COLUMN <table_name>.created_at  IS 'Row creation timestamp with timezone';
COMMENT ON COLUMN <table_name>.updated_at  IS 'Last update timestamp with timezone';

-- Indexes (add for every FK and high-cardinality filter column)
-- CREATE UNIQUE INDEX idx_<table>_<col>  ON <table_name> (<col>);
-- CREATE INDEX        idx_<table>_fk     ON <table_name> (foreign_key_col);

-- ── down.sql ─────────────────────────────────────────────────────────────────

-- DROP TABLE IF EXISTS <table_name>;
`
}

func skillTemplateAsset() string {
	return `---
name: mr-agent-{skill-name}
description: >
  {One-line description of what this skill enables}.
  Trigger: {When the AI should load this skill — be specific}.
license: MIT
allowed-tools: Read, Edit, Write, Glob, Grep, Bash
metadata:
  version: "1.0"
---

## When to Use

Use this skill when:
- {Condition 1}
- {Condition 2}
- {Condition 3}

**Don't use when:** {Counter-cases}

---

## Decision Tree

` + "```" + `
{Question 1}?  → {Action A}
{Question 2}?  → {Action B}
Otherwise      → {Default action}
` + "```" + `

---

## Critical Patterns

### Pattern 1: {Name}

` + "```" + `{language}
// ✅ GOOD
{correct example}

// ❌ BAD
{incorrect example}
` + "```" + `

### Pattern 2: {Name}

` + "```" + `{language}
{example}
` + "```" + `

---

## Commands

` + "```" + `bash
{command 1}  # {description}
{command 2}  # {description}
` + "```" + `

---

## Definition of Done

- [ ] {Requirement 1}
- [ ] {Requirement 2}
- [ ] {Requirement 3}

---

## Resources

- **Templates**: See [assets/](assets/) for {description}
`
}

// ─── Orchestrator ─────────────────────────────────────────────────────────────

func orchestratorSkill() string {
	return skillFrontmatter(
		"orchestrator",
		"Central coordination agent. Routes tasks to sub-agents, tracks state, propagates Trace-ID.",
		"When coordinating multiple agents or managing workflow state.",
		"Read, Write, Bash, Task",
	) + bt(`# Orchestrator Agent

## When to Use

- Coordinating multiple sub-agents for a single request
- Managing workflow state across steps
- Propagating Trace-ID to all downstream agents

**Don't use when:** The task is isolated to one layer — use that layer's skill directly.

---

## Decision Tree

`) + fence + `
Has multiple sub-agents involved?  → Generate Trace-ID, delegate with it
Sub-agent returned an error?       → Log + map to standard error shape, don't retry silently
Is this business logic?            → Delegate to Service Layer — NEVER implement here
` + fence + bt(`

---

## Non-Negotiable Rules

1. Generate a unique **Trace-ID** (TICKuuid v4TICK) for every incoming request — pass it to every sub-agent call.
2. Never contain business logic — only routing and state transitions.
3. On sub-agent failure, log with Trace-ID and return:

`) + fence + `json
{ "success": false, "code": "ORCHESTRATION_ERROR", "message": "Sub-agent failed", "detail": "<agent>: <reason>" }
` + fence + bt(`

4. Keep workflow state immutable per step — no shared mutable state between agents.

---

## Delegation Pattern

`) + fence + `
[Orchestrator]
  ├── receives: { intent, payload, traceId }
  ├── dispatches to: Controller → Service → Repository
  ├── each hop receives: { ...payload, traceId }
  └── returns: { success, code, message, detail, data }
` + fence + bt(`

---

## Log Convention

`) + fence + `
[<traceId>] INFO  orchestrator — dispatching to <agent>, intent=<intent>
[<traceId>] ERROR orchestrator — <agent> failed: <reason>
` + fence + `

---

## Commands

` + fence + `bash
# Generate a Trace-ID manually for testing
node -e "console.log(require('crypto').randomUUID())"
` + fence + `

## Definition of Done

- [ ] Every request has a Trace-ID from this layer down
- [ ] No business logic in orchestrator files
- [ ] Sub-agent failures caught and mapped to standard error shape
- [ ] Workflow state logged at each transition
`
}

// ─── Frontend ─────────────────────────────────────────────────────────────────

func frontendSkill() string {
	return skillFrontmatter(
		"frontend",
		"React 19 agent. No manual memoization, Server Components first, Tailwind cn(), Zustand 5, XSS prevention.",
		"When building React components, implementing design systems, managing state, or handling frontend security.",
		"Read, Edit, Write, Glob, Grep, Bash, WebFetch",
	) + bt(`# Frontend Agent

## When to Use

- Building or refactoring React components
- Implementing design systems or Tailwind styles
- Managing client-side state with Zustand
- Handling frontend security (XSS, CSP, obfuscation)

> References: anthropics/frontend-design · google-labs-code/stitch-skills · vercel-labs/react-best-practices · Gentleman-Skills/react-19 · Gentleman-Skills/tailwind-4 · Gentleman-Skills/zustand-5

---

## Decision Tree

`) + fence + `
Need state that persists across pages?    → Zustand with persist middleware
Need state only inside one component?     → useState
Need async data?                          → Server Component (default) or useQuery hook
Conditional classes?                      → cn("base", condition && "variant")
Static classes only?                      → className="..." directly — no cn() needed
Dynamic non-Tailwind value?               → style={{ width: `+"`"+`${x}%`+"`"+` }}
Code that must not be readable?           → JS Obfuscator at build time (not on components)
` + fence + `

---

## React 19: No Manual Memoization (REQUIRED)

React Compiler handles optimization automatically. Never add ` + "`useMemo`" + ` or ` + "`useCallback`" + ` manually.

` + fence + `typescript
// ✅ CORRECT — React Compiler optimizes this automatically
function ProductList({ items }) {
  const active = items.filter(x => x.active);
  const handleClick = (id) => console.log(id);
  return <List items={active} onClick={handleClick} />;
}

// ❌ NEVER — unnecessary manual memoization
const active = useMemo(() => items.filter(x => x.active), [items]);
const handleClick = useCallback((id) => console.log(id), []);
` + fence + bt(`

---

## Server Components First (Next.js / React 19)

`) + fence + `typescript
// ✅ Server Component — default, no directive, async, zero client JS
export default async function Page() {
  const data = await fetchData();
  return <ClientComponent data={data} />;
}

// ✅ Client Component — only when you need interactivity
"use client";
export function Counter() {
  const [n, setN] = useState(0);
  return <button onClick={() => setN(n + 1)}>{n}</button>;
}
` + fence + bt(`

**Use TICK"use client"TICK only for:** TICKuseStateTICK, TICKuseEffectTICK, event handlers, browser APIs (TICKwindowTICK, TICKlocalStorageTICK).

---

## React 19 Patterns

**TICKuse()TICK hook — read promises and context conditionally:**
`) + fence + `typescript
import { use } from "react";

function Comments({ promise }) {
  const comments = use(promise); // suspends until resolved
  return comments.map(c => <div key={c.id}>{c.text}</div>);
}
` + fence + bt(`

**TICKref TICK as prop — no more TICKforwardRefTICK:**
`) + fence + `typescript
// ✅ React 19
function Input({ ref, ...props }) {
  return <input ref={ref} {...props} />;
}

// ❌ Old pattern — no longer needed
const Input = forwardRef((props, ref) => <input ref={ref} {...props} />);
` + fence + bt(`

**TICKuseActionStateTICK for forms:**
`) + fence + `typescript
"use server";
async function submitForm(formData: FormData) {
  await saveToDatabase(formData);
  revalidatePath("/");
}

"use client";
function Form() {
  const [state, action, isPending] = useActionState(submitForm, null);
  return (
    <form action={action}>
      <button disabled={isPending}>{isPending ? "Saving..." : "Save"}</button>
    </form>
  );
}
` + fence + `

---

## Component Architecture

` + fence + `
src/
├── components/        # Pure UI — no data fetching, no business logic
│   └── Button/
│       ├── Button.tsx
│       ├── Button.test.tsx
│       └── index.ts
├── hooks/             # Logic, side effects, API calls — no JSX
├── features/          # Domain-scoped composition (auth, dashboard)
├── data/              # Static content, API schemas, type definitions
└── styles/            # Design tokens, global CSS
` + fence + bt(`

**Rules:**
- One component per file. No barrel files with side effects.
- Props must use TICKReadonly<T>TICK interfaces.
- Static content goes in TICKdata/TICK — never hardcoded in JSX.
- Imports: TICKTICK named only TICK— never TICKimport React from "react"TICK.

---

## Tailwind 4 — Critical Rules

`) + fence + `typescript
// ❌ NEVER: var() in className
<div className="bg-[var(--color-primary)]" />
<div className="text-[var(--text-color)]" />

// ✅ ALWAYS: Tailwind semantic classes
<div className="bg-primary" />
<div className="text-slate-400" />

// ❌ NEVER: hex colors in className
<p className="text-[#ffffff]" />

// ✅ ALWAYS: Tailwind color classes
<p className="text-white" />
` + fence + bt(`

**TICKcn()TICK utility — use for conditional classes only:**
`) + fence + `typescript
import { clsx } from "clsx";
import { twMerge } from "tailwind-merge";
export const cn = (...inputs: ClassValue[]) => twMerge(clsx(inputs));

// ✅ Use cn() for conditional or mergeable classes
<div className={cn("rounded-lg border", isActive && "bg-blue-500", className)} />

// ❌ Don't wrap static classes in cn() — unnecessary
<div className={cn("flex items-center gap-2")} />  // just use className directly
` + fence + `

---

## Zustand 5 — State Management

` + fence + `typescript
import { create } from "zustand";
import { useShallow } from "zustand/react/shallow";

interface UserStore {
  user: User | null;
  loading: boolean;
  fetchUser: (id: string) => Promise<void>;
}

const useUserStore = create<UserStore>((set) => ({
  user: null,
  loading: false,
  fetchUser: async (id) => {
    set({ loading: true });
    const user = await api.getUser(id);
    set({ user, loading: false });
  },
}));

// ✅ Select specific fields — prevents re-render on unrelated changes
const name = useUserStore((state) => state.name);

// ✅ Multiple fields — use useShallow
const { name, email } = useUserStore(useShallow((s) => ({ name: s.name, email: s.email })));

// ❌ Selecting entire store — re-renders on ANY state change
const store = useUserStore();
` + fence + `

---

## Performance (vercel-labs/react-best-practices — CRITICAL)

**Eliminate Waterfalls:**
` + fence + `typescript
// ❌ Sequential fetches
const user  = await getUser(id);
const posts = await getPosts(user.id);

// ✅ Parallel fetches
const [user, posts] = await Promise.all([getUser(id), getPosts(id)]);
` + fence + bt(`

**Bundle Size:**
- Lazy-load every route: TICKconst Page = lazy(() => import('./Page'))TICK
- Audit with TICK@next/bundle-analyzerTICK before every release.
- No barrel re-exports that pull in unused modules.

---

## Security

**XSS Prevention:**
- Never use TICKdangerouslySetInnerHTMLTICK without DOMPurify.
- Set Content-Security-Policy: restrict TICKscript-srcTICK to TICK'self'TICK.

**JS Obfuscation (for sensitive isolated modules only):**
`) + fence + `bash
npm install --save-dev javascript-obfuscator

# Build step — NEVER obfuscate component trees
javascript-obfuscator src/sensitive/module.js \
  --output dist/sensitive/module.js \
  --compact true \
  --control-flow-flattening true \
  --string-array true
` + fence + bt(`

**Rules:** Obfuscate at build time only. Never obfuscate React components (breaks SSR + DevTools). Keep source maps private.

---

## Commands

`) + fence + `bash
npm run dev                         # Start dev server
npm run build                       # Production build
npm run lint                        # ESLint
npx @next/bundle-analyzer           # Analyze bundle (Next.js)
npx javascript-obfuscator <input> --output <output>  # Obfuscate module
` + fence + `

---

## Definition of Done

- [ ] No ` + "`useMemo`" + ` / ` + "`useCallback`" + ` added manually (React 19 Compiler handles it)
- [ ] Server Components used by default, ` + "`"+"use client"+"` only where required" + `
- [ ] No ` + "`dangerouslySetInnerHTML`" + ` without DOMPurify
- [ ] No ` + "`var()`" + ` in Tailwind ` + "`className`" + `, no hardcoded hex colors
- [ ] Zustand selectors use ` + "`useShallow`" + ` for multiple fields
- [ ] Bundle analyzed — no unexpected growth
- [ ] Sensitive modules obfuscated at build time, not in component tree
`
}

// ─── API Security ─────────────────────────────────────────────────────────────

func apiSecuritySkill() string {
	return skillFrontmatter(
		"api-security",
		"API Security agent. CORS allowlists, token-bucket rate-limiting, JWT issuance and validation.",
		"When configuring CORS, rate limits, JWT tokens, or API authentication.",
		"Read, Edit, Write, Glob, Grep",
	) + `# API Security Agent

## When to Use

- Configuring CORS for a new API
- Adding rate limiting to endpoints
- Implementing JWT authentication
- Setting security headers

---

## Decision Tree

` + fence + `
Need cross-origin requests?         → CORS allowlist — NEVER wildcard in production
Public endpoint getting hammered?   → Rate limiting with token-bucket strategy
Need to identify the caller?        → JWT with explicit algorithm validation
Storing sensitive tokens?           → httpOnly + Secure + SameSite=Strict cookie
` + fence + `

---

## CORS

` + fence + `typescript
const allowedOrigins = process.env.CORS_ORIGINS?.split(',') ?? [];

app.use(cors({
  origin: (origin, callback) => {
    if (!origin || allowedOrigins.includes(origin)) return callback(null, true);
    callback(new Error('Not allowed by CORS'));
  },
  methods: ['GET', 'POST', 'PUT', 'PATCH', 'DELETE'],
  allowedHeaders: ['Content-Type', 'Authorization', 'X-Trace-ID'],
  credentials: true,
}));
` + fence + bt(`

**Rules:** Origins from env var per environment. TICKcredentials: trueTICK only when cookies/auth headers required.

---

## Rate Limiting

`) + fence + `typescript
import rateLimit from 'express-rate-limit';

export const apiLimiter = rateLimit({
  windowMs: 60 * 1000,
  max: 100,
  standardHeaders: true,
  legacyHeaders: false,
  handler: (req, res) => res.status(429).json({
    success: false,
    code: 'RATE_LIMIT_EXCEEDED',
    message: 'Too many requests',
    detail: 'Retry after ' + req.rateLimit.resetTime?.toISOString(),
  }),
});
` + fence + `

| Endpoint type | Limit |
|---------------|-------|
| Public read | 200 req/min |
| Authenticated write | 60 req/min |
| Auth (login, refresh) | 10 req/min |
| Webhooks | 5 req/min |

---

## JWT

` + fence + `typescript
// Issuance
jwt.sign({ sub: userId, role, traceId }, process.env.JWT_SECRET, { expiresIn: '15m', algorithm: 'HS256' });

// Validation middleware
const verifyToken = (req, res, next) => {
  const token = req.headers.authorization?.split(' ')[1];
  if (!token) return res.status(401).json({ success: false, code: 'MISSING_TOKEN', message: 'Unauthorized', detail: '' });
  try {
    req.user = jwt.verify(token, process.env.JWT_SECRET, { algorithms: ['HS256'] });
    next();
  } catch (err) {
    return res.status(401).json({ success: false, code: 'INVALID_TOKEN', message: 'Token invalid or expired', detail: err.message });
  }
};
` + fence + bt(`

**Rules:** Access token TTL 15m. Refresh token TTL 7d in TICKhttpOnlyTICK+TICKSecureTICK+TICKSameSite=StrictTICK cookie. Rotate refresh tokens on every use. Always validate TICKalgTICK explicitly.

---

## Security Headers

`) + fence + `typescript
app.use((req, res, next) => {
  res.setHeader('Strict-Transport-Security', 'max-age=63072000; includeSubDomains');
  res.setHeader('X-Content-Type-Options', 'nosniff');
  res.setHeader('X-Frame-Options', 'DENY');
  res.setHeader('Referrer-Policy', 'strict-origin-when-cross-origin');
  next();
});
` + fence + `

---

## Commands

` + fence + `bash
npm install cors express-rate-limit jsonwebtoken
npm install -D @types/cors @types/jsonwebtoken
` + fence + `

## Definition of Done

- [ ] CORS uses explicit allowlist from env — no wildcards
- [ ] Rate limiting at router level, not per-handler
- [ ] JWT validates ` + "`alg`" + ` explicitly
- [ ] Refresh tokens in httpOnly cookies, rotated on use
- [ ] All 4 security headers set on every response
`
}

// ─── Controller ───────────────────────────────────────────────────────────────

func controllerSkill() string {
	return skillFrontmatter(
		"controller",
		"Controller agent. OpenAI-compatible endpoints, Zod v4 DTOs, TypeScript const types, standard error shape.",
		"When creating API endpoints, validating DTOs, or standardizing error responses.",
		"Read, Edit, Write, Glob, Grep",
	) + bt(`# Controller Agent

## When to Use

- Creating REST API endpoints
- Validating request DTOs
- Standardizing error responses
- Writing OpenAPI documentation

**Don't use when:** You need business logic — delegate to Service Layer.

> Assets: See [assets/ERROR_RESPONSE.ts](assets/ERROR_RESPONSE.ts) for the full error response template.

---

## Decision Tree

`) + fence + `
Request failed DTO validation?      → 400 VALIDATION_ERROR
No auth token?                      → 401 UNAUTHORIZED
Valid token, wrong role?            → 403 FORBIDDEN
Resource doesn't exist?             → 404 NOT_FOUND
Duplicate / stale version?          → 409 CONFLICT
Business rule violation?            → 422 UNPROCESSABLE (from DomainException)
Unexpected server error?            → 500 INTERNAL_ERROR
` + fence + bt(`

---

## TypeScript Const Types (REQUIRED for codes and enums)

`) + fence + `typescript
// ✅ ALWAYS: const object → extract type (single source of truth)
const ERROR_CODES = {
  VALIDATION_ERROR: 'VALIDATION_ERROR',
  NOT_FOUND:        'NOT_FOUND',
  CONFLICT:         'CONFLICT',
} as const;

type ErrorCode = (typeof ERROR_CODES)[keyof typeof ERROR_CODES];

// ❌ NEVER: direct union types
type ErrorCode = 'VALIDATION_ERROR' | 'NOT_FOUND' | 'CONFLICT';
` + fence + `

---

## Standard Error Shape

Every error response — without exception:

` + fence + `typescript
interface ErrorResponse {
  success: false;
  code: ErrorCode;   // from ERROR_CODES const
  message: string;   // safe to show in UI
  detail: string;    // technical context — never PII
}
` + fence + `

---

## DTO Validation — Zod v4

` + fence + `typescript
import { z } from "zod";

// ✅ Zod v4 — top-level validators
const CreateUserDTO = z.object({
  email: z.email({ error: "Invalid email address" }),   // NOT z.string().email()
  name:  z.string().min(2, { error: "Name too short" }).max(100),
  role:  z.enum(["admin", "user"]),
  id:    z.uuid().optional(),                           // NOT z.string().uuid()
});

// ❌ Zod v3 (OLD — do not use)
// z.string().email()
// z.string().uuid()
// z.string().nonempty()

export async function createUser(req, res) {
  const parsed = CreateUserDTO.safeParse(req.body);
  if (!parsed.success) {
    return res.status(400).json({
      success: false,
      code: ERROR_CODES.VALIDATION_ERROR,
      message: 'Invalid request payload',
      detail: parsed.error.issues.map(i => i.path.join('.') + ': ' + i.message).join('; '),
    });
  }
  const result = await userService.create(parsed.data, req.traceId);
  return res.status(201).json({ success: true, data: result });
}
` + fence + `

---

## Flat TypeScript Interfaces for DTOs

` + fence + `typescript
// ✅ One level deep — nested objects get their own interface
interface CreateAddressDTO {
  street: string;
  city: string;
}
interface CreateUserDTO {
  email: string;
  name: string;
  address: CreateAddressDTO;
}

// ❌ NEVER inline nested objects
interface CreateUserDTO {
  address: { street: string; city: string };  // NO
}
` + fence + `

---

## OpenAPI 3.0 Annotation (required per endpoint)

` + fence + `yaml
/users:
  post:
    summary: Create a new user
    operationId: createUser
    requestBody:
      required: true
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/CreateUserDTO'
    responses:
      '201':
        description: User created
      '400':
        $ref: '#/components/responses/ValidationError'
      '409':
        $ref: '#/components/responses/Conflict'
` + fence + bt(`

---

## Rules

1. Controllers call **services only** — never repositories directly.
2. Pass TICKtraceIdTICK from request header (or generate one) to every service call.
3. Return TICK{ success: true, data: ... }TICK on success.
4. Never leak stack traces in responses.
5. Rate limiting at router level — not inside handlers.

---

## Commands

`) + fence + `bash
npm install zod
npx ts-node scripts/generate-openapi.ts   # regenerate OpenAPI spec
` + fence + `

## Definition of Done

- [ ] Every endpoint has an OpenAPI 3.0 annotation
- [ ] DTOs use Zod v4 (` + "`z.email()`" + `, ` + "`z.uuid()`" + `, not ` + "`z.string().email()`" + `)
- [ ] Error codes use ` + "`as const`" + ` pattern — no raw string unions
- [ ] All error paths return standard ` + "`ErrorResponse`" + ` shape
- [ ] No service or repository logic inside controller
- [ ] ` + "`traceId`" + ` passed to every service call
`
}

// ─── Service Layer ────────────────────────────────────────────────────────────

func serviceLayerSkill() string {
	return skillFrontmatter(
		"service-layer",
		"Service Layer agent. Business logic, DI/IoC, TypeScript const types, domain exception handling.",
		"When implementing business logic, applying DI/IoC, or handling domain exceptions.",
		"Read, Edit, Write, Glob, Grep",
	) + bt(`# Service Layer Agent

## When to Use

- Implementing business rules and domain logic
- Wiring dependencies (DI/IoC)
- Throwing domain exceptions

**Don't use when:** You need SQL or HTTP — those belong in Repository and Controller.

---

## Decision Tree

`) + fence + `
Business rule violated?             → throw DomainException with typed code
Need a repository?                  → inject via constructor interface — never import directly
Need to call external service?      → inject as interface — mock in tests
Got a raw DB error?                 → should have been caught in Repository, not here
` + fence + `

---

## TypeScript Const Types for Domain Codes (REQUIRED)

` + fence + `typescript
// ✅ Single source of truth — runtime + type
const DOMAIN_CODES = {
  USER_ALREADY_EXISTS:  'USER_ALREADY_EXISTS',
  INSUFFICIENT_BALANCE: 'INSUFFICIENT_BALANCE',
  ACCOUNT_SUSPENDED:    'ACCOUNT_SUSPENDED',
} as const;

type DomainCode = (typeof DOMAIN_CODES)[keyof typeof DOMAIN_CODES];

// Use in exceptions
throw new DomainException(DOMAIN_CODES.USER_ALREADY_EXISTS, 'Email already registered', dto.email);
` + fence + `

---

## Dependency Injection / IoC

` + fence + `typescript
// Repository interface — service depends on abstraction
interface UserRepository {
  findByEmail(email: string, traceId: string): Promise<User | null>;
  save(user: User, traceId: string): Promise<User>;
}

class UserService {
  constructor(
    private readonly userRepo: UserRepository,  // injected
    private readonly logger: Logger,            // injected
  ) {}

  async register(dto: CreateUserDTO, traceId: string): Promise<User> {
    this.logger.info({ traceId, op: 'register', email: dto.email });

    const existing = await this.userRepo.findByEmail(dto.email, traceId);
    if (existing) throw new DomainException(DOMAIN_CODES.USER_ALREADY_EXISTS, 'Email already registered', dto.email);

    return this.userRepo.save(User.create(dto), traceId);
  }
}
` + fence + `

---

## Domain Exception

` + fence + `typescript
export class DomainException extends Error {
  constructor(
    public readonly code: DomainCode,
    public readonly message: string,
    public readonly detail: string = '',
  ) {
    super(message);
    this.name = 'DomainException';
  }
}
` + fence + bt(`

---

## Rules

1. **No SQL, no ORM calls** — delegate to repository interfaces.
2. **No TICKreqTICK, TICKresTICK, TICKheadersTICK** — services are HTTP-agnostic.
3. Every public method receives TICKtraceIdTICK and passes it to every log and repo call.
4. One service, one domain concept.
5. Throw TICKDomainExceptionTICK for business violations. Let unexpected errors bubble.
6. Use TICKas constTICK for all domain code enums — no raw string unions.

---

## Logger Convention

`) + fence + `typescript
this.logger.info({ traceId, op: 'methodName', ...context });
this.logger.warn({ traceId, op: 'methodName', code: DOMAIN_CODES.X, detail });
this.logger.error({ traceId, op: 'methodName', err: error.message });
` + fence + `

---

## Definition of Done

- [ ] Dependencies injected via constructor interfaces — never imported directly
- [ ] Zero SQL or ORM imports in service files
- [ ] Zero HTTP types imported
- [ ] Domain codes use ` + "`as const`" + ` pattern
- [ ] All business violations throw ` + "`DomainException`" + ` with typed code
- [ ] Every method logs entry with ` + "`traceId`" + `
`
}

// ─── Repository ───────────────────────────────────────────────────────────────

func repositorySkill() string {
	return skillFrontmatter(
		"repository",
		"Repository agent. SQL optimization, idempotency, transactions, DB exception mapping.",
		"When performing database operations, managing transactions, or optimizing SQL queries.",
		"Read, Edit, Write, Glob, Grep, Bash",
	) + `# Repository Agent

## When to Use

- Writing database queries
- Managing transactions
- Mapping DB errors to domain exceptions
- Optimizing slow queries

> Reference: github/awesome-copilot/sql-optimization

---

## Decision Tree

` + fence + `
Need to read multiple related rows?     → JOIN with explicit column list — no SELECT *
Updating multiple tables?               → Wrap in transaction with FOR UPDATE
Same operation could be called twice?   → Use ON CONFLICT DO NOTHING (upsert)
Got a DB error?                         → Map to DomainException before it leaves this layer
Large result set / paginated?           → Cursor-based (WHERE id > $cursor) — not OFFSET
` + fence + `

---

## SQL Rules

**No SELECT ***
` + fence + `sql
-- ❌ BAD
SELECT * FROM users WHERE id = $1;

-- ✅ GOOD
SELECT id, email, name, role, created_at FROM users WHERE id = $1;
` + fence + `

**No functions on indexed columns in WHERE**
` + fence + `sql
-- ❌ BAD — prevents index use
WHERE DATE(created_at) = '2024-01-01'

-- ✅ GOOD
WHERE created_at >= '2024-01-01' AND created_at < '2024-01-02'
` + fence + `

**Cursor pagination over OFFSET**
` + fence + `sql
-- ❌ BAD — O(n) cost
SELECT id, total FROM orders ORDER BY id LIMIT 20 OFFSET 10000;

-- ✅ GOOD — O(1) with index
SELECT id, total FROM orders WHERE id > $1 ORDER BY id LIMIT 20;
` + fence + `

**Batch inserts over loops**
` + fence + `sql
INSERT INTO events (user_id, type, payload)
SELECT * FROM UNNEST($1::uuid[], $2::text[], $3::jsonb[]);
` + fence + `

---

## Transactions

` + fence + `typescript
async transferFunds(fromId: string, toId: string, amount: number) {
  return this.db.transaction(async (trx) => {
    const from = await trx.query('SELECT balance FROM accounts WHERE id = $1 FOR UPDATE', [fromId]);
    if (from.rows[0].balance < amount)
      throw new DomainException('INSUFFICIENT_BALANCE', 'Transfer failed', 'from=' + fromId);
    await trx.query('UPDATE accounts SET balance = balance - $1 WHERE id = $2', [amount, fromId]);
    await trx.query('UPDATE accounts SET balance = balance + $1 WHERE id = $2', [amount, toId]);
  });
}
` + fence + `

---

## Idempotency

` + fence + `sql
INSERT INTO events (id, user_id, type, created_at)
VALUES ($1, $2, $3, NOW())
ON CONFLICT (id) DO NOTHING;
` + fence + `

---

## DB Exception Mapping

` + fence + `typescript
try {
  return await this.db.query(sql, params);
} catch (err) {
  if (err.code === '23505') throw new DomainException('DUPLICATE_ENTRY',       'Record already exists',        err.detail);
  if (err.code === '23503') throw new DomainException('FOREIGN_KEY_VIOLATION', 'Referenced record not found',  err.detail);
  if (err.code === '23502') throw new DomainException('NOT_NULL_VIOLATION',    'Required field missing',       err.detail);
  if (err.code === '40001') throw new DomainException('SERIALIZATION_FAILURE', 'Transaction conflict, retry',  err.detail);
  throw err;
}
` + fence + `

| PG Code | Domain Code |
|---------|-------------|
| ` + "`23505`" + ` | ` + "`DUPLICATE_ENTRY`" + ` |
| ` + "`23503`" + ` | ` + "`FOREIGN_KEY_VIOLATION`" + ` |
| ` + "`23502`" + ` | ` + "`NOT_NULL_VIOLATION`" + ` |
| ` + "`40001`" + ` | ` + "`SERIALIZATION_FAILURE`" + ` |

---

## Commands

` + fence + `bash
psql -d mydb -c "EXPLAIN ANALYZE <your query>"   # check query plan
psql -d mydb -c "\d <table>"                     # show table structure + indexes
` + fence + `

## Definition of Done

- [ ] No ` + "`SELECT *`" + ` anywhere
- [ ] No functions wrapping indexed columns in WHERE
- [ ] Pagination uses cursor, not OFFSET
- [ ] Multi-step writes wrapped in transactions with FOR UPDATE
- [ ] Upserts for idempotent operations
- [ ] All DB errors mapped to ` + "`DomainException`" + ` before leaving this layer
`
}

// ─── Migration ────────────────────────────────────────────────────────────────

func migrationSkill() string {
	return skillFrontmatter(
		"migration",
		"Migration agent. PostgreSQL schema design, indexing strategies, cardinality validation, reversible migrations.",
		"When creating DB migrations, designing schemas, or managing indexes.",
		"Read, Edit, Write, Glob, Grep, Bash",
	) + bt(`# Migration Agent

## When to Use

- Creating new tables or altering existing ones
- Adding indexes
- Designing PostgreSQL schemas

> Reference: wshobson/agents/postgresql-table-design
> Assets: See [assets/MIGRATION_TEMPLATE.sql](assets/MIGRATION_TEMPLATE.sql) for the standard template.

---

## Decision Tree

`) + fence + `
New table?                              → Use BIGINT GENERATED ALWAYS AS IDENTITY for PK
Need global / opaque ID?                → Add UUID external_id with DEFAULT gen_random_uuid()
Storing timestamps?                     → TIMESTAMPTZ — never plain TIMESTAMP
Storing money?                          → NUMERIC(p,s) — never FLOAT or MONEY
High-cardinality filter column?         → B-tree index
JSONB search?                           → GIN index
Time-series / append-only?             → BRIN index
Low-cardinality (boolean, status)?      → Partial index or skip
Adding FK column?                       → ALWAYS add index manually (PG doesn't auto-index FKs)
` + fence + bt(`

---

## Data Types

| Use case | Type | Avoid |
|----------|------|-------|
| Primary keys | TICKBIGINT GENERATED ALWAYS AS IDENTITYTICK | TICKSERIALTICK |
| Global / opaque IDs | TICKUUIDTICK | — |
| Timestamps | TICKTIMESTAMPTZTICK | TICKTIMESTAMPTICK |
| Money / precision | TICKNUMERIC(p, s)TICK | TICKFLOATTICK, TICKMONEYTICK |
| Strings | TICKTEXTTICK | TICKVARCHAR(n)TICK unless constraint is meaningful |
| JSON | TICKJSONBTICK | TICKJSONTICK |

---

## Schema Template

`) + fence + `sql
-- up.sql
CREATE TABLE IF NOT EXISTS users (
  id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  external_id UUID        NOT NULL DEFAULT gen_random_uuid() UNIQUE,
  email       TEXT        NOT NULL,
  name        TEXT        NOT NULL,
  role        TEXT        NOT NULL DEFAULT 'user',
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- REQUIRED: descriptions for every column
COMMENT ON COLUMN users.id          IS 'Internal surrogate key — never exposed externally';
COMMENT ON COLUMN users.external_id IS 'UUID exposed in API responses — immutable after creation';
COMMENT ON COLUMN users.email       IS 'Unique user email — normalized to lowercase on insert';
COMMENT ON COLUMN users.role        IS 'Access role: admin | user | viewer';

-- Indexes
CREATE UNIQUE INDEX idx_users_email       ON users (LOWER(email));
CREATE INDEX        idx_users_external_id ON users (external_id);
CREATE INDEX        idx_users_role        ON users (role) WHERE role != 'user'; -- partial index

-- down.sql (always present)
DROP TABLE IF EXISTS users;
` + fence + `

---

## Indexing Reference

| Cardinality | Type | Example |
|-------------|------|---------|
| High (email, UUID) | ` + "`B-tree unique`" + ` | ` + "`CREATE UNIQUE INDEX ON users (email)`" + ` |
| Medium filter (status) | ` + "`Partial B-tree`" + ` | ` + "`WHERE status = 'active'`" + ` |
| JSONB search | ` + "`GIN`" + ` | ` + "`USING GIN (payload)`" + ` |
| Time-series | ` + "`BRIN`" + ` | ` + "`USING BRIN (created_at)`" + ` |
| FK column | ` + "`B-tree`" + ` | Always — PG doesn't auto-index FKs |

---

## Migration Rules

1. Every ` + "`up`" + ` has a matching ` + "`down`" + `.
2. Use ` + "`IF NOT EXISTS`" + ` / ` + "`IF EXISTS`" + ` — migrations must be idempotent.
3. Never alter a column type directly — add column, backfill, drop old (3-step).
4. ` + "`COMMENT ON COLUMN`" + ` is required for every field — no undocumented columns.
5. Tables > 100M rows need a partitioning plan.

---

## Commands

` + fence + `bash
# Prisma
npx prisma migrate dev --name <migration_name>
npx prisma migrate status

# Knex
npx knex migrate:make <migration_name>
npx knex migrate:latest

# Check table indexes
psql -d mydb -c "\d <table_name>"
` + fence + `

## Definition of Done

- [ ] Every column has ` + "`COMMENT ON COLUMN`" + `
- [ ] Every FK column has an explicit index
- [ ] Index type matches cardinality
- [ ] ` + "`down`" + ` migration present and tested
- [ ] Migration is idempotent
- [ ] Tables > 100M rows have a partitioning plan
`
}

// ─── Observability ────────────────────────────────────────────────────────────

func observabilitySkill() string {
	return skillFrontmatter(
		"observability",
		"Observability agent. Structured logging, Trace-ID propagation, metrics, no PII in logs.",
		"When adding logging, propagating Trace-IDs, or collecting metrics.",
		"Read, Edit, Write, Glob, Grep",
	) + `# Observability Agent

## When to Use

- Adding logging to any layer
- Propagating Trace-IDs
- Exposing metrics
- Debugging production issues

---

## Decision Tree

` + fence + `
New request entering the system?       → Generate or forward X-Trace-ID header
Passing to downstream service?         → Include X-Trace-ID in call
Logging a state change?               → INFO with traceId, op, context
Logging a recoverable issue?          → WARN with traceId, code, detail
Logging a failure?                    → ERROR with traceId, op, file:line, code
About to log user data?               → Stop — use userId reference instead
` + fence + `

---

## Trace-ID Propagation

` + fence + `typescript
// Entry point middleware — generate or forward
app.use((req, res, next) => {
  req.traceId = req.headers['x-trace-id'] ?? crypto.randomUUID();
  res.setHeader('X-Trace-ID', req.traceId);
  next();
});
` + fence + `

Propagation chain:
` + fence + `
HTTP Request
  → Controller   (req.traceId)
  → Service      (traceId param)
  → Repository   (traceId param)
  → Logger       (traceId field in every entry)
  → External API (X-Trace-ID header)
` + fence + `

---

## Structured Log Format

` + fence + `typescript
interface LogEntry {
  traceId:     string;
  level:       'debug' | 'info' | 'warn' | 'error';
  op:          string;   // createUser, fetchOrder
  message:     string;
  code?:       string;   // error code
  detail?:     string;   // technical context — never PII
  durationMs?: number;
}

logger.info({ traceId, op: 'createUser', message: 'User registered', userId });
logger.error({ traceId, op: 'transfer', message: 'Transaction failed', code: 'SERIALIZATION_FAILURE', detail: err.message });
` + fence + `

**Levels:**

| Level | When |
|-------|------|
| ` + "`debug`" + ` | Internal state. Off in production. |
| ` + "`info`" + `  | State transitions: user created, order placed. |
| ` + "`warn`" + `  | Recoverable: retry, near rate-limit threshold. |
| ` + "`error`" + ` | Failures needing investigation. |

---

## Error Log Format

` + fence + `
[abc-123] ERROR transferFunds (src/services/account.service.ts:42) — SERIALIZATION_FAILURE: retry transaction
` + fence + `

Format: ` + "`[traceId] LEVEL op (file:line) — CODE: message`" + `

---

## Never Log

- Passwords, tokens, secrets
- Full credit card / PAN
- PII (email body, phone, SSN) — use ` + "`userId`" + ` reference only
- Full request bodies in production

---

## Metrics to Expose

| Metric | Type | Labels |
|--------|------|--------|
| ` + "`http_requests_total`" + ` | Counter | ` + "`method, route, status`" + ` |
| ` + "`http_request_duration_ms`" + ` | Histogram | ` + "`method, route`" + ` |
| ` + "`db_query_duration_ms`" + ` | Histogram | ` + "`operation`" + ` |
| ` + "`domain_errors_total`" + ` | Counter | ` + "`code`" + ` |

---

## Commands

` + fence + `bash
npm install pino pino-pretty        # Structured logger
npm install prom-client             # Prometheus metrics
` + fence + `

## Definition of Done

- [ ] Every request has Trace-ID from HTTP entry to DB query
- [ ] All logs are structured JSON in production
- [ ] No PII in any log line
- [ ] Error logs include ` + "`code`" + `, ` + "`op`" + `, ` + "`file:line`" + `
- [ ] Minimum 4 metrics exposed on ` + "`/metrics`" + `
`
}

// ─── QA ───────────────────────────────────────────────────────────────────────

func qaSkill() string {
	return skillFrontmatter(
		"qa",
		"QA agent. Unit tests with parametrize patterns, Playwright E2E POM, minimal mocking, Trace-ID in error paths.",
		"When writing unit tests, E2E tests, covering edge cases, or reviewing test coverage.",
		"Read, Edit, Write, Glob, Grep, Bash",
	) + bt(`# QA Agent

## When to Use

- Writing unit tests for services, controllers, repositories
- Writing E2E tests with Playwright
- Reviewing test coverage
- Choosing what to mock

**Don't use when:** You're writing integration tests against real infrastructure — there's no special skill for that, just use the real thing.

---

## Decision Tree

`) + fence + `
Testing business logic?             → Unit test — mock only I/O interfaces
Testing DB queries?                 → Integration test — real DB (test schema)
Testing user flows?                 → Playwright E2E — use MCP workflow if available
Need to mock something internal?    → Stop — you're testing the wrong thing
Multiple similar edge cases?        → Use test.each / parametrize — one test block
` + fence + bt(`

---

## Test Naming Convention

`) + fence + `
Test<Unit>_<Scenario>_<ExpectedOutcome>

TestCreateUser_DuplicateEmail_ReturnsDomainException
TestTransfer_InsufficientBalance_ThrowsError
TestGetOrder_NotFound_Returns404
TestCreateUser_ValidInput_ReturnsCreatedUser
` + fence + bt(`

---

## Parametrize / TICKtest.eachTICK (REQUIRED for edge cases)

Instead of repeating similar tests, use table-driven tests:

`) + fence + `typescript
// ✅ test.each — one block covers all edge cases
test.each([
  { email: '',              expected: 'VALIDATION_ERROR' },
  { email: 'not-an-email',  expected: 'VALIDATION_ERROR' },
  { email: 'a@b.c',         expected: 'VALIDATION_ERROR' },  // too short domain
  { email: 'valid@test.com',expected: null },
])('validateEmail($email) → $expected', ({ email, expected }) => {
  const result = validateEmail(email);
  expect(result?.code ?? null).toBe(expected);
});

// Python equivalent
@pytest.mark.parametrize("email,expected", [
  ("",              "VALIDATION_ERROR"),
  ("not-an-email",  "VALIDATION_ERROR"),
  ("valid@test.com", None),
])
def test_validate_email(email, expected):
    result = validate_email(email)
    assert (result.code if result else None) == expected
` + fence + `

---

## Minimal Mocking Philosophy

` + fence + `typescript
// ✅ GOOD — mock only external I/O at the interface boundary
const mockRepo = { findByEmail: jest.fn(), save: jest.fn() };
const service  = new UserService(mockRepo, noopLogger);

// ❌ BAD — mocking your own services defeats the purpose
jest.mock('../services/user.service');
` + fence + bt(`

**Mock only:** DB/repository interfaces, external HTTP, file system, TICKDate.now()TICK / TICKcrypto.randomUUID()TICK.
**Never mock:** Your own business logic, pure utility functions, validation.

---

## Playwright — E2E Tests

### MCP Workflow (use if Playwright MCP is available)

1. Navigate to target page.
2. Take snapshot — see real DOM structure and element labels.
3. Interact to verify the exact user flow.
4. Document actual selectors from snapshot (not assumed).
5. Only then write the test code.

### Selector Priority

`) + fence + `typescript
// 1. Best — role-based (accessible)
page.getByRole("button", { name: "Submit" });
page.getByLabel("Email");

// 2. Acceptable — text content
page.getByText("Invalid credentials");

// 3. Last resort — test ID
page.getByTestId("date-picker");

// ❌ Fragile — never use
page.locator(".btn-primary");
page.locator("#email");
` + fence + `

### Page Object Pattern

` + fence + `typescript
export class BasePage {
  constructor(protected page: Page) {}

  async goto(path: string) {
    await this.page.goto(path);
    await this.page.waitForLoadState("networkidle");
  }
}

export class LoginPage extends BasePage {
  readonly emailInput    = this.page.getByLabel("Email");
  readonly passwordInput = this.page.getByLabel("Password");
  readonly submitButton  = this.page.getByRole("button", { name: "Sign in" });

  async goto()  { await super.goto("/login"); }
  async login(email: string, password: string) {
    await this.emailInput.fill(email);
    await this.passwordInput.fill(password);
    await this.submitButton.click();
  }
}

// Test
test("User can login", { tag: ["@critical", "@e2e", "@LOGIN-E2E-001"] }, async ({ page }) => {
  const loginPage = new LoginPage(page);
  await loginPage.goto();
  await loginPage.login("user@test.com", "pass123");
  await expect(page).toHaveURL("/dashboard");
});
` + fence + `

---

## Test Structure

` + fence + `
tests/
├── unit/          # Pure logic — no I/O
│   ├── services/
│   └── domain/
├── integration/   # Real DB (test schema) — no HTTP mocks
│   └── repositories/
└── e2e/           # Full HTTP stack via Playwright
    └── login/
        ├── login-page.ts
        ├── login.spec.ts
        └── login.md
` + fence + `

---

## Coverage Targets

| Layer | Minimum |
|-------|---------|
| Service Layer | 85% branch |
| Repository | 70% |
| Controller | 80% |
| Domain models | 90% |

---

## Commands

` + fence + `bash
npx jest --coverage                  # unit tests with coverage
npx jest --watch                     # watch mode
npx playwright test                  # E2E all
npx playwright test --ui             # interactive UI
npx playwright test --debug          # debug mode
npx playwright test tests/login/     # specific folder
npx playwright test --grep "@critical"  # filter by tag
pytest -v                            # Python — verbose
pytest -m "not slow"                 # Python — skip slow
pytest --cov=src                     # Python — coverage
` + fence + `

## Definition of Done

- [ ] Every public function: ≥1 success test, ≥2 failure tests
- [ ] Edge cases use ` + "`test.each`" + ` / ` + "`parametrize`" + ` — no copy-paste test blocks
- [ ] Only external I/O boundaries mocked
- [ ] E2E tests use Page Object Model
- [ ] Playwright selectors use ` + "`getByRole`" + ` / ` + "`getByLabel`" + ` first
- [ ] Coverage meets layer targets
- [ ] Error paths verified to include ` + "`traceId`" + ` in logs
`
}

// ─── Skill Creator ────────────────────────────────────────────────────────────

func skillCreatorSkill() string {
	return skillFrontmatter(
		"skill-creator",
		"Creates new agent skills following the mr-agent-ai spec. Maintains consistency across all skills.",
		"When asked to create a new skill, add agent instructions, or document patterns for AI.",
		"Read, Edit, Write, Glob, Grep, Bash, WebFetch, WebSearch",
	) + bt(`# Skill Creator Agent

## When to Use

Create a new skill when:
- A pattern is used repeatedly and the AI needs explicit guidance
- Project conventions differ from generic best practices
- A complex workflow needs step-by-step instructions
- A decision tree would reduce ambiguity

**Don't create a skill when:**
- Documentation already exists — create a TICK`+"`"+`references/`+"`"+`TICK entry instead
- The pattern is trivial or self-explanatory
- It's a one-off task

---

## Skill Structure

`) + fence + `
skills/{skill-name}/
├── SKILL.md              # Required — main skill file
├── assets/               # Optional — templates, schemas, examples
│   ├── TEMPLATE.ts
│   └── schema.json
└── references/           # Optional — links to local docs
    └── docs.md
` + fence + bt(`

TICKassets/TICK = code templates, JSON schemas, example configs.
TICKreferences/TICK = pointers to LOCAL files (not web URLs).

---

## Decision Tree

`) + fence + `
Pattern reused 3+ times?            → Create a skill
Project-specific convention?        → Create a prowler-{name} or project-{name} skill
Need code templates?                → Add to assets/
Linking to existing docs?           → Use references/ with local path
Already exists in skills/?          → Update it — don't duplicate
` + fence + bt(`

---

## Naming Conventions

| Type | Pattern | Examples |
|------|---------|----------|
| Generic skill | TICK{technology}TICK | TICKtypescriptTICK, TICKplaywriterTICK |
| Project-specific | TICK{project}-{component}TICK | TICKmy-app-authTICK |
| Workflow skill | TICK{action}-{target}TICK | TICKskill-creatorTICK, TICKgithub-prTICK |
| Testing skill | TICKtest-{component}TICK | TICKtest-servicesTICK |

---

## Required Frontmatter Fields

| Field | Required | Value |
|-------|----------|-------|
| TICKnameTICK | Yes | Lowercase, hyphens |
| TICKdescriptionTICK | Yes | What + Trigger in one block |
| TICKlicenseTICK | Yes | TICKMITTICK |
| TICKallowed-toolsTICK | Yes | Comma-separated tool list |
| TICKmetadata.versionTICK | Yes | Semantic version as string |

---

## Content Guidelines

**DO:**
- Start with the most critical patterns
- Use decision trees for ambiguous choices
- Include ✅ / ❌ code examples
- Add TICKCommands TICK section with copy-paste bash

**DON'T:**
- Duplicate content from existing skills — reference instead
- Add lengthy prose explanations — bullet points and code
- Include troubleshooting sections
- Use web URLs in TICKreferences/TICK — use local paths

---

## After Creating a Skill

Add it to TICKAGENTS.mdTICK:

`) + fence + `markdown
| `+"`"+`mr-agent-{skill-name}`+"`"+` | {trigger description} | [`+"`"+`skills/{skill-name}/SKILL.md`+"`"+`](skills/{skill-name}/SKILL.md) |
` + fence + bt(`

---

## Checklist Before Creating

- [ ] Skill doesn't already exist (checked TICKskills/TICK with Glob)
- [ ] Pattern is reusable — not a one-off
- [ ] Name follows conventions
- [ ] Frontmatter includes TICKallowed-toolsTICK
- [ ] Decision tree covers the most common branching point
- [ ] ✅ / ❌ code examples are minimal and focused
- [ ] TICKCommands TICK section is present
- [ ] TICKDefinition of DoneTICK checklist present
- [ ] Added to TICKAGENTS.mdTICK

---

## Commands

`) + fence + `bash
# List existing skills
ls skills/

# Check if skill already exists
ls skills/ | grep <skill-name>
` + fence + `

## Resources

- **Template**: See [assets/SKILL-TEMPLATE.md](assets/SKILL-TEMPLATE.md) for the full template
`
}

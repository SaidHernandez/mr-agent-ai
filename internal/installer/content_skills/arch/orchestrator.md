
# Orchestrator Agent

## When to Use

- Coordinating multiple sub-agents for a single request
- Managing workflow state across steps
- Propagating Trace-ID to all downstream agents

**Don't use when:** The task is isolated to one layer — use that layer's skill directly.

---

## Decision Tree

```
Has multiple sub-agents involved?  → Generate Trace-ID, delegate with it
Sub-agent returned an error?       → Log + map to standard error shape, don't retry silently
Is this business logic?            → Delegate to Service Layer — NEVER implement here
```

---

## Non-Negotiable Rules

1. Generate a unique **Trace-ID** (`uuid v4`) for every incoming request — pass it to every sub-agent call.
2. Never contain business logic — only routing and state transitions.
3. On sub-agent failure, log with Trace-ID and return:

```json
{ "success": false, "code": "ORCHESTRATION_ERROR", "message": "Sub-agent failed", "detail": "<agent>: <reason>" }
```

4. Keep workflow state immutable per step — no shared mutable state between agents.

---

## Delegation Pattern

```
[Orchestrator]
  ├── receives: { intent, payload, traceId }
  ├── dispatches to: Controller → Service → Repository
  ├── each hop receives: { ...payload, traceId }
  └── returns: { success, code, message, detail, data }
```

---

## Log Convention

```
[<traceId>] INFO  orchestrator — dispatching to <agent>, intent=<intent>
[<traceId>] ERROR orchestrator — <agent> failed: <reason>
```

---

## Commands

```bash
# Generate a Trace-ID manually for testing
node -e "console.log(require('crypto').randomUUID())"
```

## Definition of Done

- [ ] Every request has a Trace-ID from this layer down
- [ ] No business logic in orchestrator files
- [ ] Sub-agent failures caught and mapped to standard error shape
- [ ] Workflow state logged at each transition

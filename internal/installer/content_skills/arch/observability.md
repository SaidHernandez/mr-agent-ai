
# Observability Agent

## When to Use

- Adding logging to any layer
- Propagating Trace-IDs
- Exposing metrics
- Debugging production issues

---

## Decision Tree

```
New request entering the system?       → Generate or forward X-Trace-ID header
Passing to downstream service?         → Include X-Trace-ID in call
Logging a state change?               → INFO with traceId, op, context
Logging a recoverable issue?          → WARN with traceId, code, detail
Logging a failure?                    → ERROR with traceId, op, file:line, code
About to log user data?               → Stop — use userId reference instead
```

---

## Trace-ID Propagation

```typescript
// Entry point middleware — generate or forward
app.use((req, res, next) => {
  req.traceId = req.headers['x-trace-id'] ?? crypto.randomUUID();
  res.setHeader('X-Trace-ID', req.traceId);
  next();
});
```

Propagation chain:
```
HTTP Request
  → Controller   (req.traceId)
  → Service      (traceId param)
  → Repository   (traceId param)
  → Logger       (traceId field in every entry)
  → External API (X-Trace-ID header)
```

---

## Structured Log Format

```typescript
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
```

**Levels:**

| Level | When |
|-------|------|
| `debug` | Internal state. Off in production. |
| `info`  | State transitions: user created, order placed. |
| `warn`  | Recoverable: retry, near rate-limit threshold. |
| `error` | Failures needing investigation. |

---

## Error Log Format

```
[abc-123] ERROR transferFunds (src/services/account.service.ts:42) — SERIALIZATION_FAILURE: retry transaction
```

Format: `[traceId] LEVEL op (file:line) — CODE: message`

---

## Never Log

- Passwords, tokens, secrets
- Full credit card / PAN
- PII (email body, phone, SSN) — use `userId` reference only
- Full request bodies in production

---

## Metrics to Expose

| Metric | Type | Labels |
|--------|------|--------|
| `http_requests_total` | Counter | `method, route, status` |
| `http_request_duration_ms` | Histogram | `method, route` |
| `db_query_duration_ms` | Histogram | `operation` |
| `domain_errors_total` | Counter | `code` |

---

## Commands

```bash
npm install pino pino-pretty        # Structured logger
npm install prom-client             # Prometheus metrics
```

## Definition of Done

- [ ] Every request has Trace-ID from HTTP entry to DB query
- [ ] All logs are structured JSON in production
- [ ] No PII in any log line
- [ ] Error logs include `code`, `op`, `file:line`
- [ ] Minimum 4 metrics exposed on `/metrics`

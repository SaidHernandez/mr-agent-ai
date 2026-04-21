
# API Security Agent

## When to Use

- Configuring CORS for a new API
- Adding rate limiting to endpoints
- Implementing JWT authentication
- Setting security headers

---

## Decision Tree

```
Need cross-origin requests?         → CORS allowlist — NEVER wildcard in production
Public endpoint getting hammered?   → Rate limiting with token-bucket strategy
Need to identify the caller?        → JWT with explicit algorithm validation
Storing sensitive tokens?           → httpOnly + Secure + SameSite=Strict cookie
```

---

## CORS

```typescript
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
```

**Rules:** Origins from env var per environment. `credentials: true` only when cookies/auth headers required.

---

## Rate Limiting

```typescript
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
```

| Endpoint type | Limit |
|---------------|-------|
| Public read | 200 req/min |
| Authenticated write | 60 req/min |
| Auth (login, refresh) | 10 req/min |
| Webhooks | 5 req/min |

---

## JWT

```typescript
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
```

**Rules:** Access token TTL 15m. Refresh token TTL 7d in `httpOnly`+`Secure`+`SameSite=Strict` cookie. Rotate refresh tokens on every use. Always validate `alg` explicitly.

---

## Security Headers

```typescript
app.use((req, res, next) => {
  res.setHeader('Strict-Transport-Security', 'max-age=63072000; includeSubDomains');
  res.setHeader('X-Content-Type-Options', 'nosniff');
  res.setHeader('X-Frame-Options', 'DENY');
  res.setHeader('Referrer-Policy', 'strict-origin-when-cross-origin');
  next();
});
```

---

## Commands

```bash
npm install cors express-rate-limit jsonwebtoken
npm install -D @types/cors @types/jsonwebtoken
```

## Definition of Done

- [ ] CORS uses explicit allowlist from env — no wildcards
- [ ] Rate limiting at router level, not per-handler
- [ ] JWT validates `alg` explicitly
- [ ] Refresh tokens in httpOnly cookies, rotated on use
- [ ] All 4 security headers set on every response

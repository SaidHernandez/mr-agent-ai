
# Controller Agent

## When to Use

- Creating REST API endpoints
- Validating request DTOs
- Standardizing error responses
- Writing OpenAPI documentation

**Don't use when:** You need business logic — delegate to Service Layer.

> Assets: See [assets/ERROR_RESPONSE.ts](assets/ERROR_RESPONSE.ts) for the full error response template.

---

## Decision Tree

```
Request failed DTO validation?      → 400 VALIDATION_ERROR
No auth token?                      → 401 UNAUTHORIZED
Valid token, wrong role?            → 403 FORBIDDEN
Resource doesn't exist?             → 404 NOT_FOUND
Duplicate / stale version?          → 409 CONFLICT
Business rule violation?            → 422 UNPROCESSABLE (from DomainException)
Unexpected server error?            → 500 INTERNAL_ERROR
```

---

## TypeScript Const Types (REQUIRED for codes and enums)

```typescript
// ✅ ALWAYS: const object → extract type (single source of truth)
const ERROR_CODES = {
  VALIDATION_ERROR: 'VALIDATION_ERROR',
  NOT_FOUND:        'NOT_FOUND',
  CONFLICT:         'CONFLICT',
} as const;

type ErrorCode = (typeof ERROR_CODES)[keyof typeof ERROR_CODES];

// ❌ NEVER: direct union types
type ErrorCode = 'VALIDATION_ERROR' | 'NOT_FOUND' | 'CONFLICT';
```

---

## Standard Error Shape

Every error response — without exception:

```typescript
interface ErrorResponse {
  success: false;
  code: ErrorCode;   // from ERROR_CODES const
  message: string;   // safe to show in UI
  detail: string;    // technical context — never PII
}
```

---

## DTO Validation — Zod v4

```typescript
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
```

---

## Flat TypeScript Interfaces for DTOs

```typescript
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
```

---

## OpenAPI 3.0 Annotation (required per endpoint)

```yaml
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
```

---

## Rules

1. Controllers call **services only** — never repositories directly.
2. Pass `traceId` from request header (or generate one) to every service call.
3. Return `{ success: true, data: ... }` on success.
4. Never leak stack traces in responses.
5. Rate limiting at router level — not inside handlers.

---

## Commands

```bash
npm install zod
npx ts-node scripts/generate-openapi.ts   # regenerate OpenAPI spec
```

## Definition of Done

- [ ] Every endpoint has an OpenAPI 3.0 annotation
- [ ] DTOs use Zod v4 (`z.email()`, `z.uuid()`, not `z.string().email()`)
- [ ] Error codes use `as const` pattern — no raw string unions
- [ ] All error paths return standard `ErrorResponse` shape
- [ ] No service or repository logic inside controller
- [ ] `traceId` passed to every service call

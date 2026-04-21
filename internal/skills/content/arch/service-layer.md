
# Service Layer Agent

## When to Use

- Implementing business rules and domain logic
- Wiring dependencies (DI/IoC)
- Throwing domain exceptions

**Don't use when:** You need SQL or HTTP — those belong in Repository and Controller.

---

## Decision Tree

```
Business rule violated?             → throw DomainException with typed code
Need a repository?                  → inject via constructor interface — never import directly
Need to call external service?      → inject as interface — mock in tests
Got a raw DB error?                 → should have been caught in Repository, not here
```

---

## TypeScript Const Types for Domain Codes (REQUIRED)

```typescript
// ✅ Single source of truth — runtime + type
const DOMAIN_CODES = {
  USER_ALREADY_EXISTS:  'USER_ALREADY_EXISTS',
  INSUFFICIENT_BALANCE: 'INSUFFICIENT_BALANCE',
  ACCOUNT_SUSPENDED:    'ACCOUNT_SUSPENDED',
} as const;

type DomainCode = (typeof DOMAIN_CODES)[keyof typeof DOMAIN_CODES];

// Use in exceptions
throw new DomainException(DOMAIN_CODES.USER_ALREADY_EXISTS, 'Email already registered', dto.email);
```

---

## Dependency Injection / IoC

```typescript
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
```

---

## Domain Exception

```typescript
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
```

---

## Rules

1. **No SQL, no ORM calls** — delegate to repository interfaces.
2. **No `req`, `res`, `headers`** — services are HTTP-agnostic.
3. Every public method receives `traceId` and passes it to every log and repo call.
4. One service, one domain concept.
5. Throw `DomainException` for business violations. Let unexpected errors bubble.
6. Use `as const` for all domain code enums — no raw string unions.

---

## Logger Convention

```typescript
this.logger.info({ traceId, op: 'methodName', ...context });
this.logger.warn({ traceId, op: 'methodName', code: DOMAIN_CODES.X, detail });
this.logger.error({ traceId, op: 'methodName', err: error.message });
```

---

## Definition of Done

- [ ] Dependencies injected via constructor interfaces — never imported directly
- [ ] Zero SQL or ORM imports in service files
- [ ] Zero HTTP types imported
- [ ] Domain codes use `as const` pattern
- [ ] All business violations throw `DomainException` with typed code
- [ ] Every method logs entry with `traceId`

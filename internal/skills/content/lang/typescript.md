
# TypeScript Best Practices


## When to Use

- Writing or reviewing any TypeScript file
- Configuring a new TypeScript project (tsconfig.json)
- Refactoring JavaScript to TypeScript
- Code review of TS patterns

**Don't use when:** Working in plain JavaScript files with no TS compilation.

---

## Decision Tree

```text
Need a type for a plain object shape?  → interface
Need a union, tuple, or mapped type?   → type alias
Value that should never change?        → as const
Unsure of incoming type?               → unknown (never any)
Async operation?                       → async/await + try/catch
```

---

## Critical Patterns

### 1. Enable Strict Mode

```json
// ✅ tsconfig.json
{
  "compilerOptions": {
    "strict": true,
    "noImplicitAny": true,
    "strictNullChecks": true,
    "strictFunctionTypes": true
  }
}

// ❌ BAD — missing strict mode allows silent type errors
{
  "compilerOptions": {}
}
```

### 2. Never Use `any` — Prefer `unknown`

```typescript
// ✅ GOOD
function parse(input: unknown): string {
  if (typeof input !== "string") throw new Error("Expected string");
  return input.trim();
}

// ❌ BAD — any disables all type checking
function parse(input: any): string {
  return input.trim(); // runtime error if input is a number
}
```

### 3. Interfaces vs Type Aliases

```typescript
// ✅ GOOD — interface for extensible object shapes
interface User {
  id: string;
  name: string;
}
interface AdminUser extends User {
  role: "admin";
}

// ✅ GOOD — type for unions and mapped types
type Status = "active" | "inactive" | "pending";
type Nullable<T> = T | null;

// ❌ BAD — type alias where interface is clearer for object shapes
type User = { id: string; name: string };
```

### 4. Use `as const` for Literal Types

```typescript
// ✅ GOOD
const DIRECTIONS = ["north", "south", "east", "west"] as const;
type Direction = (typeof DIRECTIONS)[number]; // "north" | "south" | "east" | "west"

const HTTP_STATUS = { OK: 200, NOT_FOUND: 404 } as const;
type StatusCode = (typeof HTTP_STATUS)[keyof typeof HTTP_STATUS];

// ❌ BAD — inferred as string[], losing literal precision
const DIRECTIONS = ["north", "south", "east", "west"];
```

### 5. Explicit Return Types on Public APIs

```typescript
// ✅ GOOD — contract is explicit, catches return mismatches early
function getUserById(id: string): Promise<User | null> {
  return db.users.findOne({ id });
}

// ❌ BAD — return type inferred, breaks silently when DB changes
function getUserById(id: string) {
  return db.users.findOne({ id });
}
```

### 6. Type Guards for Safe Narrowing

```typescript
// ✅ GOOD
function isUser(value: unknown): value is User {
  return (
    typeof value === "object" &&
    value !== null &&
    "id" in value &&
    "name" in value
  );
}

if (isUser(payload)) {
  console.log(payload.name); // safe
}

// ❌ BAD — casting without checking
const user = payload as User;
console.log(user.name); // may blow up at runtime
```

### 7. Async/Await Error Handling & Parallelism

```typescript
// ✅ GOOD — wrap async in try/catch
async function fetchUser(id: string): Promise<User> {
  try {
    const res = await api.get<User>("/users/" + id);
    if (!res.ok) throw new Error("HTTP " + res.status);
    return await res.json();
  } catch (err) {
    logger.error("fetchUser failed", { id, err });
    throw err;
  }
}

// ✅ GOOD — run independent calls in parallel
const [user, posts] = await Promise.all([
  fetchUser(id),
  fetchPosts(id),
]);

// ❌ BAD — sequential await when operations are independent
const user  = await fetchUser(id);
const posts = await fetchPosts(id);
```

### 8. Type-Only Imports

```typescript
// ✅ GOOD — removed at compile time, reduces bundle size
import type { User, Post } from "./models";

// ❌ BAD — may include runtime code in the bundle
import { User, Post } from "./models";
```

### 9. Dependency Injection for Testability

```typescript
// ✅ GOOD
class UserService {
  constructor(private readonly repo: UserRepository) {}
  async getUser(id: string) { return this.repo.findById(id); }
}

// ❌ BAD — hard-coded dependency, impossible to unit test
class UserService {
  private repo = new PostgresUserRepository();
}
```

### 10. Optional Chaining & Nullish Coalescing

```typescript
// ✅ GOOD
const city = user?.address?.city ?? "Unknown";

// ❌ BAD — verbose and error-prone
const city = user && user.address && user.address.city
  ? user.address.city
  : "Unknown";
```

---

## Commands

```bash
tsc --noEmit          # type-check without emitting files
npx tsc --init        # generate tsconfig.json
```

---

## Definition of Done

- [ ] `strict: true` enabled in tsconfig.json
- [ ] No `any` usage — replaced with `unknown` or specific types
- [ ] Public functions have explicit return types
- [ ] All async functions have try/catch error handling
- [ ] Type guards used when narrowing union types
- [ ] `import type` used for type-only imports
- [ ] `as const` used for literal constant objects/arrays

---


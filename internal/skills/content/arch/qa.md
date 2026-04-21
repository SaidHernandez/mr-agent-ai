
# QA Agent

## When to Use

- Writing unit tests for services, controllers, repositories
- Writing E2E tests with Playwright
- Reviewing test coverage
- Choosing what to mock

**Don't use when:** You're writing integration tests against real infrastructure — there's no special skill for that, just use the real thing.

---

## Decision Tree

```
Testing business logic?             → Unit test — mock only I/O interfaces
Testing DB queries?                 → Integration test — real DB (test schema)
Testing user flows?                 → Playwright E2E — use MCP workflow if available
Need to mock something internal?    → Stop — you're testing the wrong thing
Multiple similar edge cases?        → Use test.each / parametrize — one test block
```

---

## Test Naming Convention

```
Test<Unit>_<Scenario>_<ExpectedOutcome>

TestCreateUser_DuplicateEmail_ReturnsDomainException
TestTransfer_InsufficientBalance_ThrowsError
TestGetOrder_NotFound_Returns404
TestCreateUser_ValidInput_ReturnsCreatedUser
```

---

## Parametrize / `test.each` (REQUIRED for edge cases)

Instead of repeating similar tests, use table-driven tests:

```typescript
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
```

---

## Minimal Mocking Philosophy

```typescript
// ✅ GOOD — mock only external I/O at the interface boundary
const mockRepo = { findByEmail: jest.fn(), save: jest.fn() };
const service  = new UserService(mockRepo, noopLogger);

// ❌ BAD — mocking your own services defeats the purpose
jest.mock('../services/user.service');
```

**Mock only:** DB/repository interfaces, external HTTP, file system, `Date.now()` / `crypto.randomUUID()`.
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

```typescript
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
```

### Page Object Pattern

```typescript
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
```

---

## Test Structure

```
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
```

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

```bash
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
```

## Definition of Done

- [ ] Every public function: ≥1 success test, ≥2 failure tests
- [ ] Edge cases use `test.each` / `parametrize` — no copy-paste test blocks
- [ ] Only external I/O boundaries mocked
- [ ] E2E tests use Page Object Model
- [ ] Playwright selectors use `getByRole` / `getByLabel` first
- [ ] Coverage meets layer targets
- [ ] Error paths verified to include `traceId` in logs

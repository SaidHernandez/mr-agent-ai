
# Python Best Practices


## When to Use

- Writing or reviewing any Python file
- Designing Python modules or packages
- Optimizing Python code for performance
- Code review of Pythonic patterns

**Don't use when:** Working in frameworks with their own enforced conventions (e.g., auto-generated Django migrations).

---

## Decision Tree

```text
Need to represent structured data?     → dataclass or Pydantic model
Opening a file or connection?          → context manager (with)
Building a list from iteration?        → list comprehension
Building a dict from iteration?        → dict comprehension
Need a default mutable argument?       → use None + default inside body
Catching exceptions?                   → catch specific exceptions, not bare except
Logging output?                        → logging module, not print()
```

---

## Critical Patterns

### 1. Always Use Type Hints

```python
# ✅ GOOD
def get_user(user_id: str) -> dict[str, str] | None:
    return db.find(user_id)

# ✅ GOOD — structured data with dataclass
from dataclasses import dataclass

@dataclass
class User:
    id: str
    name: str
    email: str

# ❌ BAD — no hints, intent is impossible to infer
def get_user(user_id):
    return db.find(user_id)
```

### 2. Never Use Mutable Default Arguments

```python
# ✅ GOOD
def append_item(item: str, items: list[str] | None = None) -> list[str]:
    if items is None:
        items = []
    items.append(item)
    return items

# ❌ BAD — default list is shared across ALL calls
def append_item(item: str, items: list[str] = []) -> list[str]:
    items.append(item)  # mutates the shared default object!
    return items
```

### 3. Use Context Managers for Resources

```python
# ✅ GOOD — file is closed even if an exception occurs
with open("data.csv", "r", encoding="utf-8") as f:
    data = f.read()

# ✅ GOOD — custom context manager
from contextlib import contextmanager

@contextmanager
def db_transaction(conn):
    try:
        yield conn
        conn.commit()
    except Exception:
        conn.rollback()
        raise

# ❌ BAD — resource leak if exception occurs before close()
f = open("data.csv")
data = f.read()
f.close()
```

### 4. Prefer List/Dict Comprehensions

```python
# ✅ GOOD
emails = [user["email"] for user in users if user["active"]]
index  = {user["id"]: user for user in users}

# ❌ BAD — verbose loop with no benefit
emails = []
for user in users:
    if user["active"]:
        emails.append(user["email"])
```

### 5. Catch Specific Exceptions

```python
# ✅ GOOD
try:
    result = int(user_input)
except ValueError as e:
    logger.warning("Invalid integer input", extra={"input": user_input, "error": str(e)})
    return None

# ❌ BAD — swallows ALL exceptions including KeyboardInterrupt
try:
    result = int(user_input)
except:
    pass
```

### 6. Use f-strings for Formatting

```python
# ✅ GOOD
msg   = f"User {user.name!r} logged in from {ip}"
price = f"Total: ${amount:.2f}"

# ❌ BAD — old-style formatting
msg = "User %s logged in from %s" % (user.name, ip)
msg = "User {} logged in from {}".format(user.name, ip)
```

### 7. Use the `logging` Module, Not `print()`

```python
# ✅ GOOD
import logging

logger = logging.getLogger(__name__)
logger.info("Processing order", extra={"order_id": order_id, "user_id": user_id})
logger.error("Payment failed", extra={"error": str(e)}, exc_info=True)

# ❌ BAD — unstructured, not filterable, breaks in production
print(f"Processing order {order_id}")
print(f"ERROR: {e}")
```

### 8. Use Pydantic for Data Validation

```python
# ✅ GOOD
from pydantic import BaseModel, EmailStr

class CreateUserRequest(BaseModel):
    name: str
    email: EmailStr
    age: int

    model_config = {"str_strip_whitespace": True}

# ❌ BAD — manual validation is error-prone and verbose
def create_user(data: dict):
    if "name" not in data:
        raise ValueError("name required")
    if "@" not in data.get("email", ""):
        raise ValueError("invalid email")
```

### 9. Avoid Implicit Truthiness When Distinguishing None from Empty

```python
# ✅ GOOD — explicit checks make intent clear
if items is None:
    raise ValueError("items must be provided")

if len(items) == 0:
    return default_result

# ❌ BAD — can't distinguish None vs [] vs 0 vs ""
if not items:
    raise ValueError("items must be provided")
```

### 10. Use `logging` and Structured Extras

```python
# ✅ GOOD — structured, queryable, filterable
import logging
import json

logging.basicConfig(
    level=logging.INFO,
    format="%(message)s",
)

logger = logging.getLogger(__name__)
logger.info(json.dumps({"event": "user_created", "user_id": user.id}))

# ❌ BAD — string concatenation loses structure
logging.info("User created: " + str(user.id))
```

---

## Commands

```bash
python -m mypy src/         # static type checking
python -m ruff check src/   # linting (fast, replaces flake8)
python -m ruff format src/  # formatting (replaces black)
python -m pytest            # run tests
python -m pytest --cov=src  # test coverage report
```

---

## Definition of Done

- [ ] All functions and methods have type hints
- [ ] No mutable default arguments (`def f(x=[])`)
- [ ] All file/DB/network resources use context managers (`with`)
- [ ] No bare `except:` clauses — exceptions are specific
- [ ] No `print()` in non-script code — use `logging`
- [ ] Data input/output validated with Pydantic or dataclasses
- [ ] `mypy` passes with no errors on the changed files

---



# Migration Agent

## When to Use

- Creating new tables or altering existing ones
- Adding indexes
- Designing PostgreSQL schemas

> Assets: See [assets/MIGRATION_TEMPLATE.sql](assets/MIGRATION_TEMPLATE.sql) for the standard template.

---

## Decision Tree

```
New table?                              → Use BIGINT GENERATED ALWAYS AS IDENTITY for PK
Need global / opaque ID?                → Add UUID external_id with DEFAULT gen_random_uuid()
Storing timestamps?                     → TIMESTAMPTZ — never plain TIMESTAMP
Storing money?                          → NUMERIC(p,s) — never FLOAT or MONEY
High-cardinality filter column?         → B-tree index
JSONB search?                           → GIN index
Time-series / append-only?             → BRIN index
Low-cardinality (boolean, status)?      → Partial index or skip
Adding FK column?                       → ALWAYS add index manually (PG doesn't auto-index FKs)
```

---

## Data Types

| Use case | Type | Avoid |
|----------|------|-------|
| Primary keys | `BIGINT GENERATED ALWAYS AS IDENTITY` | `SERIAL` |
| Global / opaque IDs | `UUID` | — |
| Timestamps | `TIMESTAMPTZ` | `TIMESTAMP` |
| Money / precision | `NUMERIC(p, s)` | `FLOAT`, `MONEY` |
| Strings | `TEXT` | `VARCHAR(n)` unless constraint is meaningful |
| JSON | `JSONB` | `JSON` |

---

## Schema Template

```sql
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
```

---

## Indexing Reference

| Cardinality | Type | Example |
|-------------|------|---------|
| High (email, UUID) | `B-tree unique` | `CREATE UNIQUE INDEX ON users (email)` |
| Medium filter (status) | `Partial B-tree` | `WHERE status = 'active'` |
| JSONB search | `GIN` | `USING GIN (payload)` |
| Time-series | `BRIN` | `USING BRIN (created_at)` |
| FK column | `B-tree` | Always — PG doesn't auto-index FKs |

---

## Migration Rules

1. Every `up` has a matching `down`.
2. Use `IF NOT EXISTS` / `IF EXISTS` — migrations must be idempotent.
3. Never alter a column type directly — add column, backfill, drop old (3-step).
4. `COMMENT ON COLUMN` is required for every field — no undocumented columns.
5. Tables > 100M rows need a partitioning plan.

---

## Commands

```bash
# Prisma
npx prisma migrate dev --name <migration_name>
npx prisma migrate status

# Knex
npx knex migrate:make <migration_name>
npx knex migrate:latest

# Check table indexes
psql -d mydb -c "\d <table_name>"
```

## Definition of Done

- [ ] Every column has `COMMENT ON COLUMN`
- [ ] Every FK column has an explicit index
- [ ] Index type matches cardinality
- [ ] `down` migration present and tested
- [ ] Migration is idempotent
- [ ] Tables > 100M rows have a partitioning plan

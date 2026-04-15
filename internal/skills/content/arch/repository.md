
# Repository Agent

## When to Use

- Writing database queries
- Managing transactions
- Mapping DB errors to domain exceptions
- Optimizing slow queries


---

## Decision Tree

```
Need to read multiple related rows?     → JOIN with explicit column list — no SELECT *
Updating multiple tables?               → Wrap in transaction with FOR UPDATE
Same operation could be called twice?   → Use ON CONFLICT DO NOTHING (upsert)
Got a DB error?                         → Map to DomainException before it leaves this layer
Large result set / paginated?           → Cursor-based (WHERE id > $cursor) — not OFFSET
```

---

## SQL Rules

**No SELECT ***
```sql
-- ❌ BAD
SELECT * FROM users WHERE id = $1;

-- ✅ GOOD
SELECT id, email, name, role, created_at FROM users WHERE id = $1;
```

**No functions on indexed columns in WHERE**
```sql
-- ❌ BAD — prevents index use
WHERE DATE(created_at) = '2024-01-01'

-- ✅ GOOD
WHERE created_at >= '2024-01-01' AND created_at < '2024-01-02'
```

**Cursor pagination over OFFSET**
```sql
-- ❌ BAD — O(n) cost
SELECT id, total FROM orders ORDER BY id LIMIT 20 OFFSET 10000;

-- ✅ GOOD — O(1) with index
SELECT id, total FROM orders WHERE id > $1 ORDER BY id LIMIT 20;
```

**Batch inserts over loops**
```sql
INSERT INTO events (user_id, type, payload)
SELECT * FROM UNNEST($1::uuid[], $2::text[], $3::jsonb[]);
```

---

## Transactions

```typescript
async transferFunds(fromId: string, toId: string, amount: number) {
  return this.db.transaction(async (trx) => {
    const from = await trx.query('SELECT balance FROM accounts WHERE id = $1 FOR UPDATE', [fromId]);
    if (from.rows[0].balance < amount)
      throw new DomainException('INSUFFICIENT_BALANCE', 'Transfer failed', 'from=' + fromId);
    await trx.query('UPDATE accounts SET balance = balance - $1 WHERE id = $2', [amount, fromId]);
    await trx.query('UPDATE accounts SET balance = balance + $1 WHERE id = $2', [amount, toId]);
  });
}
```

---

## Idempotency

```sql
INSERT INTO events (id, user_id, type, created_at)
VALUES ($1, $2, $3, NOW())
ON CONFLICT (id) DO NOTHING;
```

---

## DB Exception Mapping

```typescript
try {
  return await this.db.query(sql, params);
} catch (err) {
  if (err.code === '23505') throw new DomainException('DUPLICATE_ENTRY',       'Record already exists',        err.detail);
  if (err.code === '23503') throw new DomainException('FOREIGN_KEY_VIOLATION', 'Referenced record not found',  err.detail);
  if (err.code === '23502') throw new DomainException('NOT_NULL_VIOLATION',    'Required field missing',       err.detail);
  if (err.code === '40001') throw new DomainException('SERIALIZATION_FAILURE', 'Transaction conflict, retry',  err.detail);
  throw err;
}
```

| PG Code | Domain Code |
|---------|-------------|
| `23505` | `DUPLICATE_ENTRY` |
| `23503` | `FOREIGN_KEY_VIOLATION` |
| `23502` | `NOT_NULL_VIOLATION` |
| `40001` | `SERIALIZATION_FAILURE` |

---

## Commands

```bash
psql -d mydb -c "EXPLAIN ANALYZE <your query>"   # check query plan
psql -d mydb -c "\d <table>"                     # show table structure + indexes
```

## Definition of Done

- [ ] No `SELECT *` anywhere
- [ ] No functions wrapping indexed columns in WHERE
- [ ] Pagination uses cursor, not OFFSET
- [ ] Multi-step writes wrapped in transactions with FOR UPDATE
- [ ] Upserts for idempotent operations
- [ ] All DB errors mapped to `DomainException` before leaving this layer

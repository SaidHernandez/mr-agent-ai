-- Migration: <migration_name>
-- Created:   2026-03-26
-- Description: <what this migration does>

-- ── up.sql ───────────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS <table_name> (
  id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  external_id UUID        NOT NULL DEFAULT gen_random_uuid() UNIQUE,
  -- Add your columns here
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Column descriptions (REQUIRED for every field)
COMMENT ON COLUMN <table_name>.id          IS 'Internal surrogate key — never exposed externally';
COMMENT ON COLUMN <table_name>.external_id IS 'UUID exposed in API responses — immutable after creation';
COMMENT ON COLUMN <table_name>.created_at  IS 'Row creation timestamp with timezone';
COMMENT ON COLUMN <table_name>.updated_at  IS 'Last update timestamp with timezone';

-- Indexes (add for every FK and high-cardinality filter column)
-- CREATE UNIQUE INDEX idx_<table>_<col>  ON <table_name> (<col>);
-- CREATE INDEX        idx_<table>_fk     ON <table_name> (foreign_key_col);

-- ── down.sql ─────────────────────────────────────────────────────────────────

-- DROP TABLE IF EXISTS <table_name>;

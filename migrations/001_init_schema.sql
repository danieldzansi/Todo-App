-- Initial schema for Todo API
-- This file is executed by the Postgres container on first init

-- Ensure UUID generation helper is available (for gen_random_uuid)
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Core table
CREATE TABLE IF NOT EXISTS todos (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title       text        NOT NULL,
    description text,
    completed   boolean     NOT NULL DEFAULT false,
    due_date    timestamptz,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

-- Helpful indexes for common queries
CREATE INDEX IF NOT EXISTS idx_todos_created_at ON todos (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_todos_due_date   ON todos (due_date);
CREATE INDEX IF NOT EXISTS idx_todos_completed  ON todos (completed);

CREATE SCHEMA IF NOT EXISTS lab_migrations;

CREATE TABLE IF NOT EXISTS lab_migrations.versions (
    version integer PRIMARY KEY,
    name text NOT NULL,
    checksum text NOT NULL,
    applied_at timestamptz NOT NULL DEFAULT now()
);

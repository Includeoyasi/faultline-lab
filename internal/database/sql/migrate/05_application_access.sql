-- Reapply after migrations, including a run with no pending versions.
-- $1: business schema; $2: application role. The journal has no app grants.
SELECT format($sql$
GRANT USAGE ON SCHEMA %1$I TO %2$I;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA %1$I TO %2$I;
ALTER DEFAULT PRIVILEGES IN SCHEMA %1$I
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO %2$I;
$sql$, $1::text, $2::text);

-- $1: role name. Empty output means the role already exists.
SELECT CASE
    WHEN EXISTS (SELECT FROM pg_roles WHERE rolname = $1::text) THEN ''
    ELSE format('CREATE ROLE %I', $1::text)
END;

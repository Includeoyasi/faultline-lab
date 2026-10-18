-- $1: role name; $2: password. %I quotes an identifier, %L quotes a value.
SELECT format(
    'ALTER ROLE %I WITH LOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE '
    'NOREPLICATION NOBYPASSRLS NOINHERIT PASSWORD %L',
    $1::text, $2::text
);

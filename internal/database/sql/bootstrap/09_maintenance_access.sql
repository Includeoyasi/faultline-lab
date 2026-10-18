-- $1: maintenance database name.
SELECT format('REVOKE ALL ON DATABASE %I FROM PUBLIC', $1::text);

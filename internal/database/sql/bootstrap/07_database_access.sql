-- $1: database; $2: migration role; $3: application role.
SELECT format($sql$
REVOKE ALL ON DATABASE %1$I FROM PUBLIC;
GRANT CONNECT ON DATABASE %1$I TO %2$I, %3$I;
$sql$, $1::text, $2::text, $3::text);

-- $1: database name; $2: migration role. Execute outside a transaction.
SELECT format('CREATE DATABASE %I OWNER %I', $1::text, $2::text);

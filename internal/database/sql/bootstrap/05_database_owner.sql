SELECT pg_get_userbyid(datdba)
FROM pg_database
WHERE datname = $1;

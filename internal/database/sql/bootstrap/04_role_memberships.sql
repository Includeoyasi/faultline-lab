SELECT count(*)
FROM pg_auth_members AS membership
JOIN pg_roles AS role ON role.oid = membership.member
WHERE role.rolname = $1;

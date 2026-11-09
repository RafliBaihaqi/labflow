-- name: CreateUser :one
INSERT INTO users (xid, username, password, email, phone_number, name, company, address, created_at, updated_at)
VALUES (@xid, @username, @password, @email, @phone_number, @name, @company, @address, now(), now())
RETURNING *;

-- name: AddUserRole :exec
INSERT INTO user_roles (user_id, role_id) VALUES (@user_id, @role_id);

-- name: GetUserByID :one
SELECT * FROM users WHERE id = @id;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = @username ORDER BY id LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = @email ORDER BY id LIMIT 1;

-- name: ListUserPermissions :many
-- One row per (role, access). A role without accesses still gets one row.
SELECT r.role_id, r.role_name, a.access_id, a.access_name
FROM user_roles ur
JOIN roles r ON r.role_id = ur.role_id
LEFT JOIN role_accesses ra ON ra.role_id = r.role_id
LEFT JOIN accesses a ON a.access_id = ra.access_id
WHERE ur.user_id = @user_id
ORDER BY r.role_id, a.access_id;

-- name: SetUserVerified :one
UPDATE users SET is_verified = TRUE, updated_at = now()
WHERE email = @email
RETURNING *;

-- name: SetUserPassword :one
UPDATE users SET password = @password, updated_at = now()
WHERE email = @email
RETURNING *;

-- name: UpdateUserProfile :one
-- A nil argument sets the column to NULL, as the GORM map update did.
UPDATE users
SET phone_number = sqlc.narg(phone_number),
    company      = sqlc.narg(company),
    address      = sqlc.narg(address),
    avatar       = sqlc.narg(avatar),
    updated_at   = now(),
    updated_by   = @updated_by
WHERE xid = @xid
RETURNING *;

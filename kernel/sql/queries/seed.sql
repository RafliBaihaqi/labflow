-- Queries used only by the seeders. Re-running a seed is safe.

-- name: SeedAccess :exec
INSERT INTO accesses (access_id, access_name) VALUES (@access_id, @access_name)
ON CONFLICT DO NOTHING;

-- name: SeedRole :exec
INSERT INTO roles (role_id, role_name) VALUES (@role_id, @role_name)
ON CONFLICT DO NOTHING;

-- name: SeedCategory :exec
INSERT INTO categories (xid, name, label, created_at, updated_at)
VALUES (@xid, @name, sqlc.narg(label), now(), now())
ON CONFLICT DO NOTHING;

-- name: SeedTestList :exec
INSERT INTO test_lists (xid, category_id, value, label, created_at, updated_at)
VALUES (@xid, @category_id, @value, @label, now(), now())
ON CONFLICT DO NOTHING;

-- name: SeedUser :exec
INSERT INTO users (xid, username, password, email, phone_number, name, company, address, is_verified, created_at, updated_at)
VALUES (@xid, @username, @password, @email, @phone_number, @name, @company, @address,
        COALESCE(sqlc.narg(is_verified), FALSE), now(), now())
ON CONFLICT (username) DO NOTHING;

-- name: SeedUserRole :exec
INSERT INTO user_roles (user_id, role_id) VALUES (@user_id, @role_id)
ON CONFLICT DO NOTHING;

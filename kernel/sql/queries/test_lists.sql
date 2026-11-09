-- name: ListTestListsByCategory :many
SELECT sqlc.embed(tl), sqlc.embed(c)
FROM test_lists tl
JOIN categories c ON c.id = tl.category_id
WHERE tl.category_id = @category_id
ORDER BY tl.id;

-- name: CountTestLists :one
SELECT count(*) FROM test_lists
WHERE (sqlc.narg(search)::text IS NULL OR value ILIKE sqlc.narg(search) OR label ILIKE sqlc.narg(search))
  AND (sqlc.narg(category_id)::int IS NULL OR category_id = sqlc.narg(category_id));

-- name: ListTestLists :many
SELECT sqlc.embed(tl), sqlc.embed(c)
FROM test_lists tl
JOIN categories c ON c.id = tl.category_id
WHERE (sqlc.narg(search)::text IS NULL OR tl.value ILIKE sqlc.narg(search) OR tl.label ILIKE sqlc.narg(search))
  AND (sqlc.narg(category_id)::int IS NULL OR tl.category_id = sqlc.narg(category_id))
ORDER BY tl.created_at DESC
LIMIT @page_limit::int OFFSET @page_offset::int;

-- name: CreateTestList :one
-- version 0 and a nil is_current fall back to the column defaults, as GORM did.
INSERT INTO test_lists (
    xid, category_id, value, label, sample_amount_liquid, sample_amount_solid,
    sample_amount_powder, leadtime, version, is_current, created_at, updated_at
) VALUES (
    @xid, @category_id, @value, @label, sqlc.narg(sample_amount_liquid), sqlc.narg(sample_amount_solid),
    sqlc.narg(sample_amount_powder), sqlc.narg(leadtime), COALESCE(NULLIF(@version::int, 0), 1),
    COALESCE(sqlc.narg(is_current), TRUE), now(), now()
)
RETURNING id;

-- name: GetTestListByID :one
SELECT sqlc.embed(tl), sqlc.embed(c)
FROM test_lists tl
JOIN categories c ON c.id = tl.category_id
WHERE tl.id = @id;

-- name: GetTestListByXid :one
SELECT sqlc.embed(tl), sqlc.embed(c)
FROM test_lists tl
JOIN categories c ON c.id = tl.category_id
WHERE tl.xid = @xid
ORDER BY tl.id
LIMIT 1;

-- name: DeactivateTestList :exec
UPDATE test_lists SET is_current = FALSE, updated_at = now() WHERE xid = @xid;

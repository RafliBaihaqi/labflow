-- name: GetTestWithCategory :one
SELECT sqlc.embed(t), sqlc.embed(tl), sqlc.embed(c)
FROM tests t
JOIN test_lists tl ON tl.id = t.test_list_id
JOIN categories c ON c.id = tl.category_id
WHERE t.xid = @xid
ORDER BY t.id
LIMIT 1;

-- name: UpdateTestResult :one
-- A nil argument sets the column to NULL, as the GORM map update did.
UPDATE tests
SET test_result = sqlc.narg(test_result),
    test_summary = sqlc.narg(test_summary),
    notes = sqlc.narg(notes),
    updated_by = @updated_by,
    updated_at = now()
WHERE xid = @xid
RETURNING *;

-- name: UpdateTestApproval :one
UPDATE tests
SET approved_date = now(),
    updated_by = @updated_by,
    is_approved = sqlc.narg(is_approved),
    updated_at = now()
WHERE xid = @xid
RETURNING *;

-- name: CreateReservation :one
-- feedback_required is TRUE although the column default is FALSE: the GORM
-- model's default:true tag was written into every INSERT, so existing data
-- and the app expect TRUE.
INSERT INTO reservations (
    xid, user_id, reservation_date, sample_name, sample_amount, received_note,
    test_request_number, sample_residue, special_instructions, sample_type, note,
    batch_number, feedback_required, sample_taken_date, sample_taken_by, created_at, updated_at
) VALUES (
    @xid, @user_id, @reservation_date, @sample_name, @sample_amount, sqlc.narg(received_note),
    sqlc.narg(test_request_number), @sample_residue, sqlc.narg(special_instructions), @sample_type, sqlc.narg(note),
    COALESCE(sqlc.narg(batch_number), '-'), TRUE, @sample_taken_date, @sample_taken_by, now(), now()
)
RETURNING *;

-- name: CreateTests :many
-- Inserts all test lines of a reservation in one statement.
INSERT INTO tests (xid, reservation_id, test_list_id, method, criteria, test_unit, created_at, updated_at)
SELECT unnest(@xids::uuid[]), @reservation_id::int, unnest(@test_list_ids::int[]),
       unnest(@methods::text[]), unnest(@criteria::text[]), unnest(@test_units::text[]), now(), now()
RETURNING *;

-- name: CountReservations :one
SELECT count(*) FROM reservations
WHERE (sqlc.narg(search)::text IS NULL OR sample_name ILIKE sqlc.narg(search) OR test_request_number ILIKE sqlc.narg(search))
  AND (cardinality(@statuses::text[]) = 0 OR status = ANY(@statuses::text[]));

-- name: ListReservations :many
-- sqlc cannot bind a column name in ORDER BY, so each allowed sort is a CASE.
-- @sort is "<column>:<asc|desc>"; anything else falls through to created_at DESC.
SELECT sqlc.embed(r), sqlc.embed(u)
FROM reservations r
JOIN users u ON u.id = r.user_id
WHERE (sqlc.narg(search)::text IS NULL OR r.sample_name ILIKE sqlc.narg(search) OR r.test_request_number ILIKE sqlc.narg(search))
  AND (cardinality(@statuses::text[]) = 0 OR r.status = ANY(@statuses::text[]))
ORDER BY
    CASE WHEN @sort::text = 'test_request_number:asc' THEN r.test_request_number END ASC,
    CASE WHEN @sort::text = 'test_request_number:desc' THEN r.test_request_number END DESC,
    CASE WHEN @sort::text = 'sample_name:asc' THEN r.sample_name END ASC,
    CASE WHEN @sort::text = 'sample_name:desc' THEN r.sample_name END DESC,
    CASE WHEN @sort::text = 'status:asc' THEN r.status END ASC,
    CASE WHEN @sort::text = 'status:desc' THEN r.status END DESC,
    CASE WHEN @sort::text = 'reservation_date:asc' THEN r.reservation_date END ASC,
    CASE WHEN @sort::text = 'reservation_date:desc' THEN r.reservation_date END DESC,
    CASE WHEN @sort::text = 'created_at:asc' THEN r.created_at END ASC,
    r.created_at DESC
LIMIT sqlc.narg(page_limit)::int OFFSET @page_offset::int;

-- name: ListTestsByReservations :many
SELECT sqlc.embed(t), sqlc.embed(tl)
FROM tests t
JOIN test_lists tl ON tl.id = t.test_list_id
WHERE t.reservation_id = ANY(@reservation_ids::int[])
ORDER BY t.id;

-- name: CountReservationsByStatus :many
SELECT status, count(*) AS count FROM reservations GROUP BY status;

-- name: GetReservationWithUser :one
SELECT sqlc.embed(r), sqlc.embed(u)
FROM reservations r
JOIN users u ON u.id = r.user_id
WHERE r.xid = @xid
ORDER BY r.id
LIMIT 1;

-- name: ListTestsWithCategory :many
SELECT sqlc.embed(t), sqlc.embed(tl), sqlc.embed(c)
FROM tests t
JOIN test_lists tl ON tl.id = t.test_list_id
JOIN categories c ON c.id = tl.category_id
WHERE t.reservation_id = @reservation_id
ORDER BY t.id;

-- name: UpdateReservationDocuments :one
-- A nil argument sets the column to NULL, as the GORM map update did.
UPDATE reservations
SET po_url = sqlc.narg(po_url),
    po_number = sqlc.narg(po_number),
    sample_photo = sqlc.narg(sample_photo),
    updated_at = now(),
    updated_by = @updated_by
WHERE xid = @xid
RETURNING *;

-- name: UpdateReservationStatus :one
-- sample_photo only changes when a new one is given.
UPDATE reservations
SET status = @status,
    updated_by = @updated_by,
    sample_photo = COALESCE(sqlc.narg(sample_photo), sample_photo),
    updated_at = now()
WHERE xid = @xid
RETURNING *;

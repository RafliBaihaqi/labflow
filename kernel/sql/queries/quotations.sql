-- name: CreateQuotation :one
INSERT INTO quotations (xid, email, company_name, phone_number, message, created_at, updated_at)
VALUES (@xid, @email, @company_name, sqlc.narg(phone_number), sqlc.narg(message), now(), now())
RETURNING *;

-- name: CountQuotations :one
SELECT count(*) FROM quotations
WHERE sqlc.narg(search)::text IS NULL
   OR email ILIKE sqlc.narg(search) OR company_name ILIKE sqlc.narg(search) OR phone_number ILIKE sqlc.narg(search);

-- name: ListQuotations :many
SELECT * FROM quotations
WHERE sqlc.narg(search)::text IS NULL
   OR email ILIKE sqlc.narg(search) OR company_name ILIKE sqlc.narg(search) OR phone_number ILIKE sqlc.narg(search)
ORDER BY created_at DESC
LIMIT @page_limit::int OFFSET @page_offset::int;

-- name: UpdateQuotationStatus :one
UPDATE quotations SET status = @status, updated_at = now()
WHERE xid = @xid
RETURNING *;

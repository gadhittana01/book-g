-- name: CreateBook :one
INSERT INTO "book"(title, description, author, price) VALUES
($1, $2, $3, $4) RETURNING *;

-- name: CheckBookExists :one
SELECT EXISTS(SELECT id FROM "book" WHERE id=$1);

-- name: FindBookByID :one
SELECT * FROM "book" WHERE id=$1;

-- name: FindBook :many
SELECT * FROM "book" AS b
ORDER BY b.created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetBookCount :one
SELECT COUNT(o.*) FROM (SELECT * FROM "book" AS b) AS o;
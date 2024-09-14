-- name: CheckBookExists :one
SELECT EXISTS(SELECT id FROM "book" WHERE id=$1);
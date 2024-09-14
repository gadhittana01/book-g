-- name: CreateUser :one
INSERT INTO "user"(name, email, password) VALUES
($1, $2, $3) RETURNING *;

-- name: FindUserByEmail :one
SELECT * FROM "user" WHERE email=$1;

-- name: CheckEmailExists :one
SELECT EXISTS(SELECT id FROM "user" WHERE email=$1);
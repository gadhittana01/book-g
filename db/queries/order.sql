-- name: CreateOrder :one
INSERT INTO "order"(user_id, date) VALUES
($1, $2) RETURNING *;

-- name: CreateOrderDetail :one
INSERT INTO "order_detail"(order_id, book_id, quantity) VALUES
($1, $2, $3) RETURNING *;

-- name: FindOrderByUserID :many
SELECT * FROM "order" AS o
WHERE o.user_id=$1
ORDER BY o.date DESC
LIMIT $2 OFFSET $3;

-- name: FindOrderByID :one
SELECT * FROM "order" AS o
WHERE o.user_id=$1 AND o.id=$2;

-- name: GetOrderCountByUserId :one
SELECT COUNT(o.*) FROM (SELECT * FROM "order" AS o
WHERE o.user_id=$1) AS o;

-- name: CheckOrderExists :one
SELECT EXISTS(SELECT id FROM "order" WHERE id=$1);

-- name: FindOrderDetailByOrderID :many
SELECT 
    o.id, o.date, od.book_id, b.title,
    b.description, b.author, od.quantity
FROM "order" AS o
JOIN "order_detail" od 
ON o.id = od.order_id JOIN "book" AS b
ON od.book_id = b.id
WHERE o.user_id=$1 AND o.id=$2;
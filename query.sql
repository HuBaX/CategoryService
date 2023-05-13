-- name: AddCategory :execresult
INSERT INTO category(name) VALUES(?);

-- name: GetCategory :one
SELECT * FROM category WHERE id=?;

-- name: GetCategories :many
SELECT * FROM category;

-- name: GetCategoryByName :many
SELECT * FROM category WHERE name=?;

-- name: DelCategory :exec
DELETE FROM category WHERE id=?;

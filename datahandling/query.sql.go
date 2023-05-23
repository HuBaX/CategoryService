// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: query.sql

package datahandling

import (
	"context"
)

const addCategory = `-- name: AddCategory :exec
INSERT INTO category(name) VALUES(?)
`

func (q *Queries) AddCategory(ctx context.Context, name string) error {
	_, err := q.db.ExecContext(ctx, addCategory, name)
	return err
}

const delCategory = `-- name: DelCategory :exec
DELETE FROM category WHERE id=?
`

func (q *Queries) DelCategory(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, delCategory, id)
	return err
}

const getCategories = `-- name: GetCategories :many
SELECT id, name FROM category
`

func (q *Queries) GetCategories(ctx context.Context) ([]Category, error) {
	rows, err := q.db.QueryContext(ctx, getCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Category
	for rows.Next() {
		var i Category
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCategory = `-- name: GetCategory :one
SELECT id, name FROM category WHERE id=?
`

func (q *Queries) GetCategory(ctx context.Context, id int32) (Category, error) {
	row := q.db.QueryRowContext(ctx, getCategory, id)
	var i Category
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getCategoryByName = `-- name: GetCategoryByName :many
SELECT id, name FROM category WHERE name=?
`

func (q *Queries) GetCategoryByName(ctx context.Context, name string) ([]Category, error) {
	rows, err := q.db.QueryContext(ctx, getCategoryByName, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Category
	for rows.Next() {
		var i Category
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

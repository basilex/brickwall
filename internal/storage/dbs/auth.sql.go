// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: auth.sql

package dbs

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const authSelectUserCredentials = `-- name: AuthSelectUserCredentials :one
select id, username, password, is_blocked, is_checked, blocked_at, checked_at
  from users u
 where u.username = $1
`

type AuthSelectUserCredentialsRow struct {
	ID        string           `json:"id"`
	Username  string           `json:"username"`
	Password  string           `json:"password"`
	IsBlocked bool             `json:"is_blocked"`
	IsChecked bool             `json:"is_checked"`
	BlockedAt pgtype.Timestamp `json:"blocked_at"`
	CheckedAt pgtype.Timestamp `json:"checked_at"`
}

// AuthSelectUserCredentials
//
//	select id, username, password, is_blocked, is_checked, blocked_at, checked_at
//	  from users u
//	 where u.username = $1
func (q *Queries) AuthSelectUserCredentials(ctx context.Context, username string) (*AuthSelectUserCredentialsRow, error) {
	row := q.db.QueryRow(ctx, authSelectUserCredentials, username)
	var i AuthSelectUserCredentialsRow
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.IsBlocked,
		&i.IsChecked,
		&i.BlockedAt,
		&i.CheckedAt,
	)
	return &i, err
}

const authUpdateVisitedAt = `-- name: AuthUpdateVisitedAt :one
update users set visited_at = now() where id = $1
       returning id, username, checked_at, visited_at, created_at
`

type AuthUpdateVisitedAtRow struct {
	ID        string           `json:"id"`
	Username  string           `json:"username"`
	CheckedAt pgtype.Timestamp `json:"checked_at"`
	VisitedAt pgtype.Timestamp `json:"visited_at"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

// AuthUpdateVisitedAt
//
//	update users set visited_at = now() where id = $1
//	       returning id, username, checked_at, visited_at, created_at
func (q *Queries) AuthUpdateVisitedAt(ctx context.Context, id string) (*AuthUpdateVisitedAtRow, error) {
	row := q.db.QueryRow(ctx, authUpdateVisitedAt, id)
	var i AuthUpdateVisitedAtRow
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.CheckedAt,
		&i.VisitedAt,
		&i.CreatedAt,
	)
	return &i, err
}

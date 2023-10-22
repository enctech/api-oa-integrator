// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package database

import (
	"context"

	"github.com/sqlc-dev/pqtype"
)

const createLog = `-- name: CreateLog :one
INSERT INTO logs (module, info, extra)
VALUES ($1, $2, $3)
RETURNING id, module, info, extra
`

type CreateLogParams struct {
	Module string
	Info   string
	Extra  pqtype.NullRawMessage
}

func (q *Queries) CreateLog(ctx context.Context, arg CreateLogParams) (Log, error) {
	row := q.db.QueryRowContext(ctx, createLog, arg.Module, arg.Info, arg.Extra)
	var i Log
	err := row.Scan(
		&i.ID,
		&i.Module,
		&i.Info,
		&i.Extra,
	)
	return i, err
}

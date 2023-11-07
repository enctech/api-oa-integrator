// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

const createConfig = `-- name: CreateConfig :one
insert into snb_config (endpoint, facility, device)
values ($1, $2, $3)
returning id, endpoint, facility, device
`

type CreateConfigParams struct {
	Endpoint sql.NullString
	Facility sql.NullString
	Device   sql.NullString
}

func (q *Queries) CreateConfig(ctx context.Context, arg CreateConfigParams) (SnbConfig, error) {
	row := q.db.QueryRowContext(ctx, createConfig, arg.Endpoint, arg.Facility, arg.Device)
	var i SnbConfig
	err := row.Scan(
		&i.ID,
		&i.Endpoint,
		&i.Facility,
		&i.Device,
	)
	return i, err
}

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

const createUser = `-- name: CreateUser :one
insert into users (username, password, permission)
values ($1, $2, $3)
returning id, username, password, permission
`

type CreateUserParams struct {
	Username   sql.NullString
	Password   sql.NullString
	Permission sql.NullString
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.Password, arg.Permission)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Permission,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :execresult
delete
from users
where id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) (sql.Result, error) {
	return q.db.ExecContext(ctx, deleteUser, id)
}

const getConfig = `-- name: GetConfig :one
select id, endpoint, facility, device
from snb_config
where facility = $1
  and device = $2
`

type GetConfigParams struct {
	Facility sql.NullString
	Device   sql.NullString
}

func (q *Queries) GetConfig(ctx context.Context, arg GetConfigParams) (SnbConfig, error) {
	row := q.db.QueryRowContext(ctx, getConfig, arg.Facility, arg.Device)
	var i SnbConfig
	err := row.Scan(
		&i.ID,
		&i.Endpoint,
		&i.Facility,
		&i.Device,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
select id, username, password, permission
from users
where username = $1
`

func (q *Queries) GetUser(ctx context.Context, username sql.NullString) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Permission,
	)
	return i, err
}

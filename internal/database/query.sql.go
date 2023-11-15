// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/sqlc-dev/pqtype"
)

const createConfig = `-- name: CreateConfig :one
insert into snb_config (endpoint, facility, device)
values ($1, $2, $3)
returning id, endpoint, facility, device
`

type CreateConfigParams struct {
	Endpoint sql.NullString
	Facility []string
	Device   []string
}

// -----------Region S&B Config start-------------
func (q *Queries) CreateConfig(ctx context.Context, arg CreateConfigParams) (SnbConfig, error) {
	row := q.db.QueryRowContext(ctx, createConfig, arg.Endpoint, pq.Array(arg.Facility), pq.Array(arg.Device))
	var i SnbConfig
	err := row.Scan(
		&i.ID,
		&i.Endpoint,
		pq.Array(&i.Facility),
		pq.Array(&i.Device),
	)
	return i, err
}

const createLog = `-- name: CreateLog :one
INSERT INTO logs (level, message, fields, created_at)
VALUES ($1, $2, $3, $4)
RETURNING id, level, message, fields, created_at
`

type CreateLogParams struct {
	Level     sql.NullString
	Message   sql.NullString
	Fields    pqtype.NullRawMessage
	CreatedAt time.Time
}

func (q *Queries) CreateLog(ctx context.Context, arg CreateLogParams) (Log, error) {
	row := q.db.QueryRowContext(ctx, createLog,
		arg.Level,
		arg.Message,
		arg.Fields,
		arg.CreatedAt,
	)
	var i Log
	err := row.Scan(
		&i.ID,
		&i.Level,
		&i.Message,
		&i.Fields,
		&i.CreatedAt,
	)
	return i, err
}

const createOATransaction = `-- name: CreateOATransaction :one

insert into oa_transactions (businesstransactionid, lpn, customerid, jobid, facility, device, extra, entry_lane,
                             exit_lane)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
returning id, businesstransactionid, lpn, customerid, jobid, facility, device, extra, entry_lane, exit_lane, created_at, updated_at
`

type CreateOATransactionParams struct {
	Businesstransactionid string
	Lpn                   sql.NullString
	Customerid            sql.NullString
	Jobid                 sql.NullString
	Facility              sql.NullString
	Device                sql.NullString
	Extra                 pqtype.NullRawMessage
	EntryLane             sql.NullString
	ExitLane              sql.NullString
}

// -----------Region S&B Config end---------------
// -----------Region OA Transaction start-------------
func (q *Queries) CreateOATransaction(ctx context.Context, arg CreateOATransactionParams) (OaTransaction, error) {
	row := q.db.QueryRowContext(ctx, createOATransaction,
		arg.Businesstransactionid,
		arg.Lpn,
		arg.Customerid,
		arg.Jobid,
		arg.Facility,
		arg.Device,
		arg.Extra,
		arg.EntryLane,
		arg.ExitLane,
	)
	var i OaTransaction
	err := row.Scan(
		&i.ID,
		&i.Businesstransactionid,
		&i.Lpn,
		&i.Customerid,
		&i.Jobid,
		&i.Facility,
		&i.Device,
		&i.Extra,
		&i.EntryLane,
		&i.ExitLane,
		&i.CreatedAt,
		&i.UpdatedAt,
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

// -----------Region Integrator Config end---------------
// -----------Region Authentication start-------------
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
where facility in ($1)
  and device in ($2)
`

type GetConfigParams struct {
	Facility []string
	Device   []string
}

func (q *Queries) GetConfig(ctx context.Context, arg GetConfigParams) (SnbConfig, error) {
	row := q.db.QueryRowContext(ctx, getConfig, pq.Array(arg.Facility), pq.Array(arg.Device))
	var i SnbConfig
	err := row.Scan(
		&i.ID,
		&i.Endpoint,
		pq.Array(&i.Facility),
		pq.Array(&i.Device),
	)
	return i, err
}

const getIntegratorConfig = `-- name: GetIntegratorConfig :one

select id, client_id, name, sp_id, plaza_id, url, insecure_skip_verify, created_at, updated_at
from integrator_config
where client_id = $1
`

// -----------Region OA Transaction end-------------
// -----------Region Integrator Config start-------------
func (q *Queries) GetIntegratorConfig(ctx context.Context, clientID sql.NullString) (IntegratorConfig, error) {
	row := q.db.QueryRowContext(ctx, getIntegratorConfig, clientID)
	var i IntegratorConfig
	err := row.Scan(
		&i.ID,
		&i.ClientID,
		&i.Name,
		&i.SpID,
		&i.PlazaID,
		&i.Url,
		&i.InsecureSkipVerify,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getOATransaction = `-- name: GetOATransaction :one
select id, businesstransactionid, lpn, customerid, jobid, facility, device, extra, entry_lane, exit_lane, created_at, updated_at
from oa_transactions
where businesstransactionid = $1
`

func (q *Queries) GetOATransaction(ctx context.Context, businesstransactionid string) (OaTransaction, error) {
	row := q.db.QueryRowContext(ctx, getOATransaction, businesstransactionid)
	var i OaTransaction
	err := row.Scan(
		&i.ID,
		&i.Businesstransactionid,
		&i.Lpn,
		&i.Customerid,
		&i.Jobid,
		&i.Facility,
		&i.Device,
		&i.Extra,
		&i.EntryLane,
		&i.ExitLane,
		&i.CreatedAt,
		&i.UpdatedAt,
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

const updateOATransaction = `-- name: UpdateOATransaction :one
update oa_transactions
set lpn        = coalesce($2, lpn),
    customerid = coalesce($3, customerid),
    jobid      = coalesce($4, jobid),
    facility   = coalesce($5, facility),
    device     = coalesce($6, device),
    extra      = coalesce($7, extra),
    exit_lane  = coalesce($8, extra)
where businesstransactionid = $1
returning id, businesstransactionid, lpn, customerid, jobid, facility, device, extra, entry_lane, exit_lane, created_at, updated_at
`

type UpdateOATransactionParams struct {
	Businesstransactionid string
	Lpn                   sql.NullString
	Customerid            sql.NullString
	Jobid                 sql.NullString
	Facility              sql.NullString
	Device                sql.NullString
	Extra                 pqtype.NullRawMessage
	ExitLane              sql.NullString
}

func (q *Queries) UpdateOATransaction(ctx context.Context, arg UpdateOATransactionParams) (OaTransaction, error) {
	row := q.db.QueryRowContext(ctx, updateOATransaction,
		arg.Businesstransactionid,
		arg.Lpn,
		arg.Customerid,
		arg.Jobid,
		arg.Facility,
		arg.Device,
		arg.Extra,
		arg.ExitLane,
	)
	var i OaTransaction
	err := row.Scan(
		&i.ID,
		&i.Businesstransactionid,
		&i.Lpn,
		&i.Customerid,
		&i.Jobid,
		&i.Facility,
		&i.Device,
		&i.Extra,
		&i.EntryLane,
		&i.ExitLane,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

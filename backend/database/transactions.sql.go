// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: transactions.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

const createIntegratorTransaction = `-- name: CreateIntegratorTransaction :one
with inserted_transaction as (
    insert into integrator_transactions (business_transaction_id, lpn, integrator_id, status, amount, error, tax_data,
                                         extra)
        values ($1, $2, $3, $4, $5, $6, $7, $8)
        returning id, business_transaction_id, lpn, integrator_id, status, amount, error, extra, tax_data, created_at, updated_at)
select inserted_transaction.id, business_transaction_id, lpn, integrator_id, status, amount, error, inserted_transaction.extra, tax_data, inserted_transaction.created_at, inserted_transaction.updated_at, integrator_config.id, client_id, provider_id, name, display_name, integrator_name, sp_id, plaza_id_map, integrator_config.extra, url, tax_rate, surcharge, surchange_type, insecure_skip_verify, integrator_config.created_at, integrator_config.updated_at
from inserted_transaction
         inner join integrator_config on integrator_config.id = $3
`

type CreateIntegratorTransactionParams struct {
	BusinessTransactionID uuid.UUID
	Lpn                   sql.NullString
	ID                    uuid.UUID
	Status                sql.NullString
	Amount                sql.NullString
	Error                 sql.NullString
	TaxData               pqtype.NullRawMessage
	Extra                 pqtype.NullRawMessage
}

type CreateIntegratorTransactionRow struct {
	ID                    uuid.UUID
	BusinessTransactionID uuid.UUID
	Lpn                   sql.NullString
	IntegratorID          uuid.NullUUID
	Status                sql.NullString
	Amount                sql.NullString
	Error                 sql.NullString
	Extra                 pqtype.NullRawMessage
	TaxData               pqtype.NullRawMessage
	CreatedAt             time.Time
	UpdatedAt             time.Time
	ID_2                  uuid.UUID
	ClientID              sql.NullString
	ProviderID            sql.NullInt32
	Name                  sql.NullString
	DisplayName           sql.NullString
	IntegratorName        sql.NullString
	SpID                  sql.NullString
	PlazaIDMap            pqtype.NullRawMessage
	Extra_2               pqtype.NullRawMessage
	Url                   sql.NullString
	TaxRate               sql.NullString
	Surcharge             sql.NullString
	SurchangeType         NullSurchargeType
	InsecureSkipVerify    sql.NullBool
	CreatedAt_2           time.Time
	UpdatedAt_2           time.Time
}

func (q *Queries) CreateIntegratorTransaction(ctx context.Context, arg CreateIntegratorTransactionParams) (CreateIntegratorTransactionRow, error) {
	row := q.queryRow(ctx, q.createIntegratorTransactionStmt, createIntegratorTransaction,
		arg.BusinessTransactionID,
		arg.Lpn,
		arg.ID,
		arg.Status,
		arg.Amount,
		arg.Error,
		arg.TaxData,
		arg.Extra,
	)
	var i CreateIntegratorTransactionRow
	err := row.Scan(
		&i.ID,
		&i.BusinessTransactionID,
		&i.Lpn,
		&i.IntegratorID,
		&i.Status,
		&i.Amount,
		&i.Error,
		&i.Extra,
		&i.TaxData,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ID_2,
		&i.ClientID,
		&i.ProviderID,
		&i.Name,
		&i.DisplayName,
		&i.IntegratorName,
		&i.SpID,
		&i.PlazaIDMap,
		&i.Extra_2,
		&i.Url,
		&i.TaxRate,
		&i.Surcharge,
		&i.SurchangeType,
		&i.InsecureSkipVerify,
		&i.CreatedAt_2,
		&i.UpdatedAt_2,
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

func (q *Queries) CreateOATransaction(ctx context.Context, arg CreateOATransactionParams) (OaTransaction, error) {
	row := q.queryRow(ctx, q.createOATransactionStmt, createOATransaction,
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

const getIntegratorTransactions = `-- name: GetIntegratorTransactions :many
select it.id, it.business_transaction_id, it.lpn, it.integrator_id, it.status, it.amount, it.error, it.extra, it.tax_data, it.created_at, it.updated_at, ic.name as integrator_name
from integrator_transactions it
         inner join public.integrator_config ic on ic.id = it.integrator_id
where lpn like concat('%', $3::text, '%')
  and name::text like concat('%', $4::text, '%')
  and status::text like concat('%', $5::text, '%')
  and it.created_at >= $6
  and it.created_at <= $7
order by it.created_at desc
limit $1 offset $2
`

type GetIntegratorTransactionsParams struct {
	Limit   int32
	Offset  int32
	Lpn     string
	Name    string
	Status  string
	StartAt time.Time
	EndAt   time.Time
}

type GetIntegratorTransactionsRow struct {
	ID                    uuid.UUID
	BusinessTransactionID uuid.UUID
	Lpn                   sql.NullString
	IntegratorID          uuid.NullUUID
	Status                sql.NullString
	Amount                sql.NullString
	Error                 sql.NullString
	Extra                 pqtype.NullRawMessage
	TaxData               pqtype.NullRawMessage
	CreatedAt             time.Time
	UpdatedAt             time.Time
	IntegratorName        sql.NullString
}

func (q *Queries) GetIntegratorTransactions(ctx context.Context, arg GetIntegratorTransactionsParams) ([]GetIntegratorTransactionsRow, error) {
	rows, err := q.query(ctx, q.getIntegratorTransactionsStmt, getIntegratorTransactions,
		arg.Limit,
		arg.Offset,
		arg.Lpn,
		arg.Name,
		arg.Status,
		arg.StartAt,
		arg.EndAt,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetIntegratorTransactionsRow{}
	for rows.Next() {
		var i GetIntegratorTransactionsRow
		if err := rows.Scan(
			&i.ID,
			&i.BusinessTransactionID,
			&i.Lpn,
			&i.IntegratorID,
			&i.Status,
			&i.Amount,
			&i.Error,
			&i.Extra,
			&i.TaxData,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.IntegratorName,
		); err != nil {
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

const getIntegratorTransactionsCount = `-- name: GetIntegratorTransactionsCount :one
select count(*)
from integrator_transactions it
         inner join public.integrator_config ic on ic.id = it.integrator_id
where lpn like concat('%', $1::text, '%')
  and name::text like concat('%', $2::text, '%')
  and status::text like concat('%', $3::text, '%')
  and it.created_at >= $4
  and it.created_at <= $5
`

type GetIntegratorTransactionsCountParams struct {
	Lpn     string
	Name    string
	Status  string
	StartAt time.Time
	EndAt   time.Time
}

func (q *Queries) GetIntegratorTransactionsCount(ctx context.Context, arg GetIntegratorTransactionsCountParams) (int64, error) {
	row := q.queryRow(ctx, q.getIntegratorTransactionsCountStmt, getIntegratorTransactionsCount,
		arg.Lpn,
		arg.Name,
		arg.Status,
		arg.StartAt,
		arg.EndAt,
	)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getLatestOATransaction = `-- name: GetLatestOATransaction :one
select id, businesstransactionid, lpn, customerid, jobid, facility, device, extra, entry_lane, exit_lane, created_at, updated_at
from oa_transactions
where businesstransactionid = $1
order by created_at desc
limit 1
`

func (q *Queries) GetLatestOATransaction(ctx context.Context, businesstransactionid string) (OaTransaction, error) {
	row := q.queryRow(ctx, q.getLatestOATransactionStmt, getLatestOATransaction, businesstransactionid)
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

const getLatestOATransactions = `-- name: GetLatestOATransactions :many
select id, businesstransactionid, lpn, customerid, jobid, facility, device, extra, entry_lane, exit_lane, created_at, updated_at
from oa_transactions
where updated_at >= $3
  and updated_at <= $4
order by updated_at desc
limit $1 offset $2
`

type GetLatestOATransactionsParams struct {
	Limit   int32
	Offset  int32
	StartAt time.Time
	EndAt   time.Time
}

func (q *Queries) GetLatestOATransactions(ctx context.Context, arg GetLatestOATransactionsParams) ([]OaTransaction, error) {
	rows, err := q.query(ctx, q.getLatestOATransactionsStmt, getLatestOATransactions,
		arg.Limit,
		arg.Offset,
		arg.StartAt,
		arg.EndAt,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []OaTransaction{}
	for rows.Next() {
		var i OaTransaction
		if err := rows.Scan(
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
		); err != nil {
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

const getOAEntryTransactions = `-- name: GetOAEntryTransactions :one
select count(*)
from oa_transactions
where extra ->> 'steps' = 'leave_loop_entry_done'
  and created_at >= $1
  and created_at <= $2
`

type GetOAEntryTransactionsParams struct {
	StartAt time.Time
	EndAt   time.Time
}

func (q *Queries) GetOAEntryTransactions(ctx context.Context, arg GetOAEntryTransactionsParams) (int64, error) {
	row := q.queryRow(ctx, q.getOAEntryTransactionsStmt, getOAEntryTransactions, arg.StartAt, arg.EndAt)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getOAExitTransactions = `-- name: GetOAExitTransactions :one
select count(*)
from oa_transactions
where extra ->> 'steps' = 'exit_leave_loop_done'
  and created_at >= $1
  and created_at <= $2
`

type GetOAExitTransactionsParams struct {
	StartAt time.Time
	EndAt   time.Time
}

func (q *Queries) GetOAExitTransactions(ctx context.Context, arg GetOAExitTransactionsParams) (int64, error) {
	row := q.queryRow(ctx, q.getOAExitTransactionsStmt, getOAExitTransactions, arg.StartAt, arg.EndAt)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getOATransaction = `-- name: GetOATransaction :one
select id, businesstransactionid, lpn, customerid, jobid, facility, device, extra, entry_lane, exit_lane, created_at, updated_at
from oa_transactions
where businesstransactionid = $1
`

func (q *Queries) GetOATransaction(ctx context.Context, businesstransactionid string) (OaTransaction, error) {
	row := q.queryRow(ctx, q.getOATransactionStmt, getOATransaction, businesstransactionid)
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

const getOATransactions = `-- name: GetOATransactions :many
select id, businesstransactionid, lpn, customerid, jobid, facility, device, extra, entry_lane, exit_lane, created_at, updated_at
from oa_transactions
where lpn like concat('%', $3::text, '%')
  and jobid::text like concat('%', $4::text, '%')
  and facility::text like concat('%', $5::text, '%')
  and entry_lane::text like concat('%', $6::text, '%')
  and (exit_lane::text LIKE concat('%', $7::text, '%') or
       (exit_lane is null and ($7::text) = ''))
  and created_at >= $8
  and created_at <= $9
order by created_at desc
limit $1 offset $2
`

type GetOATransactionsParams struct {
	Limit     int32
	Offset    int32
	Lpn       string
	Jobid     string
	Facility  string
	EntryLane string
	ExitLane  string
	StartAt   time.Time
	EndAt     time.Time
}

func (q *Queries) GetOATransactions(ctx context.Context, arg GetOATransactionsParams) ([]OaTransaction, error) {
	rows, err := q.query(ctx, q.getOATransactionsStmt, getOATransactions,
		arg.Limit,
		arg.Offset,
		arg.Lpn,
		arg.Jobid,
		arg.Facility,
		arg.EntryLane,
		arg.ExitLane,
		arg.StartAt,
		arg.EndAt,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []OaTransaction{}
	for rows.Next() {
		var i OaTransaction
		if err := rows.Scan(
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
		); err != nil {
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

const getOATransactionsCount = `-- name: GetOATransactionsCount :one
select count(*)
from oa_transactions
where lpn like concat('%', $1::text, '%')
  and jobid::text like concat('%', $2::text, '%')
  and facility::text like concat('%', $3::text, '%')
  and entry_lane::text like concat('%', $4::text, '%')
  and (exit_lane::text LIKE concat('%', $5::text, '%') or
       (exit_lane is null and ($5::text) = ''))
  and created_at >= $6
  and created_at <= $7
`

type GetOATransactionsCountParams struct {
	Lpn       string
	Jobid     string
	Facility  string
	EntryLane string
	ExitLane  string
	StartAt   time.Time
	EndAt     time.Time
}

func (q *Queries) GetOATransactionsCount(ctx context.Context, arg GetOATransactionsCountParams) (int64, error) {
	row := q.queryRow(ctx, q.getOATransactionsCountStmt, getOATransactionsCount,
		arg.Lpn,
		arg.Jobid,
		arg.Facility,
		arg.EntryLane,
		arg.ExitLane,
		arg.StartAt,
		arg.EndAt,
	)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getTotalTransactionAmount = `-- name: GetTotalTransactionAmount :one
select sum(amount)::numeric
from integrator_transactions
where status::text like concat('%', $1::text, '%')
  and created_at >= $2
  and created_at <= $3
`

type GetTotalTransactionAmountParams struct {
	Status  string
	StartAt time.Time
	EndAt   time.Time
}

func (q *Queries) GetTotalTransactionAmount(ctx context.Context, arg GetTotalTransactionAmountParams) (string, error) {
	row := q.queryRow(ctx, q.getTotalTransactionAmountStmt, getTotalTransactionAmount, arg.Status, arg.StartAt, arg.EndAt)
	var column_1 string
	err := row.Scan(&column_1)
	return column_1, err
}

const updateOATransaction = `-- name: UpdateOATransaction :one
update oa_transactions
set lpn        = coalesce($2, lpn),
    customerid = coalesce($3, customerid),
    jobid      = coalesce($4, jobid),
    facility   = coalesce($5, facility),
    device     = coalesce($6, device),
    extra      = coalesce($7, extra),
    exit_lane  = coalesce($8, exit_lane)
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
	row := q.queryRow(ctx, q.updateOATransactionStmt, updateOATransaction,
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

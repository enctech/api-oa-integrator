// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package database

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.countLogsStmt, err = db.PrepareContext(ctx, countLogs); err != nil {
		return nil, fmt.Errorf("error preparing query CountLogs: %w", err)
	}
	if q.createIntegratorConfigStmt, err = db.PrepareContext(ctx, createIntegratorConfig); err != nil {
		return nil, fmt.Errorf("error preparing query CreateIntegratorConfig: %w", err)
	}
	if q.createLogStmt, err = db.PrepareContext(ctx, createLog); err != nil {
		return nil, fmt.Errorf("error preparing query CreateLog: %w", err)
	}
	if q.createOATransactionStmt, err = db.PrepareContext(ctx, createOATransaction); err != nil {
		return nil, fmt.Errorf("error preparing query CreateOATransaction: %w", err)
	}
	if q.createSnbConfigStmt, err = db.PrepareContext(ctx, createSnbConfig); err != nil {
		return nil, fmt.Errorf("error preparing query CreateSnbConfig: %w", err)
	}
	if q.createUserStmt, err = db.PrepareContext(ctx, createUser); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUser: %w", err)
	}
	if q.deleteIntegratorConfigStmt, err = db.PrepareContext(ctx, deleteIntegratorConfig); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteIntegratorConfig: %w", err)
	}
	if q.deleteSnbConfigStmt, err = db.PrepareContext(ctx, deleteSnbConfig); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteSnbConfig: %w", err)
	}
	if q.deleteUserStmt, err = db.PrepareContext(ctx, deleteUser); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteUser: %w", err)
	}
	if q.getAllSnbConfigStmt, err = db.PrepareContext(ctx, getAllSnbConfig); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllSnbConfig: %w", err)
	}
	if q.getIntegratorConfigStmt, err = db.PrepareContext(ctx, getIntegratorConfig); err != nil {
		return nil, fmt.Errorf("error preparing query GetIntegratorConfig: %w", err)
	}
	if q.getIntegratorConfigByClientStmt, err = db.PrepareContext(ctx, getIntegratorConfigByClient); err != nil {
		return nil, fmt.Errorf("error preparing query GetIntegratorConfigByClient: %w", err)
	}
	if q.getIntegratorConfigByNameStmt, err = db.PrepareContext(ctx, getIntegratorConfigByName); err != nil {
		return nil, fmt.Errorf("error preparing query GetIntegratorConfigByName: %w", err)
	}
	if q.getIntegratorConfigsStmt, err = db.PrepareContext(ctx, getIntegratorConfigs); err != nil {
		return nil, fmt.Errorf("error preparing query GetIntegratorConfigs: %w", err)
	}
	if q.getLogsStmt, err = db.PrepareContext(ctx, getLogs); err != nil {
		return nil, fmt.Errorf("error preparing query GetLogs: %w", err)
	}
	if q.getOATransactionStmt, err = db.PrepareContext(ctx, getOATransaction); err != nil {
		return nil, fmt.Errorf("error preparing query GetOATransaction: %w", err)
	}
	if q.getOATransactionsStmt, err = db.PrepareContext(ctx, getOATransactions); err != nil {
		return nil, fmt.Errorf("error preparing query GetOATransactions: %w", err)
	}
	if q.getOATransactionsCountStmt, err = db.PrepareContext(ctx, getOATransactionsCount); err != nil {
		return nil, fmt.Errorf("error preparing query GetOATransactionsCount: %w", err)
	}
	if q.getSnbConfigStmt, err = db.PrepareContext(ctx, getSnbConfig); err != nil {
		return nil, fmt.Errorf("error preparing query GetSnbConfig: %w", err)
	}
	if q.getSnbConfigByFacilityAndDeviceStmt, err = db.PrepareContext(ctx, getSnbConfigByFacilityAndDevice); err != nil {
		return nil, fmt.Errorf("error preparing query GetSnbConfigByFacilityAndDevice: %w", err)
	}
	if q.getUserStmt, err = db.PrepareContext(ctx, getUser); err != nil {
		return nil, fmt.Errorf("error preparing query GetUser: %w", err)
	}
	if q.updateIntegratorConfigStmt, err = db.PrepareContext(ctx, updateIntegratorConfig); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateIntegratorConfig: %w", err)
	}
	if q.updateOATransactionStmt, err = db.PrepareContext(ctx, updateOATransaction); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateOATransaction: %w", err)
	}
	if q.updateSnbConfigStmt, err = db.PrepareContext(ctx, updateSnbConfig); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateSnbConfig: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.countLogsStmt != nil {
		if cerr := q.countLogsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing countLogsStmt: %w", cerr)
		}
	}
	if q.createIntegratorConfigStmt != nil {
		if cerr := q.createIntegratorConfigStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createIntegratorConfigStmt: %w", cerr)
		}
	}
	if q.createLogStmt != nil {
		if cerr := q.createLogStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createLogStmt: %w", cerr)
		}
	}
	if q.createOATransactionStmt != nil {
		if cerr := q.createOATransactionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createOATransactionStmt: %w", cerr)
		}
	}
	if q.createSnbConfigStmt != nil {
		if cerr := q.createSnbConfigStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createSnbConfigStmt: %w", cerr)
		}
	}
	if q.createUserStmt != nil {
		if cerr := q.createUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUserStmt: %w", cerr)
		}
	}
	if q.deleteIntegratorConfigStmt != nil {
		if cerr := q.deleteIntegratorConfigStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteIntegratorConfigStmt: %w", cerr)
		}
	}
	if q.deleteSnbConfigStmt != nil {
		if cerr := q.deleteSnbConfigStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteSnbConfigStmt: %w", cerr)
		}
	}
	if q.deleteUserStmt != nil {
		if cerr := q.deleteUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteUserStmt: %w", cerr)
		}
	}
	if q.getAllSnbConfigStmt != nil {
		if cerr := q.getAllSnbConfigStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllSnbConfigStmt: %w", cerr)
		}
	}
	if q.getIntegratorConfigStmt != nil {
		if cerr := q.getIntegratorConfigStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getIntegratorConfigStmt: %w", cerr)
		}
	}
	if q.getIntegratorConfigByClientStmt != nil {
		if cerr := q.getIntegratorConfigByClientStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getIntegratorConfigByClientStmt: %w", cerr)
		}
	}
	if q.getIntegratorConfigByNameStmt != nil {
		if cerr := q.getIntegratorConfigByNameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getIntegratorConfigByNameStmt: %w", cerr)
		}
	}
	if q.getIntegratorConfigsStmt != nil {
		if cerr := q.getIntegratorConfigsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getIntegratorConfigsStmt: %w", cerr)
		}
	}
	if q.getLogsStmt != nil {
		if cerr := q.getLogsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getLogsStmt: %w", cerr)
		}
	}
	if q.getOATransactionStmt != nil {
		if cerr := q.getOATransactionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getOATransactionStmt: %w", cerr)
		}
	}
	if q.getOATransactionsStmt != nil {
		if cerr := q.getOATransactionsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getOATransactionsStmt: %w", cerr)
		}
	}
	if q.getOATransactionsCountStmt != nil {
		if cerr := q.getOATransactionsCountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getOATransactionsCountStmt: %w", cerr)
		}
	}
	if q.getSnbConfigStmt != nil {
		if cerr := q.getSnbConfigStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSnbConfigStmt: %w", cerr)
		}
	}
	if q.getSnbConfigByFacilityAndDeviceStmt != nil {
		if cerr := q.getSnbConfigByFacilityAndDeviceStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSnbConfigByFacilityAndDeviceStmt: %w", cerr)
		}
	}
	if q.getUserStmt != nil {
		if cerr := q.getUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserStmt: %w", cerr)
		}
	}
	if q.updateIntegratorConfigStmt != nil {
		if cerr := q.updateIntegratorConfigStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateIntegratorConfigStmt: %w", cerr)
		}
	}
	if q.updateOATransactionStmt != nil {
		if cerr := q.updateOATransactionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateOATransactionStmt: %w", cerr)
		}
	}
	if q.updateSnbConfigStmt != nil {
		if cerr := q.updateSnbConfigStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateSnbConfigStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                                  DBTX
	tx                                  *sql.Tx
	countLogsStmt                       *sql.Stmt
	createIntegratorConfigStmt          *sql.Stmt
	createLogStmt                       *sql.Stmt
	createOATransactionStmt             *sql.Stmt
	createSnbConfigStmt                 *sql.Stmt
	createUserStmt                      *sql.Stmt
	deleteIntegratorConfigStmt          *sql.Stmt
	deleteSnbConfigStmt                 *sql.Stmt
	deleteUserStmt                      *sql.Stmt
	getAllSnbConfigStmt                 *sql.Stmt
	getIntegratorConfigStmt             *sql.Stmt
	getIntegratorConfigByClientStmt     *sql.Stmt
	getIntegratorConfigByNameStmt       *sql.Stmt
	getIntegratorConfigsStmt            *sql.Stmt
	getLogsStmt                         *sql.Stmt
	getOATransactionStmt                *sql.Stmt
	getOATransactionsStmt               *sql.Stmt
	getOATransactionsCountStmt          *sql.Stmt
	getSnbConfigStmt                    *sql.Stmt
	getSnbConfigByFacilityAndDeviceStmt *sql.Stmt
	getUserStmt                         *sql.Stmt
	updateIntegratorConfigStmt          *sql.Stmt
	updateOATransactionStmt             *sql.Stmt
	updateSnbConfigStmt                 *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                                  tx,
		tx:                                  tx,
		countLogsStmt:                       q.countLogsStmt,
		createIntegratorConfigStmt:          q.createIntegratorConfigStmt,
		createLogStmt:                       q.createLogStmt,
		createOATransactionStmt:             q.createOATransactionStmt,
		createSnbConfigStmt:                 q.createSnbConfigStmt,
		createUserStmt:                      q.createUserStmt,
		deleteIntegratorConfigStmt:          q.deleteIntegratorConfigStmt,
		deleteSnbConfigStmt:                 q.deleteSnbConfigStmt,
		deleteUserStmt:                      q.deleteUserStmt,
		getAllSnbConfigStmt:                 q.getAllSnbConfigStmt,
		getIntegratorConfigStmt:             q.getIntegratorConfigStmt,
		getIntegratorConfigByClientStmt:     q.getIntegratorConfigByClientStmt,
		getIntegratorConfigByNameStmt:       q.getIntegratorConfigByNameStmt,
		getIntegratorConfigsStmt:            q.getIntegratorConfigsStmt,
		getLogsStmt:                         q.getLogsStmt,
		getOATransactionStmt:                q.getOATransactionStmt,
		getOATransactionsStmt:               q.getOATransactionsStmt,
		getOATransactionsCountStmt:          q.getOATransactionsCountStmt,
		getSnbConfigStmt:                    q.getSnbConfigStmt,
		getSnbConfigByFacilityAndDeviceStmt: q.getSnbConfigByFacilityAndDeviceStmt,
		getUserStmt:                         q.getUserStmt,
		updateIntegratorConfigStmt:          q.updateIntegratorConfigStmt,
		updateOATransactionStmt:             q.updateOATransactionStmt,
		updateSnbConfigStmt:                 q.updateSnbConfigStmt,
	}
}

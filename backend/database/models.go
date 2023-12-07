// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

type IntegratorConfig struct {
	ID                 uuid.UUID
	ClientID           sql.NullString
	ProviderID         sql.NullInt32
	Name               sql.NullString
	IntegratorName     sql.NullString
	SpID               sql.NullString
	PlazaIDMap         pqtype.NullRawMessage
	Extra              pqtype.NullRawMessage
	Url                sql.NullString
	InsecureSkipVerify sql.NullBool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type IntegratorTransaction struct {
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
}

type Log struct {
	ID        uuid.UUID
	Level     sql.NullString
	Message   sql.NullString
	Fields    pqtype.NullRawMessage
	CreatedAt time.Time
}

type OaTransaction struct {
	ID                    uuid.UUID
	Businesstransactionid string
	Lpn                   sql.NullString
	Customerid            sql.NullString
	Jobid                 sql.NullString
	Facility              sql.NullString
	Device                sql.NullString
	Extra                 pqtype.NullRawMessage
	EntryLane             sql.NullString
	ExitLane              sql.NullString
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type SnbConfig struct {
	ID       uuid.UUID
	Name     sql.NullString
	Username sql.NullString
	Password sql.NullString
	Endpoint sql.NullString
	Facility []string
	Device   []string
}

type User struct {
	ID          uuid.UUID
	Name        sql.NullString
	Username    sql.NullString
	Password    sql.NullString
	Permissions []string
}

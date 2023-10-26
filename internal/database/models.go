// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package database

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

type Log struct {
	ID     uuid.UUID
	Module string
	Info   string
	Extra  pqtype.NullRawMessage
}

type SnbConfig struct {
	ID       uuid.UUID
	Endpoint sql.NullString
	Facility sql.NullString
	Device   sql.NullString
}

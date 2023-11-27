package transactions

import (
	"time"
)

type LogData struct {
	ID        string         `json:"id"`
	Level     string         `json:"level"`
	Message   string         `json:"message"`
	Fields    map[string]any `json:"fields"`
	CreatedAt time.Time      `json:"createdAt"`
}

type OATransaction struct {
	ID                    string         `json:"id"`
	BusinessTransactionId string         `json:"businessTransactionId"`
	Lpn                   string         `json:"lpn"`
	Customerid            string         `json:"customerid"`
	Jobid                 string         `json:"jobid"`
	Facility              string         `json:"facility"`
	Device                string         `json:"device"`
	Extra                 map[string]any `json:"extra"`
	EntryLane             string         `json:"entryLane"`
	ExitLane              string         `json:"exitLane"`
	CreatedAt             time.Time      `json:"createdAt"`
	UpdatedAt             time.Time      `json:"updatedAt"`
}

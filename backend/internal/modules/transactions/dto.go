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
	BusinessTransactionID string         `json:"businessTransactionId"`
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

type IntegratorTransactions struct {
	Id                    string         `json:"id"`
	BusinessTransactionID string         `json:"businessTransactionId"`
	Lpn                   string         `json:"lpn"`
	IntegratorID          string         `json:"integratorId"`
	Status                string         `json:"status"`
	Amount                float64        `json:"amount"`
	Error                 string         `json:"error"`
	IntegratorName        string         `json:"integratorName"`
	Extra                 map[string]any `json:"extra"`
	TaxData               map[string]any `json:"taxData"`
	CreatedAt             time.Time      `json:"createdAt"`
	UpdatedAt             time.Time      `json:"updatedAt"`
}

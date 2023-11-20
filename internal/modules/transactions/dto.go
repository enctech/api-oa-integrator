package transactions

import (
	"time"
)

type LogResponse struct {
	ID        string         `json:"id"`
	Level     string         `json:"level"`
	Message   string         `json:"message"`
	Fields    map[string]any `json:"fields"`
	CreatedAt time.Time      `json:"createdAt"`
}

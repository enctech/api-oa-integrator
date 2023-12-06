package logger

import (
	"api-oa-integrator/internal/database"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/sqlc-dev/pqtype"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func CreateLogger() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.LowercaseLevelEncoder

	customCore := &CustomDatabaseCore{}
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(customCore),
		zap.DebugLevel,
	))

	return zap.Must(logger, nil)
}

type CustomDatabaseCore struct {
}

func (c *CustomDatabaseCore) Sync() error {
	// If your database core requires a Sync() method, implement it here
	return nil
}

func (c *CustomDatabaseCore) Write(p []byte) (n int, err error) {
	logEntry := string(p)
	fmt.Println(logEntry)
	output := map[string]interface{}{}
	if err := json.Unmarshal(p, &output); err != nil {
		panic(err)
	}
	outputCopy := make(map[string]any)
	for k, v := range output {
		if k == "level" || k == "msg" {
			continue
		}
		outputCopy[k] = v
	}

	go func() {
		jsonString, _ := json.Marshal(outputCopy)
		createdAt, _ := time.Parse("2006-01-02T15:04:05.999-0700", output["timestamp"].(string))
		db := database.D()
		if db != nil {
			fmt.Println("INSERTING LOG")
			result, err := database.New(db).CreateLog(context.Background(), database.CreateLogParams{
				Level:     sql.NullString{String: output["level"].(string), Valid: true},
				Message:   sql.NullString{String: output["msg"].(string), Valid: true},
				Fields:    pqtype.NullRawMessage{RawMessage: jsonString, Valid: true},
				CreatedAt: createdAt.UTC().Round(time.Microsecond),
			})
			if err != nil {
				fmt.Println(fmt.Sprintf("Error while creating log: %s", err.Error()))
			} else {
				id, _ := result.LastInsertId()
				fmt.Println(fmt.Sprintf("INSERTING LOG DONE NO ERROR %v", id))
			}
		}
	}()

	return len(p), nil
}

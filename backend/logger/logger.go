package logger

import (
	"api-oa-integrator/database"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/sqlc-dev/pqtype"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

func LogData(level string, msg string, fields map[string]interface{}) {
	out := map[string]any{
		"level":     level,
		"message":   msg,
		"fields":    fields,
		"timestamp": time.Now().UTC().Round(time.Microsecond),
	}
	logger := zap.L()

	for k, v := range out {
		logger = logger.With(zap.Any(k, v))
	}

	switch level {
	case "error":
		logger.Error(msg)
	case "info":
		logger.Info(msg)
	case "debug":
		logger.Debug(msg)
	case "warn":
		logger.Warn(msg)
	case "fatal":
		logger.Fatal(msg)
	case "panic":
		logger.Panic(msg)
	default:
		logger.Info(msg)
	}

	db := database.D()
	if db == nil {
		fmt.Println("db is nil")
		return
	}
	if fields == nil {
		fields = map[string]interface{}{
			"timestamp": time.Now().UTC().Round(time.Microsecond),
		}
	}
	jsonString, _ := json.Marshal(fields)
	_, err := database.New(db).CreateLog(context.Background(), database.CreateLogParams{
		Level:     sql.NullString{String: level, Valid: true},
		Message:   sql.NullString{String: msg, Valid: true},
		Fields:    pqtype.NullRawMessage{RawMessage: jsonString, Valid: true},
		CreatedAt: time.Now().UTC().Round(time.Microsecond),
	})
	if err != nil {
		fmt.Println("fail to create log", err)
	}
}

func CreateLogger() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.LowercaseLevelEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		zap.DebugLevel,
	))

	return zap.Must(logger, nil)
}

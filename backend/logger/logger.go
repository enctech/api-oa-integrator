package logger

import (
	"api-oa-integrator/database"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//
// =========================
// Configuration
// =========================
//

const (
	logQueueSize     = 10_000
	logBatchSize     = 100
	logFlushInterval = 500 * time.Millisecond
	logDBTimeout     = 3 * time.Second
)

//
// =========================
// Types
// =========================
//

type dbLogEntry struct {
	Level     string
	Message   string
	Fields    []byte
	CreatedAt time.Time
}

//
// =========================
// Globals
// =========================
//

var (
	logQueue = make(chan dbLogEntry, logQueueSize)
	db       *sql.DB // keep reference for graceful shutdown
)

//
// =========================
// Public API
// =========================
//

// Init must be called ONCE at startup
func Init(databaseInstance *sql.DB) {
	if databaseInstance == nil {
		fmt.Println("DB logger disabled: db is nil")
		return
	}

	db = databaseInstance
	startDBLogger(db)
	zap.ReplaceGlobals(CreateLogger())
}

// LogData is SAFE, FAST, and NON-BLOCKING
func LogData(level, msg string, fields map[string]interface{}) {
	now := time.Now().UTC().Round(time.Microsecond)

	if fields == nil {
		fields = map[string]interface{}{}
	}
	fields["timestamp"] = now

	// ---- zap logging (sync, fast) ----
	logger := zap.L().With(
		zap.String("level", level),
		zap.Any("fields", fields),
		zap.Time("timestamp", now),
	)

	switch level {
	case "error":
		logger.Error(msg)
	case "warn":
		logger.Warn(msg)
	case "debug":
		logger.Debug(msg)
	case "fatal":
		logger.Fatal(msg)
	case "panic":
		logger.Panic(msg)
	default:
		logger.Info(msg)
	}

	// ---- async DB logging ----
	data, err := json.Marshal(fields)
	if err != nil {
		return
	}

	select {
	case logQueue <- dbLogEntry{
		Level:     level,
		Message:   msg,
		Fields:    data,
		CreatedAt: now,
	}:
	default:
		// queue full → DROP LOG (protect app)
	}
}

//
// =========================
// Internal: DB Logger Worker
// =========================
//

func startDBLogger(db *sql.DB) {
	go func() {
		ticker := time.NewTicker(logFlushInterval)
		defer ticker.Stop()

		batch := make([]dbLogEntry, 0, logBatchSize)

		for {
			select {
			case entry := <-logQueue:
				batch = append(batch, entry)
				if len(batch) >= logBatchSize {
					flushLogs(db, batch)
					batch = batch[:0]
				}

			case <-ticker.C:
				if len(batch) > 0 {
					flushLogs(db, batch)
					batch = batch[:0]
				}
			}
		}
	}()
}

// flushLogs writes logs to DB in bulk using arrays
func flushLogs(db *sql.DB, logs []dbLogEntry) {
	if len(logs) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), logDBTimeout)
	defer cancel()

	q := database.New(db)

	levels := make([]string, 0, len(logs))
	messages := make([]string, 0, len(logs))
	fields := make([]json.RawMessage, 0, len(logs))
	timestamps := make([]time.Time, 0, len(logs))

	for _, l := range logs {
		levels = append(levels, l.Level)
		messages = append(messages, l.Message)
		fields = append(fields, l.Fields)
		timestamps = append(timestamps, l.CreatedAt)
	}

	// Bulk insert using SQLC CreateLogs (Postgres unnest with ordinality)
	_, err := q.CreateLogs(ctx, database.CreateLogsParams{
		Column1: levels,
		Column2: messages,
		Column3: fields,
		Column4: timestamps,
	})

	if err != nil {
		fmt.Println("failed to write batch logs:", err)
	}
}

//
// =========================
// Zap Logger Setup
// =========================
//

func CreateLogger() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.LowercaseLevelEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		zap.DebugLevel,
	)

	return zap.New(core)
}

//
// =========================
// Optional: Graceful Shutdown
// =========================
//

func Shutdown(ctx context.Context) {
	if db == nil {
		return
	}

	for {
		select {
		case entry := <-logQueue:
			flushLogs(db, []dbLogEntry{entry})
		case <-ctx.Done():
			return
		default:
			return
		}
	}
}

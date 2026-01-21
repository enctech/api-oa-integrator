package logger

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/sqlc-dev/pqtype"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type logEntry struct {
	level     string
	message   string
	fields    map[string]interface{}
	timestamp time.Time
}

type LogBatcher struct {
	buffer     []logEntry
	mu         sync.Mutex
	batchSize  int
	flushTimer *time.Timer
	flushDelay time.Duration
	db         *sql.DB
	stopCh     chan struct{}
	wg         sync.WaitGroup
}

var (
	batcher     *LogBatcher
	batcherOnce sync.Once
)

// InitBatcher initializes the log batcher (call once at startup)
func InitBatcher(db *sql.DB, batchSize int, flushDelay time.Duration) {
	batcherOnce.Do(func() {
		batcher = &LogBatcher{
			buffer:     make([]logEntry, 0, batchSize),
			batchSize:  batchSize,
			flushDelay: flushDelay,
			db:         db,
			stopCh:     make(chan struct{}),
		}
		batcher.start()
	})
}

// Shutdown flushes remaining logs and stops the batcher
func Shutdown() {
	if batcher != nil {
		close(batcher.stopCh)
		batcher.wg.Wait()
		batcher.flush()
	}
}

func (b *LogBatcher) start() {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		<-b.stopCh
	}()
}

func (b *LogBatcher) add(entry logEntry) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.buffer = append(b.buffer, entry)

	// Reset or create timer
	if b.flushTimer != nil {
		b.flushTimer.Stop()
	}
	b.flushTimer = time.AfterFunc(b.flushDelay, func() {
		b.flush()
	})

	// Flush if batch size reached
	if len(b.buffer) >= b.batchSize {
		if b.flushTimer != nil {
			b.flushTimer.Stop()
		}
		b.flush()
	}
}

func (b *LogBatcher) flush() {
	b.mu.Lock()
	if len(b.buffer) == 0 {
		b.mu.Unlock()
		return
	}

	// Copy buffer and clear it
	entries := make([]logEntry, len(b.buffer))
	copy(entries, b.buffer)
	b.buffer = b.buffer[:0]
	b.mu.Unlock()

	// Insert in batch
	if b.db != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := b.bulkInsert(ctx, entries)
		if err != nil {
			fmt.Printf("failed to batch insert logs: %v\n", err)
		}
	}
}

func (b *LogBatcher) bulkInsert(ctx context.Context, entries []logEntry) error {
	if len(entries) == 0 {
		return nil
	}

	// Build bulk insert query
	query := "INSERT INTO logs (level, message, fields, created_at) VALUES "
	values := []interface{}{}

	for i, entry := range entries {
		if i > 0 {
			query += ", "
		}

		paramOffset := i * 4
		query += fmt.Sprintf("($%d, $%d, $%d, $%d)",
			paramOffset+1, paramOffset+2, paramOffset+3, paramOffset+4)

		jsonString, _ := json.Marshal(entry.fields)
		values = append(values,
			sql.NullString{String: entry.level, Valid: true},
			sql.NullString{String: entry.message, Valid: true},
			pqtype.NullRawMessage{RawMessage: jsonString, Valid: true},
			entry.timestamp,
		)
	}

	_, err := b.db.ExecContext(ctx, query, values...)
	return err
}

func LogData(level string, msg string, fields map[string]interface{}) {
	// Log to stdout via zap
	out := map[string]any{
		"level":     level,
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

	// Add to batch for DB insertion
	if batcher != nil {
		if fields == nil {
			fields = map[string]interface{}{
				"timestamp": time.Now().UTC().Round(time.Microsecond),
			}
		}

		entry := logEntry{
			level:     level,
			message:   msg,
			fields:    fields,
			timestamp: time.Now().UTC().Round(time.Microsecond),
		}
		_ = entry // Temporarily disabled - logs table is too large, see PR #2
		// batcher.add(entry)
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

package logger

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

type LogCleaner struct {
	db        *sql.DB
	retention time.Duration
	interval  time.Duration
	stopCh    chan struct{}
	wg        sync.WaitGroup
}

var (
	cleaner     *LogCleaner
	cleanerOnce sync.Once
)

// InitCleaner initializes the log cleaner that runs on startup and periodically
// retention: how long to keep logs (e.g., 100 days)
// interval: how often to run cleanup (e.g., 12 hours for twice daily)
func InitCleaner(db *sql.DB, retention time.Duration, interval time.Duration) {
	cleanerOnce.Do(func() {
		cleaner = &LogCleaner{
			db:        db,
			retention: retention,
			interval:  interval,
			stopCh:    make(chan struct{}),
		}
		cleaner.start()
	})
}

func (c *LogCleaner) start() {
	// Run cleanup immediately on startup
	c.cleanup()

	// Start periodic cleanup
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		ticker := time.NewTicker(c.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.cleanup()
			case <-c.stopCh:
				return
			}
		}
	}()
}

func (c *LogCleaner) cleanup() {
	if c.db == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cutoff := time.Now().Add(-c.retention)

	result, err := c.db.ExecContext(ctx,
		"DELETE FROM logs WHERE created_at < $1", cutoff)

	if err != nil {
		fmt.Printf("log cleanup failed: %v\n", err)
		return
	}

	if rowsAffected, err := result.RowsAffected(); err == nil && rowsAffected > 0 {
		fmt.Printf("log cleanup: deleted %d old log entries\n", rowsAffected)
	}
}

// StopCleaner stops the periodic cleanup goroutine
func StopCleaner() {
	if cleaner != nil {
		close(cleaner.stopCh)
		cleaner.wg.Wait()
	}
}

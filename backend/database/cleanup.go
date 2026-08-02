package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

type DataCleaner struct {
	db        *sql.DB
	retention time.Duration
	interval  time.Duration
	stopCh    chan struct{}
	wg        sync.WaitGroup
}

var (
	cleaner     *DataCleaner
	cleanerOnce sync.Once
)

var tables = []string{
	"logs",
	"oa_transactions",
	"integrator_transactions",
}

// InitCleaner initializes the data cleaner that runs on startup and periodically
// retention: how long to keep data (e.g., 100 days)
// interval: how often to run cleanup (e.g., 12 hours for twice daily)
func InitCleaner(db *sql.DB, retention time.Duration, interval time.Duration) {
	cleanerOnce.Do(func() {
		cleaner = &DataCleaner{
			db:        db,
			retention: retention,
			interval:  interval,
			stopCh:    make(chan struct{}),
		}
		cleaner.start()
	})
}

func (c *DataCleaner) start() {
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

func (c *DataCleaner) cleanup() {
	if c.db == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cutoff := time.Now().Add(-c.retention)

	for _, table := range tables {
		result, err := c.db.ExecContext(ctx,
			fmt.Sprintf("DELETE FROM %s WHERE created_at < $1", table), cutoff)

		if err != nil {
			fmt.Printf("%s cleanup failed: %v\n", table, err)
			continue
		}

		if rowsAffected, err := result.RowsAffected(); err == nil && rowsAffected > 0 {
			fmt.Printf("%s cleanup: deleted %d old entries\n", table, rowsAffected)
		}
	}
}

// StopCleaner stops the periodic cleanup goroutine
func StopCleaner() {
	if cleaner != nil {
		close(cleaner.stopCh)
		cleaner.wg.Wait()
	}
}

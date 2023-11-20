package database

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"sync"
)

var (
	_globalMu sync.RWMutex
	_globalDb = &sql.DB{}
)

func InitDatabase() error {
	db, err := sql.Open("postgres", viper.GetString("database.url"))
	if err != nil {
		zap.L().Sugar().Errorf("failed open db %v", err)
		return err
	}

	err = db.Ping()
	if err != nil {
		zap.L().Sugar().Errorf("failed ping db %v", err)
		return err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		zap.L().Sugar().Errorf("failed postgres.WithInstance %v", err)
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./internal/database/migrations",
		"postgres", driver)
	if err != nil {
		zap.L().Sugar().Errorf("failed NewWithDatabaseInstance up %v", err)
		return err
	}
	err = m.Up()
	if err != nil {
		zap.L().Sugar().Errorf("failed migrate up %v", err)
	}
	_globalMu.Lock()
	_globalDb = db
	_globalMu.Unlock()
	return nil
}

func D() *sql.DB {
	_globalMu.RLock()
	l := _globalDb
	_globalMu.RUnlock()
	return l
}

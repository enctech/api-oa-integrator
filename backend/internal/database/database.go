package database

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"sync"
)

var lock = &sync.Mutex{}
var singleInstance *sql.DB

func InitDatabase() error {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			db, err := sql.Open("postgres", viper.GetString("database.url"))
			if err != nil {
				zap.L().Sugar().Errorf("failed open db %v", err)
				return err
			}
			loggerAdapter := zerologadapter.New(zerolog.New(os.Stdout))
			db = sqldblogger.OpenDriver(viper.GetString("database.url"), db.Driver(), loggerAdapter)
			err = db.Ping()
			if err != nil {
				zap.L().Sugar().With("URL", viper.GetString("database.url")).Errorf("failed ping db %v", err)
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
			singleInstance = db
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}
	return nil
}

func D() *sql.DB {
	return singleInstance
}

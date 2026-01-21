package database

import (
	"api-oa-integrator/tracing"
	"database/sql"
	"fmt"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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
				"file://"+viper.GetString("migrations"),
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

func TracedD() DBTX {
	return tracing.NewTracedDBTX(singleInstance)
}

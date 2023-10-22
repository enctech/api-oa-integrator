package utils

import (
	"api-oa-integrator/internal/database"
	"context"
	"github.com/sqlc-dev/pqtype"
	"go.uber.org/zap"
)

func LogToDb(module, info string, data []byte) error {
	txn, _ := database.D().Begin()
	_, err := database.New(database.D()).WithTx(txn).CreateLog(context.Background(), database.CreateLogParams{
		Module: module,
		Info:   info,
		Extra: pqtype.NullRawMessage{
			Valid:      true,
			RawMessage: data,
		},
	})
	if err != nil {
		zap.L().Sugar().Errorf("Error create log to db %v", err)
		return err
	}
	err = txn.Commit()
	if err != nil {
		zap.L().Sugar().Errorf("Error commit log txn to db %v", err)
		return err
	}
	return nil
}

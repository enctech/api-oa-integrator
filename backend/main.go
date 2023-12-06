package main

import (
	_ "api-oa-integrator/docs"
	"api-oa-integrator/internal"
	"api-oa-integrator/internal/database"
	"api-oa-integrator/logger"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/sqlc-dev/pqtype"
	"go.uber.org/zap"
	"strings"
	"time"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

//	@title			Swagger OA Integrator API
//	@version		1.0
//	@description	This is a server OA integrator.

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//
// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and JWT token.
func main() {
	log := logger.CreateLogger()
	fmt.Println(viper.GetString("database.url"))

	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(log)

	zap.ReplaceGlobals(log)
	err := database.InitDatabase()

	jsonString, _ := json.Marshal(map[string]any{})
	_, err = database.New(database.D()).CreateLog(context.Background(), database.CreateLogParams{
		Level:     sql.NullString{String: "info", Valid: true},
		Message:   sql.NullString{String: "TEST", Valid: true},
		Fields:    pqtype.NullRawMessage{RawMessage: jsonString, Valid: true},
		CreatedAt: time.Now().UTC().Round(time.Microsecond),
	})
	if err != nil {
		panic(fmt.Sprintf("init database error %v", err))
		return
	}

	internal.InitServer()
}

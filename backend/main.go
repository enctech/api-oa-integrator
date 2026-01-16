package main

import (
	"api-oa-integrator/database"
	_ "api-oa-integrator/docs"
	"api-oa-integrator/internal"
	"api-oa-integrator/logger"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func init() {
	// Define the --config flag
	pflag.String("config", "", "Path to the config file")
	pflag.String("migrations", "./database/migrations", "Path to the migrations folder")
	pflag.Parse()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Bind the --config flag to viper
	viper.BindPFlag("config", pflag.Lookup("config"))
	viper.BindPFlag("migrations", pflag.Lookup("migrations"))
	if configPath := viper.GetString("config"); configPath != "" {
		viper.SetConfigFile(configPath)
	}
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

//	@BasePath	/api/

// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//
// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and JWT token.
func main() {

	err := database.InitDatabase()

	if err != nil {
		panic(fmt.Sprintf("init database error %v", err))
		return
	}

	zapLogger := logger.CreateLogger()
	zap.ReplaceGlobals(zapLogger)

	// Initialize the database
	db := database.D()

	logger.InitBatcher(db, 50, 5*time.Second)

	fmt.Println(viper.GetString("database.url"))

	defer func(zapLogger *zap.Logger) {
		logger.Shutdown()
		_ = zapLogger.Sync()
	}(zapLogger)

	internal.InitServer()
}

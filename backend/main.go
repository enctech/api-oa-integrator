package main

import (
	"api-oa-integrator/database"
	_ "api-oa-integrator/docs"
	"api-oa-integrator/internal"
	"api-oa-integrator/logger"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
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
		panic(fmt.Sprintf("init database error: %v", err))
	}
	db := database.D()
	logger.Init(db)
	defer zap.L().Sync()

	// 2️⃣ Init Echo server
	e := internal.InitServer()

	// 3️⃣ Listen for shutdown signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit // wait for Ctrl+C or SIGTERM
	zap.L().Info("Shutting down server...")

	// 4️⃣ Graceful shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		zap.L().Error("server shutdown failed", zap.Error(err))
	}

	// 5️⃣ Flush pending logs
	logger.Shutdown(ctx)
	zap.L().Info("Server stopped gracefully")
}

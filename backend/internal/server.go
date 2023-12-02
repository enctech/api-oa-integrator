package internal

import (
	"api-oa-integrator/internal/modules/auth"
	"api-oa-integrator/internal/modules/config"
	"api-oa-integrator/internal/modules/health"
	"api-oa-integrator/internal/modules/oa"
	"api-oa-integrator/internal/modules/transactions"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
	"os"
)

func InitServer() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	health.InitController(e)
	oa.InitController(e)
	auth.InitController(e)
	config.InitController(e)
	transactions.InitController(e)
	if _, err := os.Stat("./cert/certificate.pem"); errors.Is(err, os.ErrNotExist) {
		e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", viper.GetString("app.port"))))
	} else {
		e.Logger.Fatal(e.StartTLS(fmt.Sprintf(":%v", viper.GetString("app.port")), "./cert/certificate.pem", "./cert/private-key.pem"))
	}

}

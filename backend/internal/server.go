package internal

import (
	"api-oa-integrator/internal/modules/auth"
	"api-oa-integrator/internal/modules/config"
	"api-oa-integrator/internal/modules/health"
	"api-oa-integrator/internal/modules/misc"
	"api-oa-integrator/internal/modules/oa"
	"api-oa-integrator/internal/modules/transactions"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func InitServer() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	g := e.Group("api/")
	g.GET("/swagger/*", echoSwagger.WrapHandler)
	misc.InitController(g)
	health.InitController(g)
	oa.InitController(g)
	auth.InitController(g)
	config.InitController(g)
	transactions.InitController(g)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", viper.GetString("app.port"))))
}

package internal

import (
	"api-oa-integrator/internal/modules/auth"
	"api-oa-integrator/internal/modules/config"
	"api-oa-integrator/internal/modules/oa"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

func InitServer() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "System is up and running!")
	})
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	oa.InitController(e)
	auth.InitController(e)
	config.InitController(e)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", viper.GetString("app.port"))))
}

package internal

import (
	"api-oa-integrator/internal/modules/auth"
	"api-oa-integrator/internal/modules/config"
	"api-oa-integrator/internal/modules/health"
	"api-oa-integrator/internal/modules/misc"
	"api-oa-integrator/internal/modules/oa"
	"api-oa-integrator/internal/modules/transactions"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"os"
)

func InitServer() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	misc.InitController(e)
	health.InitController(e)
	oa.InitController(e)
	auth.InitController(e)
	config.InitController(e)
	transactions.InitController(e)
	if _, err := os.Stat("./cert/certificate.pem"); errors.Is(err, os.ErrNotExist) {
		e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", viper.GetString("app.port"))))
	} else {
		cert, err := tls.LoadX509KeyPair("./cert/certificate.pem", "./cert/private-key.pem")
		if err != nil {
			e.Logger.Fatal(err)
		}
		e.TLSServer = &http.Server{
			Addr: fmt.Sprintf(":%v", viper.GetString("app.port")),
			TLSConfig: &tls.Config{
				Certificates:     []tls.Certificate{cert},
				MinVersion:       tls.VersionTLS12,
				CurvePreferences: []tls.CurveID{tls.CurveP256, tls.CurveP384, tls.CurveP521},
			},
		}
		e.Logger.Fatal(e.StartServer(e.TLSServer))
	}

}

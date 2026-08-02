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
	g := e.Group("")
	g.GET("/swagger/*", echoSwagger.WrapHandler)
	misc.InitController(g)
	health.InitController(g)
	oa.InitController(g)
	auth.InitController(g)
	config.InitController(g)
	transactions.InitController(g)
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
				CipherSuites: []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
					tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
					tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				},
			},
		}
		e.Logger.Fatal(e.StartServer(e.TLSServer))
	}
}

package config

import (
	"api-oa-integrator/internal/middlewares"
	"github.com/labstack/echo/v4"
	"net/http"
)

type controller struct {
}

func InitController(e *echo.Echo) {
	g := e.Group("config")
	c := controller{}
	g.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	g.Use(middlewares.GuardWithJWT())
	g.POST("/snb-config", c.createSnbConfig, middlewares.AdminOnlyMiddleware())
	g.GET("/snb-config", c.getSnBConfig, middlewares.AdminOnlyMiddleware())
}

// createSnbConfig godoc
//
//	@Summary		Create config for snb
//	@Description	Create configuration required for OA to works.
//	@Tags			config
//	@Accept			application/json
//	@Produce		application/json
//	@Param			request	body	SnbConfig	false	"Request Body"
//	@Security		Bearer
//	@Router			/config/snb-config [post]
func (con controller) createSnbConfig(c echo.Context) error {
	req := new(SnbConfig)
	err := c.Bind(req)
	user, err := createSnbConfig(c.Request().Context(), *req)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}
	return c.JSON(http.StatusCreated, user)
}

// getSnBConfig godoc
//
//	@Summary		Get config for snb
//	@Description	Get configuration required for OA to works.
//	@Tags			config
//	@Accept			application/json
//	@Produce		application/json
//	@Param			facility	query	string	true	"Facility"
//	@Param			device		query	string	true	"Device"
//	@Security		Bearer
//	@Router			/config/snb-config [get]
func (con controller) getSnBConfig(c echo.Context) error {
	user, err := getSnbConfig(c.Request().Context(), SnbConfig{
		Facility: c.QueryParam("facility"),
		Device:   c.QueryParam("device"),
	})
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}
	return c.JSON(http.StatusOK, user)
}

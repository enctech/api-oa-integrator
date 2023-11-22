package config

import (
	"api-oa-integrator/internal/middlewares"
	"github.com/google/uuid"
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
	//g.Use(middlewares.GuardWithJWT())
	g.POST("/snb-config", c.createSnbConfig, middlewares.AdminOnlyMiddleware())
	g.PUT("/snb-config", c.updateSnbConfig, middlewares.AdminOnlyMiddleware())
	g.GET("/snb-config", c.getAllSnBConfig, middlewares.AdminOnlyMiddleware())
	g.GET("/snb-config/:id", c.getSnBConfig, middlewares.AdminOnlyMiddleware())
	g.DELETE("/snb-config/:id", c.deleteSnbConfig, middlewares.AdminOnlyMiddleware())
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

// updateSnbConfig godoc
//
//	@Summary		Create config for snb
//	@Description	Create configuration required for OA to works.
//	@Tags			config
//	@Accept			application/json
//	@Produce		application/json
//	@Param			request	body	SnbConfig	false	"Request Body"
//	@Param			id		path	string		true	"Id"
//	@Security		Bearer
//	@Router			/config/snb-config/{id} [put]
func (con controller) updateSnbConfig(c echo.Context) error {
	req := new(SnbConfig)
	err := c.Bind(req)
	id := uuid.MustParse(c.Param("facility"))
	user, err := updateSnbConfig(c.Request().Context(), id, *req)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}
	return c.JSON(http.StatusCreated, user)
}

// getAllSnBConfig godoc
//
//	@Summary		Get all config for snb
//	@Description	Get configuration required for OA to works.
//	@Tags			config
//	@Accept			application/json
//	@Produce		application/json
//	@Security		Bearer
//	@Router			/config/snb-config [get]
func (con controller) getAllSnBConfig(c echo.Context) error {
	out, err := getAllSnbConfig(c.Request().Context())
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}
	return c.JSON(http.StatusOK, out)
}

// getSnBConfig godoc
//
//	@Summary		Get config for snb
//	@Description	Get configuration required for OA to works.
//	@Tags			config
//	@Accept			application/json
//	@Produce		application/json
//	@Param			id	path	string	true	"Id"
//	@Security		Bearer
//	@Router			/config/snb-config/{id} [get]
func (con controller) getSnBConfig(c echo.Context) error {
	id := uuid.MustParse(c.Param("id"))
	user, err := getSnbConfig(c.Request().Context(), id)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}
	return c.JSON(http.StatusOK, user)
}

// deleteSnbConfig godoc
//
//	@Summary		Get config for snb
//	@Description	Get configuration required for OA to works.
//	@Tags			config
//	@Accept			application/json
//	@Produce		application/json
//	@Param			id	path	string	true	"Id"
//	@Security		Bearer
//	@Router			/config/snb-config/{id} [delete]
func (con controller) deleteSnbConfig(c echo.Context) error {
	id := uuid.MustParse(c.Param("id"))
	err := deleteSnbConfig(c.Request().Context(), id)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}
	return c.JSON(http.StatusOK, "deleted")
}

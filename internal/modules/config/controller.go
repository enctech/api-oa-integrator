package config

import (
	"api-oa-integrator/internal/middlewares"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
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
	g.PUT("/snb-config/:id", c.updateSnbConfig, middlewares.AdminOnlyMiddleware())
	g.GET("/snb-config", c.getAllSnBConfig, middlewares.AdminOnlyMiddleware())
	g.GET("/snb-config/:id", c.getSnBConfig, middlewares.AdminOnlyMiddleware())
	g.DELETE("/snb-config/:id", c.deleteSnbConfig, middlewares.AdminOnlyMiddleware())

	g.GET("/integrator-config", c.getIntegratorConfigs, middlewares.AdminOnlyMiddleware())
	g.POST("/integrator-config", c.createIntegratorConfig, middlewares.AdminOnlyMiddleware())
	g.PUT("/integrator-config/:id", c.updateIntegratorConfig, middlewares.AdminOnlyMiddleware())
	g.DELETE("/integrator-config/:id", c.deleteIntegratorConfig, middlewares.AdminOnlyMiddleware())
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
//	@Summary		Update config for snb
//	@Description	Update configuration required for OA to works.
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
	id := uuid.MustParse(c.Param("id"))
	user, err := updateSnbConfig(c.Request().Context(), id, *req)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}
	return c.JSON(http.StatusCreated, user)
}

// getAllSnBConfig godoc
//
//	@Summary		Get all config for snb
//	@Description	Get all configuration required for OA to works.
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
	config, err := getSnbConfig(c.Request().Context(), id)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}
	return c.JSON(http.StatusOK, config)
}

// deleteSnbConfig godoc
//
//	@Summary		Delete config for snb
//	@Description	Delete configuration required for OA to works.
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

// getIntegratorConfigs godoc
//
//	@Summary		Get configs for all integrator
//	@Description	Get configurations for all integrators
//	@Tags			config
//	@Accept			application/json
//	@Produce		application/json
//	@Security		Bearer
//	@Router			/config/integrator-config [get]
func (con controller) getIntegratorConfigs(c echo.Context) error {
	req := new(IntegratorConfig)
	err := c.Bind(req)
	config, err := getIntegratorConfigs(c.Request().Context())
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}
	return c.JSON(http.StatusOK, config)
}

// createIntegratorConfig godoc
//
//	@Summary		Create config for integrator
//	@Description	Create configuration required for OA to send data to integrator.
//	@Tags			config
//	@Accept			application/json
//	@Produce		application/json
//	@Param			request	body	IntegratorConfig	false	"Request Body"
//	@Security		Bearer
//	@Router			/config/integrator-config [post]
func (con controller) createIntegratorConfig(c echo.Context) error {
	req := new(IntegratorConfig)
	err := c.Bind(req)
	configs, err := createIntegratorConfig(c.Request().Context(), *req)
	if err != nil {
		if strings.Contains(err.Error(), "violates unique constraint \"integrator_config_provider_id_key\"") {
			return c.String(http.StatusBadRequest, "Provider ID already exist")
		}

		if strings.Contains(err.Error(), "violates unique constraint \"integrator_config_name_key\"") {
			return c.String(http.StatusBadRequest, "Name already exist")
		}
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, configs)
}

// updateIntegratorConfig godoc
//
//	@Summary		Create config for integrator
//	@Description	Create configuration required for OA to send data to integrator.
//	@Tags			config
//	@Accept			application/json
//	@Produce		application/json
//	@Security		Bearer
//	@Param			id		path	string				true	"Id"
//	@Param			request	body	IntegratorConfig	false	"Request Body"
//	@Router			/config/integrator-config/{id} [put]
func (con controller) updateIntegratorConfig(c echo.Context) error {
	req := new(IntegratorConfig)
	err := c.Bind(req)
	id := uuid.MustParse(c.Param("id"))
	configs, err := updateIntegratorConfig(c.Request().Context(), id, *req)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}
	return c.JSON(http.StatusCreated, configs)
}

// deleteIntegratorConfig godoc
//
//	@Summary		Delete config for integrator
//	@Description	Create configuration required for OA to send data to integrator.
//	@Tags			config
//	@Accept			application/json
//	@Produce		application/json
//	@Security		Bearer
//	@Param			id		path	string				true	"Id"
//	@Router			/config/integrator-config/{id} [delete]
func (con controller) deleteIntegratorConfig(c echo.Context) error {
	id := uuid.MustParse(c.Param("id"))
	err := deleteIntegratorConfig(c.Request().Context(), id)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}
	return c.JSON(http.StatusOK, "deleted")
}

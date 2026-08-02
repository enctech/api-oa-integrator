package health

import (
	"api-oa-integrator/database"
	"api-oa-integrator/internal/modules/oa"
	"github.com/labstack/echo/v4"
	"net/http"
)

type controller struct {
}

func InitController(e *echo.Group) {
	g := e.Group("health")
	c := controller{}
	g.GET("", c.health)
}

// health godoc
//
//	@Summary		check system health
//	@Description	To check overall system health
//	@Tags			health
//	@Param			facility	query	int	false	"Facility"
//	@Param			device		query	int	false	"Device"
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/health [get]
func (con controller) health(c echo.Context) error {
	out := make(map[string]string)
	_, err := database.D().Exec("SELECT * FROM pg_catalog.pg_database")
	out["db"] = "up"
	if err != nil {
		out["db"] = "down"
	}
	err = oa.CheckSystemAvailability(c.QueryParam("facility"), c.QueryParam("device"))
	out["oa"] = "up"
	if err != nil {
		out["oa"] = "down"
	}
	return c.JSON(http.StatusCreated, out)
}

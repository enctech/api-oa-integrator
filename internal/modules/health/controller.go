package health

import (
	"api-oa-integrator/internal/database"
	"github.com/labstack/echo/v4"
	"net/http"
)

type controller struct {
}

func InitController(e *echo.Echo) {
	g := e.Group("health")
	c := controller{}
	g.GET("/", c.health)
}

type Response struct {
	Db string `json:"db"`
}

// health godoc
//
//	@Summary		check system health
//	@Description	To check overall system health
//	@Tags			health
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/health/ [get]
func (con controller) health(c echo.Context) error {
	out := Response{}
	_, err := database.D().Exec("SELECT * FROM pg_catalog.pg_database")
	dbStatus := "up"
	if err != nil {
		dbStatus = "down"
	}
	out.Db = dbStatus

	return c.JSON(http.StatusCreated, out)
}

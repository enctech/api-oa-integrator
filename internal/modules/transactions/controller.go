package transactions

import (
	"api-oa-integrator/internal/database"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type controller struct {
}

func InitController(e *echo.Echo) {
	g := e.Group("transactions")
	c := controller{}
	g.GET("", c.getLogs)
}

type Response struct {
	Db string `json:"db"`
}

// getLogs godoc
//
//	@Summary		get all logs
//	@Description	To check overall system health
//	@Param			before	query	string	true	"Before"	Format(dateTime)
//	@Param			after	query	string	false	"After"		Format(dateTime)
//	@Param			message	query	string	false	"Message"
//	@Param			fields	query	string	false	"Fields"
//	@Tags			transactions
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/transactions [get]
func (con controller) getLogs(c echo.Context) error {
	after, err := time.Parse(time.RFC3339, c.QueryParam("after"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid after date")
	}
	before, err := time.Parse(time.RFC3339, c.QueryParam("before"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid before date")
	}
	logs, err := database.New(database.D()).GetLog(c.Request().Context(), database.GetLogParams{
		After:   after.UTC().Round(time.Microsecond),
		Before:  before.UTC().Round(time.Microsecond),
		Message: fmt.Sprintf("%%%v%%", c.QueryParam("message")),
		Fields:  fmt.Sprintf("%%%v%%", c.QueryParam("fields")),
	})

	var out []LogResponse

	for _, log := range logs {
		var fields map[string]any
		if log.Fields.Valid {
			err = json.Unmarshal(log.Fields.RawMessage, &fields)
		}
		out = append(out, LogResponse{
			ID:        log.ID.String(),
			Level:     log.Level.String,
			Message:   log.Message.String,
			Fields:    fields,
			CreatedAt: log.CreatedAt,
		})
	}

	if out == nil {
		return c.JSON(http.StatusCreated, []LogResponse{})
	}

	return c.JSON(http.StatusCreated, out)
}

package transactions

import (
	"api-oa-integrator/internal/database"
	"api-oa-integrator/utils"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"math"
	"net/http"
	"strconv"
	"time"
)

type controller struct {
}

func InitController(e *echo.Echo) {
	g := e.Group("transactions")
	c := controller{}
	g.GET("/logs", c.getLogs)
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
//	@Param			perPage	query	int		false	"PerPage"
//	@Param			page	query	int		false	"Page"
//	@Tags			transactions
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/transactions/logs [get]
func (con controller) getLogs(c echo.Context) error {
	after, err := time.Parse(time.RFC3339, c.QueryParam("after"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid after date")
	}
	before, err := time.Parse(time.RFC3339, c.QueryParam("before"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid before date")
	}
	perPage := 50
	if c.QueryParam("perPage") != "" {
		perPage, err = strconv.Atoi(c.QueryParam("perPage"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, "invalid perPage")
		}
	}
	page := 0
	if c.QueryParam("page") != "" {
		page, err = strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, "invalid page")
		}
	}

	logs, err := database.New(database.D()).GetLogs(c.Request().Context(), database.GetLogsParams{
		After:   after.UTC().Round(time.Microsecond),
		Before:  before.UTC().Round(time.Microsecond),
		Message: fmt.Sprintf("%%%v%%", c.QueryParam("message")),
		Fields:  fmt.Sprintf("%%%v%%", c.QueryParam("fields")),
		Limit:   int32(perPage),
		Offset:  int32(page * perPage),
	})

	totalData, err := database.New(database.D()).CountLogs(c.Request().Context(), database.CountLogsParams{
		After:  after.UTC().Round(time.Microsecond),
		Before: before.UTC().Round(time.Microsecond),
	})

	var logOutput []LogData

	for _, log := range logs {
		var fields map[string]any
		if log.Fields.Valid {
			err = json.Unmarshal(log.Fields.RawMessage, &fields)
		}
		logOutput = append(logOutput, LogData{
			ID:        log.ID.String(),
			Level:     log.Level.String,
			Message:   log.Message.String,
			Fields:    fields,
			CreatedAt: log.CreatedAt,
		})
	}

	out := utils.PaginationResponse[LogData]{
		Data: logOutput,
		Metadata: utils.PaginationMetadata{
			TotalData: totalData,
			Page:      int64(page),
			PerPage:   int64(perPage),
			TotalPage: int64(math.Ceil(float64(totalData) / float64(perPage))),
		},
	}

	return c.JSON(http.StatusCreated, out)
}

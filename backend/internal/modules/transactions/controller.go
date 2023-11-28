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
	g.GET("/oa", c.getOATransaction)
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

// getOATransaction godoc
//
//	@Summary		get all OA logs
//	@Description	To get all transactions made through OA
//	@Param			startAt		query	string	false	"Start At"	Format(dateTime)
//	@Param			endAt		query	string	true	"End At"	Format(dateTime)
//	@Param			exitLane	query	string	false	"Exit Lane"
//	@Param			entryLane	query	string	false	"Entry Lane"
//	@Param			lpn			query	string	false	"Licence Plate Number"
//	@Param			facility	query	string	false	"Facility"
//	@Param			jobid		query	string	false	"Job ID"
//	@Param			perPage		query	int		false	"PerPage"
//	@Param			page		query	int		false	"Page"
//	@Tags			transactions
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/transactions/oa [get]
func (con controller) getOATransaction(c echo.Context) error {
	after, _ := time.Parse(time.RFC3339, c.QueryParam("startAt"))
	before, err := time.Parse(time.RFC3339, c.QueryParam("endAt"))
	if err != nil {
		before = time.Now()
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

	txns, err := database.New(database.D()).GetOATransactions(c.Request().Context(), database.GetOATransactionsParams{
		StartAt:   after.UTC().Round(time.Microsecond),
		EndAt:     before.UTC().Round(time.Microsecond),
		ExitLane:  c.QueryParam("exitLane"),
		EntryLane: c.QueryParam("entryLane"),
		Lpn:       c.QueryParam("lpn"),
		Facility:  c.QueryParam("facility"),
		Jobid:     c.QueryParam("jobid"),
	})

	totalData, err := database.New(database.D()).GetOATransactionsCount(c.Request().Context(), database.GetOATransactionsCountParams{
		StartAt:   after.UTC().Round(time.Microsecond),
		EndAt:     before.UTC().Round(time.Microsecond),
		ExitLane:  c.QueryParam("exitLane"),
		EntryLane: c.QueryParam("entryLane"),
		Lpn:       c.QueryParam("lpn"),
		Facility:  c.QueryParam("facility"),
		Jobid:     c.QueryParam("jobid"),
	})

	var txnData []OATransaction

	for _, log := range txns {
		var extra map[string]any
		if log.Extra.Valid {
			err = json.Unmarshal(log.Extra.RawMessage, &extra)
		}
		txnData = append(txnData, OATransaction{
			ID:                    log.ID.String(),
			BusinessTransactionId: log.Businesstransactionid,
			Lpn:                   log.Lpn.String,
			Customerid:            log.Customerid.String,
			Jobid:                 log.Jobid.String,
			Facility:              log.Facility.String,
			Device:                log.Device.String,
			Extra:                 extra,
			EntryLane:             log.EntryLane.String,
			ExitLane:              log.ExitLane.String,
			UpdatedAt:             log.UpdatedAt,
			CreatedAt:             log.CreatedAt,
		})
	}

	out := utils.PaginationResponse[OATransaction]{
		Data: txnData,
		Metadata: utils.PaginationMetadata{
			TotalData: totalData,
			Page:      int64(page),
			PerPage:   int64(perPage),
			TotalPage: int64(math.Ceil(float64(totalData) / float64(perPage))),
		},
	}

	return c.JSON(http.StatusCreated, out)
}

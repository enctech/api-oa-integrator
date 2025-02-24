package misc

import (
	"api-oa-integrator/database"
	"api-oa-integrator/internal/modules/oa"
	"api-oa-integrator/logger"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type controller struct {
}

func InitController(e *echo.Group) {
	g := e.Group("misc")
	c := controller{}
	g.GET("/", c.getData)
	g.GET("/integrator", c.getIntegratorStatus)
}

// getData godoc
//
//	@Summary		Get all misc data used in homepage
//	@Description	Get all misc data used in homepage
//	@Tags		 misc
//	@Accept			application/json
//	@Produce		application/json
//	@Param			startAt	query	string	false	"Start At"	Format(dateTime)
//	@Param			endAt	query	string	true	"End At"	Format(dateTime)
//	@Router			/misc/ [get]
func (con controller) getData(c echo.Context) error {
	after, _ := time.Parse(time.RFC3339, c.QueryParam("startAt"))
	before, err := time.Parse(time.RFC3339, c.QueryParam("endAt"))
	if err != nil {
		before = time.Now()
	}

	totalIn, _ := database.New(database.D()).GetOAEntryTransactions(c.Request().Context(), database.GetOAEntryTransactionsParams{
		StartAt: after.UTC().Round(time.Microsecond),
		EndAt:   before.UTC().Round(time.Microsecond),
	})
	totalOut, _ := database.New(database.D()).GetOAExitTransactions(c.Request().Context(), database.GetOAExitTransactionsParams{
		StartAt: after.UTC().Round(time.Microsecond),
		EndAt:   before.UTC().Round(time.Microsecond),
	})

	totalPayment, _ := database.New(database.D()).GetTotalTransactionAmount(c.Request().Context(), database.GetTotalTransactionAmountParams{
		Status:  "success",
		StartAt: after.UTC().Round(time.Microsecond),
		EndAt:   before.UTC().Round(time.Microsecond),
	})
	payment, err := strconv.ParseFloat(totalPayment, 64)
	out := map[string]any{
		"totalPayment": payment,
		"totalIn":      totalIn,
		"totalOut":     totalOut,
	}

	return c.JSON(http.StatusOK, out)
}

// getIntegratorStatus godoc
//
//	@Summary		Get all integrator status
//	@Description	Get all integrator status
//	@Tags		 misc
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/misc/integrator [get]
func (con controller) getIntegratorStatus(c echo.Context) error {
	snbStatus := getAllSnBStatus(c.Request().Context())
	integratorsStatus := getAllIntegratorStatus(c.Request().Context())
	out := map[string]any{
		"snb":         snbStatus,
		"integrators": integratorsStatus,
	}

	return c.JSON(http.StatusOK, out)
}

func getAllSnBStatus(ctx context.Context) []map[string]any {
	var out []map[string]any
	configs, err := database.New(database.D()).GetAllSnbConfig(ctx)
	if err != nil {
		return out
	}

	for _, config := range configs {
		if config.Facility == nil || len(config.Facility) == 0 || config.Device == nil || len(config.Device) == 0 {
			continue
		}
		oaStatus := "up"
		err := oa.CheckSystemAvailability(config.Facility[0], config.Device[0])

		if err != nil {
			oaStatus = "down"
		}

		for _, facility := range config.Facility {
			out = append(out, map[string]any{
				"facility": facility,
				"status":   oaStatus,
			})
		}
	}

	return out
}

func getAllIntegratorStatus(ctx context.Context) []map[string]any {
	var out []map[string]any
	configs, err := database.New(database.D()).GetIntegratorConfigs(ctx)
	if err != nil {
		return out
	}

	for _, config := range configs {
		if config.Url.String == "" {
			continue
		}
		integratorStatus := "up"
		err := ping(config.Url.String)

		if err != nil {
			logger.LogData("error", fmt.Sprintf("fail to ping %v %v", err, config.Url.String), nil)
			integratorStatus = "down"
		}

		out = append(out, map[string]any{
			"integrator": config.Name.String,
			"status":     integratorStatus,
		})
	}

	return out
}

func ping(domain string) error {
	var client = http.Client{
		Timeout:   time.Second * 10,
		Transport: &http.Transport{},
	}

	req, err := http.NewRequest("HEAD", domain, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

package oa

import (
	"api-oa-integrator/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"io"
	"net/http"
)

type controller struct {
}

func InitController(e *echo.Echo) {
	g := e.Group("oa")
	c := controller{}
	g.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	g.PUT("/version", c.version)
	g.PUT("/:facility/:device/:jobId/cancel", c.cancel)
	g.PUT("/:facility/:device/:jobId/finalmessage", c.finalMessage)
	g.PUT("/:facility/:device/:jobId/medialist", c.mediaList)
}

// version godoc
//
//	@Summary		check version
//	@Description	get the version and configuration available
//	@Tags			oa
//	@Accept			application/xml
//	@Produce		application/xml
//	@Param			request	body	VersionRequestWrapper	false	"Request Body"
//	@Router			/oa/version [put]
func (con controller) version(c echo.Context) error {
	go func() {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return
		}
		data, err := json.Marshal(map[string]any{
			"requestBody": fmt.Sprintf("%v", string(body)),
		})
		if err != nil {
			return
		}
		err = utils.LogToDb("OA", "Version", data)
		if err != nil {
			return
		}
	}()
	return c.XML(http.StatusCreated, VersionResponse{
		CustomerVersion: viper.GetString("app.version"),
		Configuration:   Configuration{SupportedFunctions: []string{"job", "version", "cancel"}},
	})
}

// cancel godoc
//
//	@Summary		Cancels a running job
//	@Description	This request cancels a running job on the 3rd party side. The job is identified by its resource /facility/device/jobid
//	@Tags			oa
//	@Accept			application/xml
//	@Produce		application/xml
//	@Param			facility	path	string				true	"Facility"
//	@Param			device		path	string				true	"Device"
//	@Param			jobId		path	string				true	"Job ID"
//	@Param			request		body	CancelJobWrapper	false	"Request Body"
//	@Router			/oa/{facility}/{device}/{jobId}/cancel [put]
func (con controller) cancel(c echo.Context) error {
	go func() {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return
		}
		data, err := json.Marshal(map[string]any{
			"facility":    c.Param("facility"),
			"device":      c.Param("device"),
			"jobId":       c.Param("jobId"),
			"requestBody": fmt.Sprintf("%v", string(body)),
		})
		if err != nil {
			return
		}
		err = utils.LogToDb("OA", "Cancel Job", data)
		if err != nil {
			return
		}
	}()
	return c.XML(http.StatusCreated, ConfirmationResponse{
		ConfirmationDetailStatus: "CANCEL_ACCEPTED",
		ConfirmationStatus:       "OK",
	})
}

// finalMessage godoc
//
//	@Summary		Receive Final Message from S&B
//	@Description	This request sends the last message for a job. The job is identified by its resources /facility/device/jobid
//	@Tags			oa
//	@Accept			application/xml
//	@Produce		application/xml
//	@Param			facility	path	string					true	"Facility"
//	@Param			device		path	string					true	"Device"
//	@Param			jobId		path	string					true	"Job ID"
//	@Param			request		body	FinalMessageSBWrapper	false	"Request Body"
//	@Router			/oa/{facility}/{device}/{jobId}/finalmessage [put]
func (con controller) finalMessage(c echo.Context) error {
	go func() {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return
		}
		data, err := json.Marshal(map[string]any{
			"facility":    c.Param("facility"),
			"device":      c.Param("device"),
			"jobId":       c.Param("jobId"),
			"requestBody": fmt.Sprintf("%v", string(body)),
		})
		if err != nil {
			return
		}
		err = utils.LogToDb("OA", "Final Message", data)
		if err != nil {
			return
		}
	}()
	return c.XML(http.StatusCreated, ConfirmationResponse{
		ConfirmationDetailStatus: "FINALMESSAGESB_RECEIVED",
		ConfirmationStatus:       "OK",
	})
}

// mediaList godoc
//
//	@Summary		Creates new media data
//	@Description	Creates new media data for an existing job and sends the required information as a <mediaData> element to the 3rd party system.
//					The mediaData element is described as a complex xml type within the job element.
//	@Tags			oa
//	@Accept			application/xml
//	@Produce		application/xml
//	@Param			facility	path	string				true	"Facility"
//	@Param			device		path	string				true	"Device"
//	@Param			jobId		path	string				true	"Job ID"
//	@Param			request		body	MediaDataWrapper	false	"Request Body"
//	@Router			/oa/{facility}/{device}/{jobId}/medialist [put]
func (con controller) mediaList(c echo.Context) error {
	go func() {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return
		}
		data, err := json.Marshal(map[string]any{
			"facility":    c.Param("facility"),
			"device":      c.Param("device"),
			"jobId":       c.Param("jobId"),
			"requestBody": fmt.Sprintf("%v", string(body)),
		})
		if err != nil {
			return
		}
		err = utils.LogToDb("OA", "Media List", data)
		if err != nil {
			return
		}
	}()
	return c.XML(http.StatusCreated, ConfirmationResponse{
		ConfirmationDetailStatus: "MEDIA_DATA_RECEIVED",
		ConfirmationStatus:       "OK",
	})
}

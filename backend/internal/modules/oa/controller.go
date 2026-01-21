package oa

import (
	"api-oa-integrator/logger"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

type controller struct {
}

func isValidUser(username, password string) bool {
	identifications := []struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		{
			Username: "aKr71dJ1~P_",
			Password: "cpVh19{3L92H",
		},
		{
			Username: "K6k6zm8JFW6X",
			Password: "3Xv7Kx8hNkIg",
		},
	}
	for _, id := range identifications {
		if username == id.Username && password == fmt.Sprintf(":%v", id.Password) {
			return true
		}
	}
	return false
}

func InitController(e *echo.Group) {
	g := e.Group("oa")
	c := controller{}
	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if isValidUser(username, password) {
			return true, nil
		}
		return false, nil
	}))
	g.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	g.PUT("/AuthorizationService3rdParty/version", c.version)
	g.PUT("/AuthorizationService3rdParty/:facility/:device/:jobId/cancel", c.cancel)
	g.PUT("/AuthorizationService3rdParty/:facility/:device/:jobId/finalmessage", c.finalMessage)
	g.POST("/AuthorizationService3rdParty/:facility/:device/:jobId/medialist", c.mediaList)
	g.POST("/AuthorizationService3rdParty/:facility/:device/:jobId/fake", c.fakeTest)
	g.POST("/AuthorizationService3rdParty/:facility/:device/:jobId", c.createJob)
}

// version godoc
//
//	@Summary		check version
//	@Description	get the version and configuration available
//	@Tags			oa
//	@Accept			application/xml
//	@Produce		application/xml
//	@Param			request	body	VersionRequestWrapper	false	"Request Body"
//	@Router			/oa/AuthorizationService3rdParty/version [put]
func (con controller) version(c echo.Context) error {
	go func() {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return
		}
		logger.LogData("info", "version", map[string]interface{}{
			"module":      "OK",
			"requestBody": string(body),
			"vendor":      c.Param("vendor"),
		})
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
//	@Router			/oa/AuthorizationService3rdParty/{facility}/{device}/{jobId}/cancel [put]
func (con controller) cancel(c echo.Context) error {
	go func() {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return
		}
		logger.LogData("info", "cancel", map[string]interface{}{
			"facility":    c.Param("facility"),
			"device":      c.Param("device"),
			"jobId":       c.Param("jobId"),
			"vendor":      c.Param("vendor"),
			"requestBody": fmt.Sprintf("%v", string(body)),
		})
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
//	@Router			/oa/AuthorizationService3rdParty/{facility}/{device}/{jobId}/finalmessage [put]
func (con controller) finalMessage(c echo.Context) error {
	go func() {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return
		}
		logger.LogData("info", "finalMessage", map[string]interface{}{
			"facility":    c.Param("facility"),
			"device":      c.Param("device"),
			"jobId":       c.Param("jobId"),
			"vendor":      c.Param("vendor"),
			"requestBody": fmt.Sprintf("%v", string(body)),
		})
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
//	@Router			/oa/AuthorizationService3rdParty/{facility}/{device}/{jobId}/medialist [post]
func (con controller) mediaList(c echo.Context) error {
	go func() {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return
		}
		logger.LogData("info", "mediaList", map[string]interface{}{
			"facility":    c.Param("facility"),
			"device":      c.Param("device"),
			"jobId":       c.Param("jobId"),
			"vendor":      c.Param("vendor"),
			"requestBody": fmt.Sprintf("%v", string(body)),
		})
	}()
	return c.XML(http.StatusCreated, ConfirmationResponse{
		ConfirmationDetailStatus: "MEDIA_DATA_RECEIVED",
		ConfirmationStatus:       "OK",
	})
}

// createJob godoc
//
//	@Summary		S&B creates new job
//	@Description	Creates new job and sends the required information as URI and <job> element to 3rd party system.
//	@Tags			oa
//	@Accept			application/xml
//	@Produce		application/xml
//	@Param			facility	path	string		true	"Facility"
//	@Param			device		path	string		true	"Device"
//	@Param			jobId		path	string		true	"Job ID"
//	@Param			request		body	JobWrapper	false	"Request Body"
//	@Router			/oa/AuthorizationService3rdParty/{facility}/{device}/{jobId} [post]
func (con controller) createJob(c echo.Context) error {
	rm := &RequestMetadata{
		facility: c.Param("facility"),
		device:   c.Param("device"),
		jobId:    c.Param("jobId"),
	}
	ctx := c.Request().Context()
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		go sendEmptyFinalMessage(ctx, rm)
		return c.XML(http.StatusCreated, ConfirmationResponse{
			ConfirmationDetailStatus: "JOB_CREATED",
			ConfirmationStatus:       "OK",
		})
	}
	logger.LogData("info", "createJob", map[string]interface{}{
		"facility":    c.Param("facility"),
		"device":      c.Param("device"),
		"jobId":       c.Param("jobId"),
		"vendor":      c.Param("vendor"),
		"requestBody": fmt.Sprintf("%v", string(body)),
	})

	req := new(Job)
	err = xml.Unmarshal(body, &req)

	if err != nil {
		go sendEmptyFinalMessage(ctx, rm)
		return c.XML(http.StatusCreated, ConfirmationResponse{
			ConfirmationDetailStatus: "JOB_CREATED",
			ConfirmationStatus:       "OK",
		})
	}
	handleIdentificationEntry(c, req, rm)
	handleLeaveLoopEntry(ctx, req, rm)
	handleIdentificationExit(ctx, req, rm)
	handlePaymentExit(ctx, req, rm)
	handleLeaveLoopExit(ctx, req, rm)
	if c.Response().Committed {
		return nil
	}
	return c.XML(http.StatusCreated, ConfirmationResponse{
		ConfirmationDetailStatus: "JOB_CREATED",
		ConfirmationStatus:       "OK",
	})
}

func (con controller) fakeTest(c echo.Context) error {
	rm := &RequestMetadata{
		facility: c.Param("facility"),
		device:   c.Param("device"),
		jobId:    c.Param("jobId"),
	}
	go sendEmptyFinalMessage(c.Request().Context(), rm)
	return c.XML(http.StatusCreated, ConfirmationResponse{
		ConfirmationDetailStatus: "JOB_CREATED",
		ConfirmationStatus:       "OK",
	})
}

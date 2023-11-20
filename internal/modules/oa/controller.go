package oa

import (
	"encoding/xml"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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
	g.PUT("/AuthorizationService3rdParty/:facility/:device/:jobId/cancel", c.cancel)
	g.PUT("/AuthorizationService3rdParty/:facility/:device/:jobId/finalmessage", c.finalMessage)
	g.POST("/AuthorizationService3rdParty/:facility/:device/:jobId/medialist", c.mediaList)
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
		zap.L().Sugar().
			With("module", "OK").
			With("requestBody", string(body)).
			Info("version")
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
		zap.L().Sugar().
			With("facility", c.Param("facility")).
			With("device", c.Param("device")).
			With("jobId", c.Param("jobId")).
			With("requestBody", fmt.Sprintf("%v", string(body))).
			Info("cancel")
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
		zap.L().Sugar().
			With("facility", c.Param("facility")).
			With("device", c.Param("device")).
			With("jobId", c.Param("jobId")).
			With("requestBody", fmt.Sprintf("%v", string(body))).
			Info("finalMessage")
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
		zap.L().Sugar().
			With("facility", c.Param("facility")).
			With("device", c.Param("device")).
			With("jobId", c.Param("jobId")).
			With("requestBody", fmt.Sprintf("%v", string(body))).
			Info("mediaList")
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
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		go sendEmptyFinalMessage(rm)
		return nil
	}
	go func() {
		zap.L().Sugar().
			With("facility", c.Param("facility")).
			With("device", c.Param("device")).
			With("jobId", c.Param("jobId")).
			With("requestBody", fmt.Sprintf("%v", string(body))).
			Info("createJob")
	}()

	req := new(Job)
	err = xml.Unmarshal(body, &req)

	if err != nil {
		go sendEmptyFinalMessage(rm)
		return c.XML(http.StatusCreated, ConfirmationResponse{
			ConfirmationDetailStatus: "JOB_CREATED",
			ConfirmationStatus:       "OK",
		})
	}
	handleIdentificationEntry(c, req, rm)
	handleLeaveLoopEntry(req, rm)
	handleIdentificationExit(req, rm)
	handlePaymentExit(req, rm)
	handleLeaveLoopExit(req, rm)
	if c.Response().Committed {
		return nil
	}
	return c.XML(http.StatusCreated, ConfirmationResponse{
		ConfirmationDetailStatus: "JOB_CREATED",
		ConfirmationStatus:       "OK",
	})
}

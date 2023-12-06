package utils

import (
	"api-oa-integrator/logger"
	"bytes"
	"io"
	"net/http"
)

type LoggingRoundTripper struct {
	Transport http.RoundTripper
}

func (lrt *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Capture the request body
	reqInfo := map[string]interface{}{
		"method": req.Method,
		"url":    req.URL.String(),
	}
	var requestBody bytes.Buffer
	if req.Body != nil {
		_, _ = io.Copy(&requestBody, req.Body)
		req.Body = io.NopCloser(bytes.NewReader(requestBody.Bytes()))
		reqInfo["body"] = requestBody.String()
	}

	logger.LogData("info", "HTTP Request", reqInfo)

	// Perform the actual HTTP request
	resp, err := lrt.Transport.RoundTrip(req)
	resInfo := map[string]interface{}{
		"method": req.Method,
		"url":    req.URL.String(),
	}
	if req.Body != nil {
		resInfo["request-body"] = requestBody.String()
	}
	if err != nil {
		resInfo["error"] = err
		logger.LogData("info", "HTTP Response", resInfo)
		return nil, err
	}

	resInfo["status"] = resp.Status

	// Capture the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		resInfo["response-body"] = string(responseBody)
	}

	logger.LogData("info", "HTTP Response", resInfo)

	// Create a new response object with the same status, headers, and body
	newResponse := &http.Response{
		Status:        resp.Status,
		StatusCode:    resp.StatusCode,
		Header:        resp.Header,
		Body:          io.NopCloser(bytes.NewReader(responseBody)),
		ContentLength: int64(len(responseBody)),
	}

	return newResponse, nil
}

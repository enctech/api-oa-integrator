package utils

import (
	"bytes"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type LoggingRoundTripper struct {
	Transport http.RoundTripper
}

func (lrt *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Capture the request body
	requestLogger := zap.L().Sugar().
		With("method", req.Method).
		With("url", req.URL.String())
	var requestBody bytes.Buffer
	if req.Body != nil {
		_, _ = io.Copy(&requestBody, req.Body)
		req.Body = io.NopCloser(bytes.NewReader(requestBody.Bytes()))
		requestLogger = requestLogger.With("body", requestBody.String())
	}

	requestLogger.Info("HTTP Request")

	// Perform the actual HTTP request
	resp, err := lrt.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	responseLogger := zap.L().Sugar().
		With("method", req.Method).
		With("url", req.URL.String()).
		With("status", resp.Status)

	// Capture the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		responseLogger = responseLogger.With("body", string(responseBody))
	}

	responseLogger.Info("HTTP Response")

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

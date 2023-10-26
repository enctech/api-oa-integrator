package utils

import (
	"go.uber.org/zap"
	"net/http"
)

type LoggingRoundTripper struct {
	Transport http.RoundTripper
}

func (lrt *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Log the request details before sending it
	zap.L().Sugar().Infof("Request: %s %s", req.Method, req.URL.String())

	// Send the request using the underlying transport
	resp, err := lrt.Transport.RoundTrip(req)

	// You can also log the response details if needed
	if resp != nil {
		zap.L().Sugar().Infof("Response Status: %s", resp.Status)
	}

	return resp, err
}

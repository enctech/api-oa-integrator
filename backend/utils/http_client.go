package utils

import (
	"crypto/tls"
	"net/http"
	"time"
)

var GlobalInsecureHttpClient = &http.Client{
	Transport: &LoggingRoundTripper{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			MaxIdleConns:        200,
			MaxIdleConnsPerHost: 50,
			MaxConnsPerHost:     100,
			IdleConnTimeout:     90 * time.Second,
		},
	},
}

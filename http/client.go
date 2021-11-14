package http

import (
	"crypto/tls"
	"net/http"
	"time"
)

const httpTimeoutDuration = 5 * time.Second

func NewHttpClient(tlsSkipVerify bool) *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: tlsSkipVerify}, //nolint:gosec
	}
	return &http.Client{
		Timeout:   httpTimeoutDuration,
		Transport: transport,
	}
}

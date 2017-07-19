package main

import (
	"net/http"
	"testing"
)

// Health Should Default to 200
func TestHealthCheckHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/health", nil)
	RunEndpointTest(t, req, StatusResponse200.Status, StatusResponse200.ToJSONString())
}

// Health should be able to be set
func TestSetHealthStatus(t *testing.T) {
	SetHealthStatus(http.StatusInternalServerError)
	req, _ := http.NewRequest("GET", "/health", nil)
	RunEndpointTest(t, req, StatusResponse500.Status, StatusResponse500.ToJSONString())
}

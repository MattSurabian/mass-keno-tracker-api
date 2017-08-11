package main

import (
	"github.com/mattsurabian/mass-keno-tracker/pkg/keno-tracker"
	"gopkg.in/gin-gonic/gin.v1"
	"log"
	"net/http"
	"sync"
)

var (
	status    = StatusResponse200.Status
	status_mu sync.RWMutex
)

func GetHealthStatus() int {
	status_mu.RLock()
	defer status_mu.RUnlock()
	return status
}

func SetHealthStatus(newStatus int) {
	status_mu.Lock()
	status = newStatus
	status_mu.Unlock()
}

var HealthEndpoint = &VersionedEndpoint{
	Handlers: map[int]func(*gin.Context){
		1: healthHandlerV1,
	},
}

func healthHandlerDefault(c *gin.Context) {
	s := GetStatusResponse(GetHealthStatus())
	c.IndentedJSON(s.Status, s)
}

func healthHandlerV1(c *gin.Context) {
	// TODO: A version specific health check
	healthHandlerDefault(c)
}

func healthCheck() {
	log.Printf("%s", keno_tracker.IsCacheVolatile())
	if keno_tracker.IsCacheVolatile() {
		// Scheduled maintenance like an etl process, temporary 503
		SetHealthStatus(http.StatusServiceUnavailable)
	} else {
		SetHealthStatus(StatusResponse200.Status)
	}
}

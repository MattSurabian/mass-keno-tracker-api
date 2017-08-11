package main

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

var MockEndpoint = &VersionedEndpoint{
	Handlers: map[int]func(*gin.Context){
		1: func(c *gin.Context) {
			c.IndentedJSON(StatusResponse200.Status, StatusResponse200)
		},
	},
}

func TestVersionedEndpointHandler(t *testing.T) {
	router := gin.Default()
	versionedRoutes := router.Group("/mock/v:version")
	{
		versionedRoutes.GET("/", VersionedEndpointHandler(MockEndpoint))
	}

	req, _ := http.NewRequest("GET", "/mock/v1/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, StatusResponse200.Status, w.Code)

	req, _ = http.NewRequest("GET", "/mock/v5/", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, StatusResponse501.Status, w.Code)
}

func RunEndpointTest(t *testing.T, request *http.Request, expectedStatus int, expectedResponse string) {
	w := httptest.NewRecorder()
	r := GetRouter()

	r.ServeHTTP(w, request)

	assert.Equal(t, expectedStatus, w.Code)
	if expectedResponse != "" {
		assert.Equal(t, string(expectedResponse), w.Body.String())
	}
}

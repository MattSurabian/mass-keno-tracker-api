package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"log"
	"net/http"
	"strconv"
)

type Endpoint interface {
	DefaultHandler(*gin.Context)
}

type VersionedEndpoint struct {
	Handlers map[int]func(*gin.Context)
}

func (ve *VersionedEndpoint) DefaultHandler(c *gin.Context) {
	// Default response for endpoints is 501
	s := http.StatusNotImplemented
	c.IndentedJSON(s, &StatusResponse{
		Status:     s,
		StatusText: http.StatusText(s),
	})
}

func VersionedEndpointHandler(endpoint *VersionedEndpoint) func(*gin.Context) {
	return func(c *gin.Context) {
		version, err := strconv.Atoi(c.Param("version"))
		if err != nil {
			log.Printf("Can't parse version string to int %s", c.Param("version"))
		}
		handler, ok := endpoint.Handlers[version]
		if ok {
			handler(c)
		} else {
			endpoint.DefaultHandler(c)
		}
	}
}

package main

import (
	"github.com/mattsurabian/mass-keno-tracker/pkg/keno-tracker"
	"gopkg.in/gin-gonic/gin.v1"
)

var TodaysEndpoint = &VersionedEndpoint{
	Handlers: map[int]func(*gin.Context){
		1: todaysHandlerV1,
	},
}

func todaysHandlerV1(c *gin.Context) {
	todaysManifest, err := keno_tracker.GetTodaysDraws()
	if err != nil {
		c.IndentedJSON(StatusResponse500.Status, StatusResponse500)
		return
	}
	c.IndentedJSON(StatusResponse200.Status, todaysManifest)
}

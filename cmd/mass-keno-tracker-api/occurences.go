package main

import (
	"github.com/mattsurabian/mass-keno-tracker/pkg/keno-tracker"
	"gopkg.in/gin-gonic/gin.v1"
)

var OccurencesEndpoint = &VersionedEndpoint{
	Handlers: map[int]func(*gin.Context){
		1: OccurencesHandlerV1,
	},
}

func OccurencesHandlerV1(c *gin.Context) {
	drawValue, _ := c.Params.Get("draw")
	manifest, err := keno_tracker.GetDrawOccurrencesByValue(drawValue)
	if err != nil {
		c.Error(err)
		c.IndentedJSON(StatusResponse500.Status, StatusResponse500)
		return
	}
	c.IndentedJSON(StatusResponse200.Status, manifest)
}

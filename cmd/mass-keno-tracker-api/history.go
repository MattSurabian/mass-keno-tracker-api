package main

import (
	"fmt"
	"github.com/mattsurabian/mass-keno-tracker/pkg/keno-tracker"
	"gopkg.in/gin-gonic/gin.v1"
)

var HistoryEndpoint = &VersionedEndpoint{
	Handlers: map[int]func(*gin.Context){
		1: historyHandlerV1,
	},
}

func historyHandlerV1(c *gin.Context) {
	month, _ := c.Params.Get("month")
	year, _ := c.Params.Get("year")
	day, _ := c.Params.Get("day")

	manifest, err := keno_tracker.GetDrawsByDate(fmt.Sprintf("%s-%s-%s", year, month, day))
	if err != nil {
		c.Error(err)
		c.IndentedJSON(StatusResponse500.Status, StatusResponse500)
		return
	}
	c.IndentedJSON(StatusResponse200.Status, manifest)
}

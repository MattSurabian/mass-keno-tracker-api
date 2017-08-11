package main

import (
	"encoding/json"
	"net/http"
)

var StatusResponse200 = GetStatusResponse(http.StatusOK)
var StatusResponse500 = GetStatusResponse(http.StatusInternalServerError)
var StatusResponse400 = GetStatusResponse(http.StatusBadRequest)
var StatusResponse501 = GetStatusResponse(http.StatusNotImplemented)

type StatusResponse struct {
	Status     int    `json:"status_code"`
	StatusText string `json:"status_text"`
}

func (sr *StatusResponse) ToJSONString() string {
	jsonString, _ := json.MarshalIndent(sr, "", "    ")
	return string(jsonString)
}

func GetStatusResponse(s int) *StatusResponse {
	return &StatusResponse{
		Status:     s,
		StatusText: http.StatusText(s),
	}
}

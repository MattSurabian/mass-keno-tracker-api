package keno_tracker_models

import (
	"encoding/json"
	"log"
	"strings"
	"time"
)

type KenoManifest struct {
	StartId string         `json:"min"`
	EndId   string         `json:"max"`
	Date    string         `json:"date,omitempty"`
	Draws   []DrawResponse `json:"draws"`
}

type DrawResponse struct {
	Id              string `json:"draw_id"`
	Date            string `json:"draw_date"`
	DayNormalizedId string `json:"day_normalized_id"`
	Value           string `json:"winning_num"`
	Bonus           string `json:"bonus"`
}

func (d *DrawResponse) String() string {
	serialized, err := json.Marshal(d)
	if err != nil {
		log.Printf("Error serializing draw")
		return "{}"
	}
	return string(serialized)
}

type DrawOccurences struct {
	Value string         `json:"winning_num"`
	Draws []DrawResponse `json:"draws"`
}

// The States API returns dates of the form YYYY-MM-DD and MM/DD/YYYY
// this helper method normalizes those values to YYYY-MM-DD
func DateStringNormalizer(dateString string) string {
	var normalizedDate time.Time
	var err error
	if strings.Contains(dateString, "/") {
		normalizedDate, err = time.Parse("01/02/2006", dateString)
		if err != nil {
			log.Printf("Error normalizing date string: %s %s| leaving as is", dateString, err)
			return dateString
		}
	} else {
		normalizedDate, err = time.Parse("2006-01-02", dateString)
		if err != nil {
			log.Printf("Error normalizing date string: %s", err)
			return dateString
		}
	}

	return normalizedDate.Format(time.RFC3339[:10])
}

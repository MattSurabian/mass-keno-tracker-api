package keno_tracker

import (
	"encoding/json"
	"fmt"
	"github.com/mattsurabian/mass-keno-tracker/pkg/keno-tracker-models"
	"github.com/mattsurabian/mass-keno-tracker/pkg/mass-state-lottery"
	"github.com/mattsurabian/mass-keno-tracker/pkg/redis-cache"
	"log"
	"strconv"
	"time"
)

var RedisCache = redis_cache.RedisCacheService{
	"mass-keno-tracker",
}

func drawIdCacheKey(drawId string) string {
	return fmt.Sprintf("draw.%s", drawId)
}

func drawValueCacheKey(drawValue string) string {
	return fmt.Sprintf("draw.%s", drawValue)
}

func dateCacheKey(date string) string {
	return fmt.Sprintf("date.%s", date)
}

func cacheDateManifest(manifest keno_tracker_models.KenoManifest, ttlSecs int) {
	RedisCache.SetObject(dateCacheKey(manifest.Date), &manifest, ttlSecs)
}

func cacheDrawById(drawId string, draw keno_tracker_models.DrawResponse) {
	RedisCache.SetObject(drawIdCacheKey(drawId), &draw, 0)
}

func cacheDrawOccurenceByValue(drawValue string, drawDateTime string) {
	RedisCache.AppendToSet(drawValueCacheKey(drawValue), drawDateTime)
}

func GetDrawOccurrencesByValue(drawValue string) (keno_tracker_models.DrawOccurences, error) {

	response := keno_tracker_models.DrawOccurences{
		Value: drawValue,
		Draws: []keno_tracker_models.DrawResponse{},
	}

	occurrences, err := RedisCache.GetSet(drawValueCacheKey(drawValue))

	if err != nil {
		// If we don't have a record of it for any reason, return empty draws
		return response, nil
	}

	for _, serializedDraw := range occurrences {
		var draw keno_tracker_models.DrawResponse
		err := json.Unmarshal([]byte(serializedDraw), &draw)
		if err != nil {
			continue
		}
		response.Draws = append(response.Draws, draw)
	}

	return response, nil
}

func GetDrawsByDate(date string) (keno_tracker_models.KenoManifest, error) {
	var manifest keno_tracker_models.KenoManifest
	normalizedDate := keno_tracker_models.DateStringNormalizer(date)

	if normalizedDate == todaysDateString() {
		todaysManifest, err := mass_state_lottery.GetTodaysManifest()
		if err != nil {
			log.Print(err)
			return todaysManifest, err
		}
		ProcessManifest(&todaysManifest)
	}

	err := RedisCache.GetObject(dateCacheKey(normalizedDate), &manifest)

	if err != nil {
		log.Print(err)
		return manifest, err
	}

	return manifest, nil
}

func GetTodaysDraws() (keno_tracker_models.KenoManifest, error) {
	return GetDrawsByDate(todaysDateString())
}

func todaysDateString() string {
	return time.Now().Format("2006-01-02")
}

func ProcessManifest(manifest *keno_tracker_models.KenoManifest) {
	var currDate string
	var normalizedDrawId int
	var id string
	var dayManifest keno_tracker_models.KenoManifest
	var dayManifests []keno_tracker_models.KenoManifest

	drawsLength := len(manifest.Draws)
	for index, manifestDraw := range manifest.Draws {
		if currDate == "" {
			currDate = manifestDraw.Date
			dayManifest.StartId = manifestDraw.Id
		}

		if currDate != manifestDraw.Date {
			normalizedDrawId = 0
			dayManifest.EndId = id
			dayManifests = append(dayManifests, dayManifest)

			currDate = manifestDraw.Date
			dayManifest.StartId = manifestDraw.Id
			dayManifest.Date = keno_tracker_models.DateStringNormalizer(manifestDraw.Date)
			dayManifest.Draws = []keno_tracker_models.DrawResponse{}
		}

		manifestDraw.Date = keno_tracker_models.DateStringNormalizer(manifestDraw.Date)
		manifestDraw.DayNormalizedId = strconv.Itoa(normalizedDrawId)

		cacheDrawById(manifestDraw.Id, manifestDraw)
		cacheDrawOccurenceByValue(manifestDraw.Value, manifestDraw.String())

		dayManifest.Draws = append(dayManifest.Draws, manifestDraw)
		id = manifestDraw.Id
		normalizedDrawId++

		if index == drawsLength-1 {
			dayManifest.EndId = id
			dayManifest.Date = keno_tracker_models.DateStringNormalizer(manifestDraw.Date)
			dayManifests = append(dayManifests, dayManifest)
		}
	}

	for _, manifest := range dayManifests {
		cacheDateManifest(manifest, 0)
	}

}

func MarkCacheAsVolatile() {
	RedisCache.SetString("cache-volatile", "true", 0)
}

func MarkCacheAsNonVolatile() {
	RedisCache.Bust("cache-volatile")
}

func IsCacheVolatile() bool {
	return RedisCache.Exists("cache-volatile")
}

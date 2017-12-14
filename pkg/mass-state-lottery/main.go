package mass_state_lottery

import (
	"fmt"
	"github.com/mattsurabian/mass-keno-tracker-api/pkg/keno-tracker-models"
	"github.com/mattsurabian/mass-keno-tracker-api/pkg/redis-cache"
	"log"
	"time"
)

var (
	stateManifestBaseUrl string = "http://www.masslottery.com/data/json/search/dailygames"
	monthManifestPrefix  string = "history/15"
	todaysManifestKey    string = "todays/15.json"
	historyManifestKey   string = "history/15-dates.json"
)

type HistoryManifest struct {
	DrawingDates []DrawingDates `json:"drawing_dates"`
}

type DrawingDates struct {
	Year       int16 `json:"year_id,string"`
	StartMonth int8  `json:"start_month,string"`
	StartDay   int8  `json:"start_day,string"`
	EndMonth   int8  `json:"end_month,string"`
	EndDay     int8  `json:"end_day,string"`
}

var RedisCache = redis_cache.RedisCacheService{
	"mass-state-lottery",
}

func GetMonthManifest(manifestKey string) (*keno_tracker_models.KenoManifest, error) {
	var manifest keno_tracker_models.KenoManifest
	manifestUrl := fmt.Sprintf("%s/%s/%s.json", stateManifestBaseUrl, monthManifestPrefix, manifestKey)
	err := get(manifestUrl, &manifest)
	now := time.Now()
	nowKey := fmt.Sprintf("%d%02d", now.Year(), now.Month())
	ttl := 0
	if nowKey == manifestKey {
		ttl = 30
	}
	RedisCache.SetObject(manifestUrl, manifest, ttl)
	return &manifest, err
}

func GetTodaysManifest() (keno_tracker_models.KenoManifest, error) {
	var manifest keno_tracker_models.KenoManifest
	manifestUrl := fmt.Sprintf("%s/%s", stateManifestBaseUrl, todaysManifestKey)
	err := get(manifestUrl, &manifest)
	RedisCache.SetObject(manifestUrl, manifest, 30)
	return manifest, err
}

func GetHistoryManifest() (*HistoryManifest, error) {
	var manifest HistoryManifest
	manifestUrl := fmt.Sprintf("%s/%s", stateManifestBaseUrl, historyManifestKey)
	err := get(manifestUrl, &manifest)
	RedisCache.SetObject(manifestUrl, manifest, 30)
	return &manifest, err
}

// Fetch from cache or remote if not in cache
func get(url string, target interface{}) (err error) {
	err = RedisCache.GetObject(url, target)
	if err != nil {
		log.Printf("Couldn't fetch from cache. Recovering by fetching %s from remote host instead...", url)
		err = getJson(url, target)
		if err != nil {
			log.Printf("Error fetching and or decoding manifest: %s %s", url, err)
			return
		}
	}
	return
}

package keno_tracker_etl

import (
	"fmt"
	"github.com/mattsurabian/mass-keno-tracker/pkg/keno-tracker"
	"github.com/mattsurabian/mass-keno-tracker/pkg/mass-state-lottery"
	"github.com/mattsurabian/mass-keno-tracker/pkg/redis-cache"
	"log"
	"sync"
	"time"
)

var RedisCache = &redis_cache.RedisCacheService{
	"mass-keno-tracker-etl",
}

func FetchAllHistoricData() {
	keno_tracker.MarkCacheAsVolatile()
	var startTime = time.Now()

	const workers = 50
	wg := new(sync.WaitGroup)
	in := make(chan string, workers)
	for k := 0; k < workers; k++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for manifestSlug := range in {
				RequestAndProcessMonthManifest(manifestSlug)
			}
		}()
	}

	historyManifest, err := mass_state_lottery.GetHistoryManifest()
	if err != nil {
		log.Fatalf("Error fetching history: %s", err)
	}

	lastIngest, _ := RedisCache.GetString("ingest-last")

	for _, history := range historyManifest.DrawingDates {
		for i := history.StartMonth; i <= history.EndMonth; i++ {
			manifestSlug := fmt.Sprintf(
				"%d%02d", history.Year, i,
			)
			if manifestSlug >= lastIngest {
				in <- manifestSlug
			} else {
				log.Printf("Already ingested: %s | Skipping...", manifestSlug)
			}
		}
	}

	close(in)
	wg.Wait()

	if len(historyManifest.DrawingDates) != 0 {
		RedisCache.SetString("ingest-last", fmt.Sprintf("%d%02d", historyManifest.DrawingDates[len(historyManifest.DrawingDates)-1].Year, historyManifest.DrawingDates[len(historyManifest.DrawingDates)-1].EndMonth), 0)

		log.Printf("Historic data sets from %d to %d loaded in: %f seconds",
			historyManifest.DrawingDates[0].Year,
			historyManifest.DrawingDates[len(historyManifest.DrawingDates)-1].Year,
			time.Since(startTime).Seconds(),
		)
	} else {
		log.Printf("The history manifest was empty. No data loaded.")
	}
	keno_tracker.MarkCacheAsNonVolatile()
}

func RequestAndProcessMonthManifest(manifestSlug string) {
	manifest, err := mass_state_lottery.GetMonthManifest(manifestSlug)
	if err != nil {
		log.Printf(
			"Error fetching manifest: %s | %s",
			manifestSlug,
			err,
		)
	}
	keno_tracker.ProcessManifest(manifest)
}

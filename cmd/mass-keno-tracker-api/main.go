package main

import (
	"fmt"
	"github.com/braintree/manners"
	"github.com/mattsurabian/mass-keno-tracker-api/pkg/keno-tracker-etl"
	"gopkg.in/gin-gonic/gin.v1"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// global config values
var GO_ENV string

func init() {
	go keno_tracker_etl.FetchAllHistoricData()
}

func main() {
	GO_ENV = os.Getenv("GO_ENV")
	if GO_ENV == "" {
		GO_ENV = "development"
	}

	if GO_ENV != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	httpAddr := os.Getenv("MASS_KENO_HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":8090"
	}

	log.Printf("HTTP Mass Keno Tracker service listening on %s", httpAddr)

	router := GetRouter()

	httpServer := manners.NewServer()
	httpServer.Addr = httpAddr
	httpServer.Handler = router

	healthCheckTicker := time.NewTicker(1 * time.Minute)
	etlTicker := time.NewTicker(12 * time.Hour)
	errChan := make(chan error, 10)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// Initialize health of service
		healthCheck()
		errChan <- httpServer.ListenAndServe()
	}()

	for {
		select {
		case <-healthCheckTicker.C:
			log.Printf("Performing self health check...")
			go healthCheck()
		case <-etlTicker.C:
			log.Printf("Scheduled update of historic data...")
			go keno_tracker_etl.FetchAllHistoricData()
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		case s := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exiting...", s))
			httpServer.BlockingClose()
			os.Exit(0)
		}
	}
}

func GetRouter() *gin.Engine {
	router := gin.Default()

	versionedRoutes := router.Group("/v:version")
	{
		versionedRoutes.GET("/health", VersionedEndpointHandler(HealthEndpoint))
		versionedRoutes.GET("/history/:year/:month/:day", VersionedEndpointHandler(HistoryEndpoint))
		versionedRoutes.GET("/todays", VersionedEndpointHandler(TodaysEndpoint))
		versionedRoutes.GET("/occurences/:draw", VersionedEndpointHandler(OccurencesEndpoint))
	}

	router.GET("/health", healthHandlerDefault)
	return router
}

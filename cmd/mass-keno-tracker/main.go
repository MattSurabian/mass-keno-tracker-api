package main

import (
	"fmt"
	"github.com/braintree/manners"
	"gopkg.in/gin-gonic/gin.v1"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// global config values
var MASS_KENO_TODAYS_DRAWS_URL, MASS_KENO_HISTORY_MANIFEST_URL, MASS_KENO_MONTHLY_MANIFEST_BASE_URL, GO_ENV string

func main() {
	GO_ENV = os.Getenv("GO_ENV")
	if GO_ENV == "" {
		GO_ENV = "development"
	}

	if GO_ENV != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	MASS_KENO_TODAYS_DRAWS_URL = os.Getenv("MASS_KENO_TODAYS_DRAWS_URL")
	if MASS_KENO_TODAYS_DRAWS_URL == "" {
		MASS_KENO_TODAYS_DRAWS_URL = "http://www.masslottery.com/data/json/search/dailygames/todays/15.json"
	}

	MASS_KENO_HISTORY_MANIFEST_URL = os.Getenv("MASS_KENO_HISTORY_MANIFEST_URL")
	if MASS_KENO_HISTORY_MANIFEST_URL == "" {
		MASS_KENO_HISTORY_MANIFEST_URL = "http://www.masslottery.com/data/json/search/dailygames/history/15-dates.json"
	}

	MASS_KENO_MONTHLY_MANIFEST_BASE_URL = os.Getenv("MASS_KENO_MONTHLY_MANIFEST_BASE_URL")
	if MASS_KENO_MONTHLY_MANIFEST_BASE_URL == "" {
		MASS_KENO_MONTHLY_MANIFEST_BASE_URL = "http://www.masslottery.com/data/json/search/dailygames/history/15/"
	}

	httpAddr := os.Getenv("MASS_KENO_HTTP_ADDR")
	if httpAddr == "" {
		// Defaulted in the docker file for discoverability
		log.Fatal("MASS_KENO_HTTP_ADDR must be set!")
	}
	log.Printf("HTTP Mass Keno Tracker service listening on %s", httpAddr)

	router := GetRouter()

	httpServer := manners.NewServer()
	httpServer.Addr = httpAddr
	httpServer.Handler = router

	errChan := make(chan error, 10)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		errChan <- httpServer.ListenAndServe()
	}()

	for {
		select {
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
	}

	router.GET("/health", healthHandlerDefault)
	return router
}

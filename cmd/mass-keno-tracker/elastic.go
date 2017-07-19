package main

import (
	"context"
	"fmt"
	elastic "gopkg.in/olivere/elastic.v5"
	"log"
)

var es_ctx context.Context
var es_client *elastic.Client

func setupEsClient() {
	var err error
	es_ctx = context.Background()

	// Create a client
	es_client, err = elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("%s", MASS_KENO_ES_ADDRESS)),
		elastic.SetBasicAuth(MASS_KENO_ES_USER, MASS_KENO_ES_PASS),
	)

	if err != nil {
		// Handle error
		log.Print(fmt.Sprintf("Error connecting to: %s", MASS_KENO_ES_ADDRESS))
		log.Fatal("ELASTICSEARCH ERROR | ", err, " | Check MASS_KENO_ES_ADDRESS, MASS_KENO_ES_USER, MASS_KENO_ES_PASS environment variables")
	}
}

func GetEsClient() *elastic.Client {
	if es_client == nil {
		setupEsClient()
	}
	return es_client
}
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type RelayStats struct {
	Requests                    map[string]int            `json:"requests"`
	ResponseCodes               map[string]int            `json:"response_codes"`
	ResponseCodesPerDomain      map[string]map[string]int `json:"response_codes_per_domain"`
	DeliveryCodes               map[string]int            `json:"delivery_codes"`
	DeliveryCodesPerDomain      map[string]map[string]int `json:"delivery_codes_per_domain"`
	Exceptions                  map[string]int            `json:"exceptions"`
	ExceptionsPerDomain         map[string]map[string]int `json:"exceptions_per_domain"`
	DeliveryExceptions          map[string]int            `json:"delivery_exceptions"`
	DeliveryExceptionsPerDomain map[string]map[string]int `json:"delivery_exceptions_per_domain"`
}

func collector() {
	for {
		logger.Debugf("Starting Collection")

		for _, relay := range relays.RelayList() {
			logger.Debugf("Starting Collector for %s", relay)
			go collectorWorker(relay)
		}

		time.Sleep(60 * time.Second)
	}
}

func collectorWorker(relay string) {
	logger.Debugf("Collector for %s starting", relay)

	// Get stats from relay
	resp, err := http.Get(fmt.Sprintf("https://%s/stats", relay))
	if err != nil {
		logger.Errorf("%s responded: %s", relay, err.Error())
		return
	}

	// Get response body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("could not get response body from %s: %s", relay, err.Error())
		return
	}
	logger.Tracef("%s responded: %s", relay, body)

	// Mashal JSON
	rs := RelayStats{}
	err = json.Unmarshal(body, &rs)
	if err != nil {
		logger.Errorf("could not unmarshal json from %s: %s", relay, err.Error())
		return
	}

	relays.Set(relay, &rs)


}
package main

import (
	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"net/http"
	"os"
	"strings"
	"sync"
)

var logger *loggo.Logger

type ThreadedRelayList struct {
	relays map[string]RelayStats
	Lock   sync.RWMutex
}

func (t *ThreadedRelayList) Set(relay string, r *RelayStats) {
	t.Lock.Lock()
	defer t.Lock.Unlock()

	t.relays[relay] = *r
}

func (t *ThreadedRelayList) RelayList() []string {
	t.Lock.RLock()
	defer t.Lock.RUnlock()

	var relayList []string
	for relay, _ := range t.relays {
		relayList = append(relayList, relay)
	}
	return relayList

}

func (t *ThreadedRelayList) Get() map[string]RelayStats {
	t.Lock.RLock()
	defer t.Lock.RUnlock()

	return t.relays
}

var relays ThreadedRelayList

func main() {
	// Collect Config
	cfg := CollectConfig()

	// Init Logging
	newLogger := loggo.GetLogger("main")
	logger = &newLogger

	err := loggo.ConfigureLoggers(cfg.LoggerConfig)
	if err != nil {
		logger.Errorf("Error configurting Logger: %s", err.Error())
		return
	}

	_, err = loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr))
	if err != nil {
		logger.Errorf("Error configurting Color Logger: %s", err.Error())
		return
	}

	// Init relay map
	relays.relays = make(map[string]RelayStats)

	relayList := strings.Split(cfg.Relays, ",")
	for _, r := range relayList {
		logger.Debugf("Adding relay %s", r)
		relays.Set(r, &RelayStats{})
	}

	// Start Collector
	go collector()

	// start web server
	http.HandleFunc("/", handler)
	logger.Errorf("%v", http.ListenAndServe(":9099", nil))
}

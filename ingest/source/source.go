package source

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sync"

	"github.com/rs/zerolog/log"
)

type ChampionData struct {
	Champions []string `json:"champions"`
	WinRates  []string `json:"winRates"`
}

type ProcessedCounters struct {
	Champion string
	Counters []string
}

func combineData(rawData ChampionData) []string {
	maxLength := len(rawData.Champions)
	var combined []string
	for i := 0; i < maxLength; i++ {
		combined = append(combined, rawData.Champions[i]+" "+rawData.WinRates[i])
	}
	return combined
}

func FetchSource(counterUrl string, champion string, wg *sync.WaitGroup, result chan<- *ProcessedCounters) {
	defer wg.Done()

	resp, err := http.Get(counterUrl + url.QueryEscape(champion))
	if err != nil {
		log.Error().Err(err).Msg("Error fetching champion: " + champion)
		result <- nil
		return
	}
	defer resp.Body.Close()

	var data ChampionData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Error().Err(err).Msg("Error decoding JSON from URL: " + champion)
		result <- nil
		return
	}

	var processed ProcessedCounters
	processed.Champion = champion
	processed.Counters = combineData(data)

	result <- &processed
}

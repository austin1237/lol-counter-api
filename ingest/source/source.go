package source

import (
	"encoding/json"
	"errors"
	"ingest/champions"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

func validateRawData(data ChampionData) error {
	if len(data.Champions) != len(data.WinRates) {
		return errors.New("champions and winRates arrays have different lengths")
	}

	// Create a map for faster lookup of championList.
	championMap := make(map[string]bool)
	for _, champion := range champions.Champions {
		championMap[champion] = true
	}

	// Check if each item in the Champions array is a valid champ.
	for _, champion := range data.Champions {
		champion = strings.ToLower(champion)
		if _, exists := championMap[champion]; !exists {
			return errors.New("champion not found in championList: " + champion)
		}
	}

	for _, winRate := range data.WinRates {
		parts := strings.Split(winRate, "%")
		// Check if the first part is a valid number.
		_, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return errors.New("invalid win rate format: " + winRate)
		}
	}

	return nil
}

func FetchSource(counterUrl string, champion string, wg *sync.WaitGroup, result chan<- *ProcessedCounters) {
	defer wg.Done()

	urlChampion := champion

	// The source site deviates for the following two champs in its url structure

	if urlChampion == "nunu & willump" {
		urlChampion = "nunu"
	}

	if urlChampion == "renata glasc" {
		urlChampion = "renata"
	}

	resp, err := http.Get(counterUrl + url.QueryEscape(urlChampion))
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

	err = validateRawData(data)
	if err != nil {
		log.Error().Err(err).Msg("Error validating raw data for " + champion)
		result <- nil
		return
	}
	var processed ProcessedCounters
	processed.Champion = champion
	processed.Counters = combineData(data)

	result <- &processed
}

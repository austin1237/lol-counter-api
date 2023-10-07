package source

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
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
		fmt.Println("Error fetching champion:", champion, "-", err)
		result <- nil
		return
	}
	defer resp.Body.Close()

	var data ChampionData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON from URL:", champion, "-", err)
		result <- nil
		return
	}

	var processed ProcessedCounters
	processed.Champion = champion
	processed.Counters = combineData(data)

	result <- &processed
}

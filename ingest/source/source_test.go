package source

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRawData(t *testing.T) {
	// Sample data for testing
	championDataValid := ChampionData{
		Champions: []string{"Darius", "Swain", "Sett"},
		WinRates:  []string{"50.5%", "60%", "45.2%"},
	}

	// Test case 1: Valid data with champion names "Darius," "Swain," and "Sett"
	assert.Nil(t, validateRawData(championDataValid))

	// Test case 2: Champions and WinRates arrays have different lengths
	championDataInvalid1 := ChampionData{
		Champions: []string{"Darius", "Swain"},
		WinRates:  []string{"50.5%", "60%", "45.2%"},
	}
	assert.Error(t, validateRawData(championDataInvalid1))

	// Test case 3: Invalid champion "Fake" in the data
	championDataInvalid2 := ChampionData{
		Champions: []string{"Darius", "Fake", "Swain"},
		WinRates:  []string{"50.5%", "60%", "45.2%"},
	}
	assert.Error(t, validateRawData(championDataInvalid2))

	// Test case 4: Invalid win rate format
	championDataInvalid3 := ChampionData{
		Champions: []string{"Darius", "Swain", "Sett"},
		WinRates:  []string{"50.5%", "Invalid", "45.2%"},
	}
	assert.Error(t, validateRawData(championDataInvalid3))
}

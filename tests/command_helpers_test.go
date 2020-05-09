package tests

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/weihesdlegend/Mew/commands"
	"testing"
)

func TestParseTicker(t *testing.T) {
	// test regular ticker input
	ticker := "QQQ"
	expectedTickers := []string{"QQQ"}
	parsedTickers, err := commands.ParseTicker(ticker)
	if err != nil {
		t.Fatal(err)
	}
	if !compareSlice(expectedTickers, parsedTickers) {
		t.Error("ticker parse error")
	}

	// test batch ticker input
	ticker = "@QQQ_MSFT_AAPL_FB_XIAOMI"
	expectedTickers = []string{"QQQ", "MSFT", "AAPL", "FB", "XIAOMI"}
	parsedTickers, err = commands.ParseTicker(ticker)
	if err != nil {
		t.Fatal(err)
	}
	if !compareSlice(expectedTickers, parsedTickers) {
		t.Error("ticker parse error")
	}

	// test error case
	ticker = "@@100"
	_, err = commands.ParseTicker(ticker)
	expectedErr := errors.New("ticker parsing error")
	if assert.Error(t, err, "an error was expected") {
		assert.Equal(t, err, expectedErr)
	}
}

// compare if two string slice are the same
func compareSlice(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

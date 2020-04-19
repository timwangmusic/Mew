package tests

import "testing"

func Quote(ticker string) string {
	return ticker
}

func TestQuote(t *testing.T) {
	if Quote("MSFT") != "MSFT" {
		t.Error("wrong ticker")
	}
}

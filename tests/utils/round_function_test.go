package utils

import (
	"github.com/stretchr/testify/assert"
	"github.com/weihesdlegend/Mew/utils"
	"math"
	"testing"
)

func TestRoundFunction(t *testing.T) {
	price := 100.55555555 // 8-digit decimal after dot

	// demonstrate float overflow for single-digit precision
	expectedOneDigit := 100.6
	assert.NotEqual(t, expectedOneDigit, math.Round(price/0.1)*0.1)

	// test keeping one decimal digit after dot
	rounded, _ := utils.Round(price, 0.1)
	if rounded != expectedOneDigit {
		t.Errorf("expected %f, got %f", expectedOneDigit, rounded)
	}

	// demonstrate float overflow for double-digit precision
	newPrice := 100.45999999
	expectedTwoDigits := 100.46
	assert.NotEqual(t, expectedTwoDigits, math.Round(newPrice/0.01)*0.01)

	// test keeping two decimal digit after dot
	rounded, _ = utils.Round(newPrice, 0.01)
	if rounded != expectedTwoDigits {
		t.Errorf("expected %.2f, got %.2f", expectedTwoDigits, rounded)
	}
}

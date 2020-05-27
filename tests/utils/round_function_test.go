package utils

import (
	"github.com/stretchr/testify/assert"
	"github.com/weihesdlegend/Mew/utils"
	"math"
	"testing"
)

func TestRoundFunction(t *testing.T) {
	price := 100.55555555 // 8-digit decimal after dot

	expectedOneDigit := 100.6

	// demonstrate that the previous method has deficiency
	assert.NotEqual(t, expectedOneDigit, math.Round(price/0.1)*0.1)

	// test keeping one decimal digit after dot
	rounded, _ := utils.Round(price, 0.1)
	if rounded != expectedOneDigit {
		t.Errorf("expected %f, got %f", expectedOneDigit, rounded)
	}

	// test keeping two decimal digit after dot
	expectedTwoDigit := 100.56

	rounded, _ = utils.Round(price, 0.01)
	if rounded != expectedTwoDigit {
		t.Errorf("expected %.2f, got %.2f", expectedTwoDigit, rounded)
	}
}

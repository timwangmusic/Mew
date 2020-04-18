package commands

import (
	"math"
)

type Util struct {
}

// TODO move this into util function
func (util Util) round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

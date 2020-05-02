package commands

import (
	"fmt"
	"math"

	"astuart.co/go-robinhood"
)

func round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

// About to place <market> <buy> for <10> shares of <QQQ> at <$200>
func OrderToString(opts robinhood.OrderOpts, ins robinhood.Instrument) string {
	ret := fmt.Sprintf("About to place %s %s for %d shares of %s at $%.2f",
		opts.Type.String(), opts.Side.String(),
		opts.Quantity, ins.OrderSymbol(), opts.Price)
	return ret
}

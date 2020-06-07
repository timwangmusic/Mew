package utils

import (
	"fmt"
	"math"
	"strconv"

	"github.com/coolboy/go-robinhood"
)

// round input to 2 digits
func Round(x, unit float64) (float64, error) {
	tmp := math.Round(x/unit) * unit
	resString := fmt.Sprintf("%.2f", tmp)
	res, err := strconv.ParseFloat(resString, 64)
	return res, err
}

// About to place <market> <buy> for <10> shares of <QQQ> at <$200>
func OrderToString(opts robinhood.OrderOpts, ins robinhood.Instrument, finalPrice float64) string {
	ret := fmt.Sprintf("About to place %s %s for %d shares of %s at $%.2f",
		opts.Type.String(), opts.Side.String(),
		opts.Quantity, ins.OrderSymbol(), finalPrice)
	return ret
}

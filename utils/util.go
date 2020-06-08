package utils

import (
	"bufio"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"reflect"
	"strconv"
	"strings"

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
func OrderToString(opts robinhood.OrderOpts, ins robinhood.Instrument) string {
	ret := fmt.Sprintf("About to place %s %s for %d shares of %s at $%.2f",
		opts.Type.String(), opts.Side.String(),
		opts.Quantity, ins.OrderSymbol(), opts.Price)
	return ret
}

func ReadUserConfirmation(bufferReader *bufio.Reader) error {
	// to simplify testing
	if reflect.ValueOf(bufferReader).IsNil() {
		log.Debug("the buffer reader is nil")
		return nil
	}
	// wait for user confirmation
	for {
		text, _ := bufferReader.ReadString('\n')
		if strings.Contains(text, "y") {
			break
		} else {
			return errors.New("the order is cancelled")
		}
	}
	return nil
}

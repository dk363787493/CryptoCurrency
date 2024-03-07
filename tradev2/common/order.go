package common

import (
	"CryptoCurrency/tradev2/exchange"
	"github.com/tidwall/gjson"
	"strconv"
)

const (
	SideSell = "SELL"
	SideBuy  = "BUY"
)
const (
	OrderTypeLimit  = "LIMIT"
	OrderTypeMarket = "MARKET"
)

func ParsePirceData(d string) exchange.Order {
	r := gjson.Parse(d)
	ts, _ := strconv.ParseInt(r.Get("s").String(), 10, 64)
	price, _ := strconv.ParseFloat(r.Get("k.o").String(), 64)
	p := exchange.Order{
		TimeStamp: ts,
		Symbol:    r.Get("s").String(),
		Price:     price,
	}
	return p
}

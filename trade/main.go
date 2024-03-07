package main

import (
	"CryptoCurrency/common"
	"CryptoCurrency/trade/market"
	"CryptoCurrency/trade/strategy"
	"CryptoCurrency/trade/symboltrade"
	"fmt"
	"github.com/tidwall/gjson"
	"strconv"
	"time"
)

func main() {
	currencys := []string{"SHIBUSDT"}

	for _, c := range currencys {
		go func(currency string, shortTermWindow int, longTermWindow int) {
			s := new(strategy.MoveAvgStrategy)
			s.ShortTermWindow = shortTermWindow
			s.LongTermWindow = longTermWindow
			trade := symboltrade.NewTrade(currency, time.Minute)
			s.SellCh = trade.SellCh
			trade.AddStrategy(s)
			data := make(chan string, 1024)
			signCh := market.WsKline(common.OneS, currency, data)
			timer := time.NewTimer(5 * time.Second)
			//trade.Exceute()
			for {
				select {
				case <-signCh:
					goto r
				case <-timer.C:
					{
						d := <-data
						klineData := ParsePirceData(d)
						fmt.Printf("data:%+v\n", klineData)
						trade.AppendData(klineData)
						trade.ExecuteStrategy()
						timer.Reset(5 * time.Second)
					}
				default:
				}

			}

		r:
			fmt.Println("exiting....")
		}(c, 3, 5)
	}

	time.Sleep(2 * time.Minute)
}

func ParsePirceData(d string) *strategy.PriceData {
	r := gjson.Parse(d)
	ts, _ := strconv.ParseInt(r.Get("s").String(), 10, 64)
	price, _ := strconv.ParseFloat(r.Get("k.o").String(), 64)
	p := strategy.PriceData{
		Timestamp: ts,
		Symbol:    r.Get("s").String(),
		Price:     price,
	}
	return &p
}

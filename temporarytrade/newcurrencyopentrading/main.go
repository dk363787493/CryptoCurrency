package main

import (
	"CryptoCurrency/model"
	"CryptoCurrency/temporarytrade/newcurrencyopentrading/trade"
	"fmt"
	"time"
)

func main() {

	//timer := time.NewTimer(1 * time.Minute)
	timer := time.NewTimer(1 * time.Second)
	// Get all the currencys which will be open trading soon from mysql
	for {
		<-timer.C
		currency := model.GetAllOpenTradingCurrency("OPEN")
		for _, c := range currency {
			fmt.Printf("currency:%s%s\n", c.CrytoCurrency, c.BaseCurrency)
			go func(o model.OpenTradingCurrency) {
				model.UpdadteStatus(o.CrytoCurrency, o.BaseCurrency, "CLOSE")
				sub := o.StatustartmissionDate.Sub(time.Now())
				if sub > 0 {
					sleepS := time.Duration(int64(sub.Seconds()) - 1)
					time.Sleep(sleepS * time.Second)
				}

				symbol := fmt.Sprint(o.CrytoCurrency, o.BaseCurrency)
				order := &trade.Order{
					Symbol:        symbol,
					QuoteOrderQty: o.QuoteOrderQty,
				}
				var buy *trade.Order
				var err error
				for i := 0; i < 25; i++ {
					buy, err = trade.Buy(order)
					//buy, err = trade.MockBuy(order)
					if err != nil {
						fmt.Printf("currency:%s,err:%+v\n", order.Symbol, err)
						continue
					}
					break
				}
				if err != nil {
					return
				}
				fmt.Printf("buy order:%+v", buy)
				defer model.UpdadteStatus(o.CrytoCurrency, o.BaseCurrency, "CLOSE")
				// defer buy
				var steopSize = 0.1000
				_ = steopSize
				buy.Price = buy.Price * (1 + o.ProfitRate)
				sell, err := trade.Sell(buy, steopSize)
				//sell, err := trade.MockSell(buy)
				sellResult := true
				if err != nil {
					fmt.Printf("err:%s\n", err.Error())
					steopSize = 1.0000
					sell, err = trade.Sell(buy, steopSize)
					//sell, err = trade.MockSell(buy)
					if err != nil {
						fmt.Printf("err:%s\n", err.Error())
						sellResult = false
					}
				}
				fmt.Println("sellResult:", sellResult)
				if sellResult {
					// 入库  sell
					fmt.Printf("sell order:%+v", sell)
				} else {
					fmt.Printf("failed selling order:%+v\n", sell)
				}
			}(c)
		}
		timer.Reset(1 * time.Minute)
	}
	//currency := model.GetAllOpenTradingCurrency("OPEN")

	//Keep making a buy order for every currency mentioned above

	//If success,make a limit sell order for the currency right now

	//send a message for every trading to me
}

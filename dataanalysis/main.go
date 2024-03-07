package main

import (
	"CryptoCurrency/common"
	currency "CryptoCurrency/dataanalysis/currencystatistics"
	"CryptoCurrency/model"
	"errors"
	"fmt"
	"time"
)

func main() {
	binaceCh := make(chan currency.KlineData, 1024)
	okxCh := make(chan currency.KlineData, 1024)

	// Get currency from mysql
	exchangeBit := common.BinanceBit | common.OKxBit
	cs := model.GetAllCurrencyInfo(exchangeBit)
	symbols := make([]string, len(cs))
	for i, c := range cs {
		symbols[i] = c.Currency
	}
	//set the currey to get from Kline for analyzing
	binanceExchange := currency.NewBinanceCurrency(binaceCh)
	okxExchange := currency.NewOkxCurrency(okxCh)
	okxExchange.GetMarketDataStreamFromExchange(symbols)
	binanceExchange.GetMarketDataStreamFromExchange(symbols)
	// update newest data for currency
	binanceExchange.AsyncUpdateSymbolPrice()
	okxExchange.AsyncUpdateSymbolPrice()

	// compare the prices of the same currey from different exchange

	go func() {
		for {
			for _, s := range symbols {
				price1, err1 := binanceExchange.GetPrice(s)
				price2, err2 := okxExchange.GetPrice(s)
				if err := errors.Join(err1, err2); err != nil {
					fmt.Printf("err:%v\n", err)
					time.Sleep(time.Second)
					continue
				}
				diff := price1 - price2
				basePirce := price2
				//if diff < 0 {
				//	diff = -diff
				//	basePirce = price2
				//}
				diffRate := diff / basePirce
				//fmt.Printf("diffRate:%f\n", diffRate)
				//fmt.Printf("diff:%f\n", diff)
				diffRate = diffRate * 10000
				ss := fmt.Sprintf("symbol:%s,diffRate:%f", s, diffRate) + "‱"
				fmt.Println(ss)
				//if diffRate >= 10 {
				//	ss := fmt.Sprintf("symbol:%s,diffRate:%f", s, diffRate) + "‱"
				//	fmt.Println(ss)
				//}
			}
		}

	}()

	// set the limit of the rate for trading,considering the fee from all of the exchange

	time.Sleep(200 * time.Second)
}

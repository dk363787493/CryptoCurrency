package main

import (
	"fmt"
	binance_connector "github.com/binance/binance-connector-go"
)

func main() {
	// 获取所有交易币对信息
	//baseURL := "https://api.binance.com"
	//
	//client := binance_connector.NewClient("", "", baseURL)

	// ExchangeInfo
	//exchangeInfo, err := client.NewExchangeInfoService().Do(context.Background())
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(binance_connector.PrettyPrint(exchangeInfo))

	//r := gjson.Parse(binance_connector.PrettyPrint(exchangeInfo))
	//symbols := r.Get("symbols").Array()
	//for _, s := range symbols {
	//	quoteAsset := s.Get("quoteAsset").String()
	//	if quoteAsset == "USDT" {
	//		baseAsset := s.Get("baseAsset").String()
	//		fmt.Println("baseAsset:", baseAsset)
	//	}
	//}
	// Websocket Stream   实时获取币对价格
	websocketStreamClient := binance_connector.NewWebsocketStreamClient(false)
	wsKlineHandler := func(event *binance_connector.WsKlineEvent) {
		fmt.Println(binance_connector.PrettyPrint(event))
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneCh, _, err := websocketStreamClient.WsKlineServe("BTCUSDT", "1s", wsKlineHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, _, err1 := websocketStreamClient.WsKlineServe("ETHUSDT", "1s", wsKlineHandler, errHandler)
	if err1 != nil {
		fmt.Println(err)
		return
	}

	<-doneCh
}

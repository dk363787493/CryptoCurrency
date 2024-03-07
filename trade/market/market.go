package market

import (
	"fmt"
	binance_connector "github.com/binance/binance-connector-go"
)

//使用窗口均值来计算买入时机，目前设置盈利1%-2%就卖出，超短线买卖

func WsKline(interval string, symbol string, dataCh chan<- string) <-chan struct{} {
	websocketStreamClient := binance_connector.NewWebsocketStreamClient(false)
	wsKlineHandler := func(event *binance_connector.WsKlineEvent) {
		dataCh <- binance_connector.PrettyPrint(event)
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneCh, _, err := websocketStreamClient.WsKlineServe(symbol, interval, wsKlineHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return doneCh
}

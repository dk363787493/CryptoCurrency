package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	okxApiBaseUrl = "https://www.okx.com/api/v5/%s"
)

func main() {
	WsGetMarketData()
}

func RestGetMarketData() {
	///api/v5/market/candles?instId=BTC-USDT
	//url := fmt.Sprintf(okxApiBaseUrl, "public/instruments?instType=SPOT")
	url := fmt.Sprintf(okxApiBaseUrl, "market/candles?instId=BTC-USDT&limit=1")
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}
	dataBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}
	data := gjson.ParseBytes(dataBytes)
	//fmt.Println("data:", data)
	currencyList := data.Get("data").Array()

	for _, c := range currencyList {
		s := c.String()
		fmt.Println("s:", s)
	}
	//i := 0
	//
	//for _, currency := range currencyList {
	//	if currency.Get("quoteCcy").String() == "USDT" {
	//		//result = append(result, currency.Get("baseCcy").String())
	//		//fmt.Println("current:", currency.Get("baseCcy").String())
	//		slog.Info(fmt.Sprintf("currency:%s", currency.Get("baseCcy").String()))
	//		i++
	//	}
	//}
	//fmt.Println("total currency:", i)
}

func WsGetMarketData() {
	// OKX WebSocket API endpoint
	wsURL := "wss://ws.okx.com:8443/ws/v5/business"

	// Create a new WebSocket connection to the OKX server
	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// Subscribe to the BTC/USDT ticker channel
	subscribe := map[string]interface{}{
		"op":   "subscribe",
		"args": []map[string]string{{"channel": "candle1s", "instId": "BTC-USDT"}, {"channel": "candle1s", "instId": "ETH-USDT"}},
	}
	if err := c.WriteJSON(subscribe); err != nil {
		log.Fatal("subscribe:", err)
	}

	done := make(chan struct{})

	// Start reading messages from the WebSocket connection
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			r := gjson.ParseBytes(message)
			array := r.Get("data").Array()
			if len(array) == 0 {
				continue
			}
			//fmt.Printf("Received: %s\n", message)
			//fmt.Printf("data: %s\n", array[0])
			fmt.Printf("O: %s\n", array[0].Array()[0])
		}
	}()
	time.Sleep(20 * time.Second)
}

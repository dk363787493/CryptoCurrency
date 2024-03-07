package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ExchangeInfo struct {
	Symbols []Symbol `json:"symbols"`
}

type Symbol struct {
	Symbol  string   `json:"symbol"`
	Filters []Filter `json:"filters"`
}

type Filter struct {
	FilterType string `json:"filterType"`
	MinQty     string `json:"minQty,omitempty"`
	MaxQty     string `json:"maxQty,omitempty"`
	StepSize   string `json:"stepSize,omitempty"`
}

func ExchangeInfoCall() {
	//baseURL := "https://api.binance.com"

	//client := binance_connector.NewClient("", "", baseURL)
	//
	//// ExchangeInfo
	//exchangeInfo, err := client.NewExchangeInfoService().Do(context.Background())
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(binance_connector.PrettyPrint(exchangeInfo))
	resp, err := http.Get("https://api.binance.com/api/v3/exchangeInfo")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var exchangeInfo ExchangeInfo
	err = json.Unmarshal(body, &exchangeInfo)
	if err != nil {
		panic(err)
	}

	for _, symbol := range exchangeInfo.Symbols {
		for _, filter := range symbol.Filters {
			if filter.FilterType == "LOT_SIZE" {
				fmt.Printf("Symbol: %s, MinQty: %s, MaxQty: %s, StepSize: %s\n",
					symbol.Symbol, filter.MinQty, filter.MaxQty, filter.StepSize)
			}
		}
	}

}

func main() {
	ExchangeInfoCall()
	//TestUSDT(2000)
}

func TestUSDT(buyAmt float64) {
	step := 0.000010000
	btcPrice := 65327.66
	quity := buyAmt / btcPrice
	a := float64(int(quity/step)) * step
	fmt.Printf("rawquity:%f\n", quity)
	fmt.Printf("a:%f\n", a)
}

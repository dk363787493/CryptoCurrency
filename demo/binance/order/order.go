package main

import (
	"context"
	"fmt"
	binance_connector "github.com/binance/binance-connector-go"
)

var apiKey = "UHX1iqrKxSS40o7KZkzRAOc3ABpfFG3fzidEY5fTVkTdTtuQ6FLk0FTCLUFIQNEA"

var secretKey = "PjmFJPSz19Rb5EooVkGE3RZolfJWCSfRCw0zwRpJy1vgD0HgCX3HIbTVeoU3tPML"
var baseURL = "https://api.binance.com"

func NewOrder() {

	client := binance_connector.NewClient(apiKey, secretKey, baseURL)

	// Binance New Order endpoint - POST /api/v3/order
	//type:LIMIT
	//newOrder, err := client.NewCreateOrderService().Symbol("BTCUSDT").
	//	Side("BUY").Type("MARKET").QuoteOrderQty(10) .
	//	Do(context.Background())
	newOrder, err := client.NewCreateOrderService().Symbol("BTCUSDT").
		Side("SELL").Type("LIMIT").
		TimeInForce("GTC").
		Price(61500).
		Quantity(0.00000100).
		Do(context.Background())
	//Quantity(0.001).
	//newOrder, err := client.NewCreateOrderService().Symbol("BTCAAUSDT").
	//	Side("BUY").Type("MARKET").QuoteOrderQty(10).
	//	Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(binance_connector.PrettyPrint(newOrder))
}
func QueryOrder(orderId int64) {

	client := binance_connector.NewClient(apiKey, secretKey, baseURL)

	// Binance Query Order (USER_DATA) - GET /api/v3/order
	queryOrder, err := client.NewGetOrderService().Symbol("BTCUSDT").
		OrderId(orderId).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(binance_connector.PrettyPrint(queryOrder))
}

func GetAllOrders() {
	client := binance_connector.NewClient(apiKey, secretKey, baseURL)

	// Binance Get all account orders; active, canceled, or filled - GET /api/v3/allOrders
	getAllOrders, err := client.NewGetAllOrdersService().Symbol("BTCUSDT").
		StartTime(1709263200000).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(binance_connector.PrettyPrint(getAllOrders))
}

func main() {
	NewOrder()
	//GetAllOrders()
	//25234040073,25235269232,
	//QueryOrder(25235269232)
}

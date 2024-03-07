package trade

import (
	"context"
	"fmt"
	binance_connector "github.com/binance/binance-connector-go"
	"github.com/tidwall/gjson"
	"strconv"
)

var client *binance_connector.Client
var apiKey = "UHX1iqrKxSS40o7KZkzRAOc3ABpfFG3fzidEY5fTVkTdTtuQ6FLk0FTCLUFIQNEA"
var secretKey = "PjmFJPSz19Rb5EooVkGE3RZolfJWCSfRCw0zwRpJy1vgD0HgCX3HIbTVeoU3tPML"
var baseURL = "https://api.binance.com"

func init() {
	client = binance_connector.NewClient(apiKey, secretKey, baseURL)
}

type Order struct {
	Symbol        string
	OrderType     string
	Price         float64 //
	Status        string
	Quantity      float64 // 交易数量
	QuoteOrderQty float64 // 订单成交金额，单位USDT
}

func Calculate(order *Order, stepSize float64) float64 {
	//
	executeQuantity := float64(int(order.Quantity/stepSize)) * stepSize
	return executeQuantity
}

func Buy(order *Order) (*Order, error) {
	newOrder, err := client.NewCreateOrderService().Symbol(order.Symbol).
		Side("BUY").
		//TimeInForce("GTC").
		Type("MARKET").
		QuoteOrderQty(order.QuoteOrderQty).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	ordrStr := binance_connector.PrettyPrint(newOrder)
	r := gjson.Parse(ordrStr)
	if r.Get("status").String() != "FILLED" {
		return nil, fmt.Errorf("faied to make a order,currency:%s", order.Symbol)
	}
	price, _ := strconv.ParseFloat(r.Get("price").String(), 64)
	executedQty, _ := strconv.ParseFloat(r.Get("executedQty").String(), 64)
	return &Order{
		Price:     price,
		Quantity:  executedQty,
		Symbol:    order.Symbol,
		OrderType: "BUY",
	}, nil
}
func MockBuy(order *Order) (*Order, error) {
	order.Price = 10.00
	return order, nil
}

func MockSell(order *Order) (*Order, error) {

	return order, nil
}

func Sell(order *Order, stepSize float64) (*Order, error) {

	qty := Calculate(order, stepSize)

	newOrder, err := client.NewCreateOrderService().Symbol(order.Symbol).
		Side("SELL").
		TimeInForce("GTC").
		Type("LIMIT").
		Price(order.Price).
		Quantity(qty).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	ordrStr := binance_connector.PrettyPrint(newOrder)
	r := gjson.Parse(ordrStr)
	if r.Get("status").String() != "FILLED" {
		return nil, fmt.Errorf("faied to make a order,currency:%s", order.Symbol)
	}

	price, _ := strconv.ParseFloat(r.Get("price").String(), 64)
	executedQty, _ := strconv.ParseFloat(r.Get("executedQty").String(), 64)
	return &Order{
		Price:     price,
		Quantity:  executedQty,
		Symbol:    order.Symbol,
		OrderType: "SELL",
	}, nil
}

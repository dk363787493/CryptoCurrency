package exchange

import (
	"CryptoCurrency/config"
	"CryptoCurrency/tradev2/common"
	"context"
	"fmt"
	binance_connector "github.com/binance/binance-connector-go"
	"github.com/tidwall/gjson"
	"strconv"
	"time"
)

type Exchange interface {
	GetName() string
	GetDefaultKlineStream(symbol string) (<-chan string, error)
	GetKlineStream(symbol string, interval string, wsKlineHandlerGenerator func() (binance_connector.WsKlineHandler, chan string), errHandlerGenerator func() binance_connector.ErrHandler) (<-chan string, error)
	MakeSpotTrading(order *Order) (*Order, error)
}

type Order struct {
	Side          string // 0 buy 1 sell
	OrderType     string // LIMIT MARKET
	Symbol        string
	Status        string
	Price         float64
	Amount        float64
	Quantity      float64
	QuoteOrderQty float64
	CreateAt      time.Time
	TimeStamp     int64
}

type BinanceExchange struct {
	Client       *binance_connector.Client
	StreamClient *binance_connector.WebsocketStreamClient
}

func NewBinanceExchange() *BinanceExchange {
	client := binance_connector.NewClient(config.Configuration.BinanceApiKey, config.Configuration.BinanceSecretKey, config.Configuration.BinanceBaseURL)
	websocketStreamClient := binance_connector.NewWebsocketStreamClient(false)
	exchange := &BinanceExchange{
		Client:       client,
		StreamClient: websocketStreamClient,
	}
	return exchange
}

func (binance *BinanceExchange) GetName() string {
	return "binance"
}

func (binance *BinanceExchange) GetDefaultKlineStream(symbol string) (<-chan string, error) {

	wsKlineHandlerGenerator := func() (binance_connector.WsKlineHandler, chan string) {
		dataCh := make(chan string, 1024)
		return func(event *binance_connector.WsKlineEvent) {
			dataCh <- binance_connector.PrettyPrint(event)
		}, dataCh
	}

	errHandlerGenerator := func() binance_connector.ErrHandler {
		return func(err error) {
			fmt.Println(err)
		}
	}
	return binance.GetKlineStream(symbol, "1s", wsKlineHandlerGenerator, errHandlerGenerator)
}

func (binance *BinanceExchange) GetKlineStream(symbol string, interval string,
	wsKlineHandlerGenerator func() (binance_connector.WsKlineHandler, chan string), errHandlerGenerator func() binance_connector.ErrHandler) (<-chan string, error) {
	wsHandler, dataCh := wsKlineHandlerGenerator()
	errHandler := errHandlerGenerator()
	_, _, err := binance.StreamClient.WsKlineServe(symbol, interval, wsHandler, errHandler)
	if err != nil {
		return nil, err
	}
	return dataCh, nil
}

func (binance *BinanceExchange) MakeSpotTrading(order *Order) (*Order, error) {
	switch order.Side {
	case common.SideSell:
		{
			newOrder, err := binance.Client.NewCreateOrderService().Symbol(order.Symbol).
				Side(order.Side).Type(order.OrderType).
				TimeInForce("GTC").
				Price(order.Price).
				Quantity(order.Quantity).
				Do(context.Background())
			if err != nil {
				return nil, err
			}
			r := gjson.Parse(binance_connector.PrettyPrint(newOrder))
			o := new(Order)
			o.Status = r.Get("status").String()
			o.OrderType = r.Get("type").String()
			o.Side = common.SideSell
			o.Price, _ = strconv.ParseFloat(r.Get("price").String(), 64)
			o.Quantity, _ = strconv.ParseFloat(r.Get("executedQty").String(), 64)
			return o, nil
		}
	case common.SideBuy:
		{
			newOrder, err := binance.Client.NewCreateOrderService().Symbol(order.Symbol).
				Side(order.Side).Type(order.OrderType).
				QuoteOrderQty(order.QuoteOrderQty).
				Do(context.Background())
			if err != nil {
				return nil, err
			}
			r := gjson.Parse(binance_connector.PrettyPrint(newOrder))
			o := new(Order)
			o.Price, _ = strconv.ParseFloat(r.Get("price").String(), 64)
			o.Quantity, _ = strconv.ParseFloat(r.Get("executedQty").String(), 64)
			o.Status = r.Get("status").String()
			o.OrderType = r.Get("type").String()
			o.Side = common.SideBuy
			return o, nil
		}
	default:
		return nil, fmt.Errorf("can not support order side:%s", order.Side)
	}

}

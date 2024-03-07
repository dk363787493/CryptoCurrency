package currency

import (
	"errors"
	"log/slog"
)

type ExchangeId uint8

const (
	Binance ExchangeId = 1
	Okx     ExchangeId = 2
)

type Exchange interface {
	GetExchangeName() string
	GetCurrencyList() []string
	GetCurrencyInfoFromExchange()
	GetMarketDataStreamFromExchange(currency []string) error
	UpdateCache(symbol string, price float64) error
	AsyncUpdateSymbolPrice()
	GetPrice(symbol string) (float64, error)
}

type KlineData struct {
	Exchange ExchangeId
	Symbol   string
	Price    float64
}

type BaseExchange struct {
	Exchange
}

func (*BaseExchange) GetExchangeName() string {
	slog.Error("unimplementation function")
	return ""
}
func (*BaseExchange) GetCurrencyList() []string {
	slog.Error("unimplementation function")
	return nil
}
func (*BaseExchange) GetCurrencyInfoFromExchange() {
	slog.Error("unimplementation function")
}
func (*BaseExchange) GetMarketDataStreamFromExchange(currency []string) error {
	slog.Error("unimplementation function")
	return errors.New("unimplementation function")
}

func (*BaseExchange) UpdateCache(symbol string, price float64) error {
	slog.Error("unimplementation function")
	return errors.New("unimplementation function")
}

func (*BaseExchange) AsyncUpdateSymbolPrice() {
	slog.Error("unimplementation function")
}

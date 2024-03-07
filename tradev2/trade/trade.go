package trade

import (
	"CryptoCurrency/tradev2/exchange"
	"CryptoCurrency/tradev2/strategy"
)

type Engine struct {
	exchange        exchange.Exchange
	defaultStrategy strategy.Strategy
	symbols         string
}

func NewEngine() *Engine {

	return &Engine{}
}

func (engin *Engine) Start() {
	engin.defaultStrategy.ExecuteStrategy(engin.exchange)
}

package strategy

import "CryptoCurrency/tradev2/exchange"

type FixationrateStrategyExecutor struct {
}

func (f *FixationrateStrategyExecutor) Execute(exchang exchange.Exchange, order *exchange.Order) *exchange.Order {
	return nil
}

package strategy

import "CryptoCurrency/tradev2/exchange"

type DefaultStrategy struct {
	Side                string
	Status              int8 // 1:BUY  2:SELL
	buyStategyExecutor  BuyStategyExecutor
	sellStategyExecutor SellStategyExecutor
}

func (d *DefaultStrategy) ExecuteStrategy(exchange exchange.Exchange) {
	for {
		buyOrder := d.buyStategyExecutor.Execute(exchange)
		d.Status = 1
		sellOrder := d.sellStategyExecutor.Execute(exchange, buyOrder)
		d.Status = 2
		_ = sellOrder
	}
}

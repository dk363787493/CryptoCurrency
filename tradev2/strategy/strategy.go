package strategy

import "CryptoCurrency/tradev2/exchange"

type Strategy interface {
	ExecuteStrategy(exchange.Exchange)
}

type BuyStategyExecutor interface {
	Execute(exchange exchange.Exchange) *exchange.Order
}

type SellStategyExecutor interface {
	Execute(exchange exchange.Exchange, order *exchange.Order) *exchange.Order
}

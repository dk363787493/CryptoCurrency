package strategy

type Strategy interface {
	ExecuteStrategy(param ...interface{})
	ExecuteBuyStrategy(param ...interface{}) *Order
	ExecuteSellStrategy(param ...interface{})
}

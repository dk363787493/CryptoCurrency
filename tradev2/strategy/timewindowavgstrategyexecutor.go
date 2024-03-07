package strategy

import (
	"CryptoCurrency/tradev2/common"
	"CryptoCurrency/tradev2/exchange"
	"sync"
	"time"
)

type TimeWindowAvgStrategyExecutor struct {
	TradeDataCh chan exchange.Order
	Side        string
	Status      int8
	Symbol      string
	ShortTerm   int
	LongTerm    int
	once        sync.Once
	exchange    exchange.Exchange
	closeSign   chan struct{}
	buffer      sync.Pool
}

func (t *TimeWindowAvgStrategyExecutor) init() {
	go func() {
		dataCh, err := t.exchange.GetDefaultKlineStream(t.Symbol)
		if err != nil {
			t.Stop()
		}
		timer := time.NewTimer(5 * time.Second)
		for {
			select {
			case <-dataCh:
			case <-timer.C:
				{
					d := <-dataCh
					t.TradeDataCh <- common.ParsePirceData(d)
				}
			}
		}

	}()
}

func (t *TimeWindowAvgStrategyExecutor) Execute(exchange exchange.Exchange) *exchange.Order {
	t.once.Do(t.init)

	priceBuffer := t.buffer.Get().([]float64)
	defer t.buffer.Put(priceBuffer)
	for i := 0; i < len(priceBuffer); i++ {
		order := <-t.TradeDataCh
		priceBuffer[i] = order.Price
		if i >= t.LongTerm-1 {
			result := t.caculate(priceBuffer, i)
			if result != nil {
				return result
			}
		}
	}
	return nil
}

func (t *TimeWindowAvgStrategyExecutor) caculate(datas []float64, size int) *exchange.Order {

	return nil
}

func (t *TimeWindowAvgStrategyExecutor) Stop() {
	close(t.closeSign)
}

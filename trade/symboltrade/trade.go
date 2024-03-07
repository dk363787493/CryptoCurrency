package symboltrade

import (
	"CryptoCurrency/trade/strategy"
	"sync"
	"time"
)

var Cache sync.Map

type Trade struct {
	Symbol    string
	Iterval   time.Duration
	CirlDatas strategy.CirlDatas
	strategy  []strategy.Strategy
	lock      sync.Mutex
	SellCh    chan *strategy.Order
}

func NewTrade(symbol string, iterval time.Duration) *Trade {
	ch := make(chan *strategy.Order, 1024)
	t := &Trade{
		Symbol:  symbol,
		Iterval: iterval,
		CirlDatas: strategy.CirlDatas{
			Datas: make([]strategy.PriceData, 200),
		},
		strategy: make([]strategy.Strategy, 0),
		SellCh:   ch,
	}
	go func() {
		timer := time.NewTimer(5 * time.Second)
		for {
			select {
			case o := <-ch:
				{
					s, ok := Cache.Load(symbol)
				}
			}
		}
	}()

	return t
}

func (t *Trade) AddStrategy(strategy strategy.Strategy) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.strategy = append(t.strategy, strategy)
}

func (t *Trade) AppendData(data *strategy.PriceData) {
	t.CirlDatas.Lock.Lock()
	defer t.CirlDatas.Lock.Unlock()
	if (t.CirlDatas.EndIndex+1)%len(t.CirlDatas.Datas) == t.CirlDatas.StartIndex%len(t.CirlDatas.Datas) {
		t.CirlDatas.StartIndex++
	}
	t.CirlDatas.EndIndex++
	if t.CirlDatas.EndIndex-t.CirlDatas.StartIndex+1 > len(t.CirlDatas.Datas) {
		t.CirlDatas.StartIndex = t.CirlDatas.EndIndex + 1 - len(t.CirlDatas.Datas)
	}
	t.CirlDatas.Datas[t.CirlDatas.EndIndex] = *data
}

func (t *Trade) ExecuteStrategy() {
	for _, s := range t.strategy {
		s.ExecuteStrategy(&t.CirlDatas)
	}
}

package strategy

import (
	"CryptoCurrency/trade/symboltrade"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type CirlDatas struct {
	Datas      []PriceData
	Lock       sync.Mutex
	Status     uint8 //0 买  1 卖
	StartIndex int
	EndIndex   int
}
type PriceData struct {
	Timestamp int64
	Price     float64
	Symbol    string
}

// 计算移动平均值
// 初步策略，可以获取历史数据，根据近n天的趋势进行买卖，例如BTC 近5天是都是处于下降趋势，可以选择这个时机买入
func movingAverage(prices []float64) float64 {
	var sum float64
	for _, price := range prices {
		sum += price
	}
	return sum / float64(len(prices))
}

type Order struct {
	Symbol string
	Price  float64
	Side   uint8
}

type MoveAvgStrategy struct {
	ShortTermWindow int
	LongTermWindow  int
	SellCh          chan *Order
	Status          atomic.Int32 //0 买  1 卖
}

func (t *MoveAvgStrategy) ExecuteBuyStrategy(param ...interface{}) *Order {
	datas := param[0].(*CirlDatas)
	datas.Lock.Lock()
	datas.Lock.Unlock()
	// 用于储存移动平均值
	var shortTermMA, longTermMA float64

	for i := datas.StartIndex; i <= datas.EndIndex; i++ {
		// 当数据足够计算长期移动平均值时
		if i-datas.StartIndex >= t.LongTermWindow-1 {
			shortPrices := make([]float64, t.ShortTermWindow)
			longPrices := make([]float64, t.LongTermWindow)

			// 获取短期和长期窗口价格
			for j := 0; j < t.ShortTermWindow; j++ {
				subId := (i - j) % len(datas.Datas)
				shortPrices[j] = datas.Datas[subId].Price
			}
			for j := 0; j < t.LongTermWindow; j++ {
				subId := (i - j) % len(datas.Datas)
				longPrices[j] = datas.Datas[subId].Price
			}

			// 计算移动平均值
			newShortTermMA := movingAverage(shortPrices)
			newLongTermMA := movingAverage(longPrices)

			// 产生信号
			if shortTermMA < longTermMA && newShortTermMA > newLongTermMA {
				fmt.Printf("买入信号 @ %v\n", datas.Datas[i].Timestamp)
				// 市价成交，按比成交价格高固定几个百分点设置限价卖单
				price := 2.0
				t.Status.Store(1)
				symbol := datas.Datas[i%len(datas.Datas)].Symbol
				o := &Order{
					Symbol: symbol,
					Price:  price,
					Side:   0,
				}
				// 修改状态
				t.Status.Store(1)
				symboltrade.Cache.Store("symbol", o)
				return o
			}

			// 更新移动平均值
			shortTermMA = newShortTermMA
			longTermMA = newLongTermMA
		}
	}
	return nil
}

//func (t *MoveAvgStrategy) ExecuteBuyStrategy(param ...interface{}) *Order {
//	datas := param[0].(*CirlDatas)
//	datas.Lock.Lock()
//	datas.Lock.Unlock()
//	fmt.Printf("买入信号 \n")
//	// 市价成交，按比成交价格高固定几个百分点设置限价卖单
//	price := 2.0
//	t.Status.Store(1)
//	o := &Order{
//		Symbol: "BTCUSDT",
//		Price:  price,
//		Side:   0,
//	}
//	// 修改状态
//	t.Status.Store(1)
//	return o
//}

func (t *MoveAvgStrategy) ExecuteSellStrategy(param ...interface{}) {
	go func() {
		// 执行限价卖出策略
		time.Sleep(15 * time.Second)
		o := param[0].(*Order)
		// 修改状态
		fmt.Println("sell order:", o.Price)
		t.SellCh <- o
		//t.Status.Store(0)
		// 执行限价卖出策略
	}()
}

// 执行策略
func (t *MoveAvgStrategy) ExecuteStrategy(param ...interface{}) {
	if t.Status.Load() == 1 {
		return
	}
	order := t.ExecuteBuyStrategy(param...)

	if t.Status.Load() == 0 || order == nil {
		return
	}
	t.ExecuteSellStrategy(order)

}

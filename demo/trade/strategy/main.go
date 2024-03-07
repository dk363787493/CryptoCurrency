package main

import (
	"fmt"
)

// 假设这是我们收到的价格数据
type PriceData struct {
	Timestamp int64
	Price     float64
}

// 计算移动平均值
func movingAverage(prices []float64) float64 {
	var sum float64
	for _, price := range prices {
		sum += price
	}
	return sum / float64(len(prices))
}

// 执行策略
func executeStrategy(data []PriceData, shortTermWindow, longTermWindow int) {

	// 用于储存移动平均值
	var shortTermMA, longTermMA float64

	for i := 0; i < len(data); i++ {
		// 当数据足够计算长期移动平均值时
		if i >= longTermWindow-1 {
			shortPrices := make([]float64, shortTermWindow)
			longPrices := make([]float64, longTermWindow)

			// 获取短期和长期窗口价格
			for j := 0; j < shortTermWindow; j++ {
				shortPrices[j] = data[i-j].Price
			}
			for j := 0; j < longTermWindow; j++ {
				longPrices[j] = data[i-j].Price
			}

			// 计算移动平均值
			newShortTermMA := movingAverage(shortPrices)
			newLongTermMA := movingAverage(longPrices)

			// 产生信号
			if shortTermMA < longTermMA && newShortTermMA > newLongTermMA {
				fmt.Printf("买入信号 @ %v\n", data[i].Timestamp)
			} else if shortTermMA > longTermMA && newShortTermMA < newLongTermMA {
				fmt.Printf("卖出信号 @ %v\n", data[i].Timestamp)
			}

			// 更新移动平均值
			shortTermMA = newShortTermMA
			longTermMA = newLongTermMA
		}
	}
}

func main() {
	// 假设的价格数据
	priceData := []PriceData{
		// ...这里应该是一系列价格数据
	}

	// 执行策略
	executeStrategy(priceData, 5, 20) // 这里的数字5和20表示短期和长期窗口的大小
}

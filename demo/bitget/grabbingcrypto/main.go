package main

import (
	"CryptoCurrency/demo/bitget/grabbingcrypto/task"
	"sync"
)

func main() {
	s1 := &task.GrabbingTask{
		ProfitRate:    0.05,
		ApiKey:        "bg_8bb5b03eb0f08065b0442560441ed912",
		SecretKey:     "a43174b22a98f0195f6d4ac887707210b3168cca92c44ae4ac6449ed4566d505",
		Passphrase:    "Bsdk19901214123",
		Symbol:        "BTCUSDT",
		Amout:         "6",
		PriceScale:    0.01,
		QuantityScale: 0.000001,
		ActionTime:    "2024-03-14 22:59:00",
		SubTime:       200,
	}
	//s2 := &task.GrabbingTask{
	//	ProfitRate:    0.05,
	//	ApiKey:        "bg_8bb5b03eb0f08065b0442560441ed912",
	//	SecretKey:     "a43174b22a98f0195f6d4ac887707210b3168cca92c44ae4ac6449ed4566d505",
	//	Passphrase:    "Bsdk19901214123",
	//	Symbol:        "ETHUSDT",
	//	Amout:         "6",
	//	PriceScale:    0.01,
	//	QuantityScale: 0.0001,
	//	ActionTime:    "2024-03-14 23:16:00",
	//	SubTime:       200,
	//}

	tasks := make([]*task.GrabbingTask, 0)
	tasks = append(tasks, s1)
	//tasks = append(tasks, s2)
	var wg sync.WaitGroup
	for _, t := range tasks {
		t.Init()
		wg.Add(1)
		go func(t *task.GrabbingTask) {
			defer wg.Done()
			t.Ready()
			var err error
			// buy order
			for i := 0; i < 6; i++ {
				err = t.MakeBuyOrder()
				if err == nil {
					break
				}
			}
			if err != nil {
				return
			}

			// get order
			for i := 0; i < 5; i++ {
				err = t.GetOrder()
				if err == nil {
					break
				}
			}
			if err != nil {
				return
			}

			//sell order
			err = t.MakeSellOrder()
			if err != nil {
				return
			}
		}(t)
	}
	wg.Wait()

}

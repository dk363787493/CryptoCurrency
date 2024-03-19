package main

import (
	"CryptoCurrency/demo/bitget/grabbingcrypto/task"
	"sync"
)

func main() {
	s1 := &task.GrabbingTask{
		ProfitRate:    0.03,
		ApiKey:        "bg_8bb5b03eb0f08065b0442560441ed912",
		SecretKey:     "a43174b22a98f0195f6d4ac887707210b3168cca92c44ae4ac6449ed4566d505",
		Passphrase:    "Bsdk19901214123",
		Symbol:        "ZKUSDT",
		Amout:         "10",
		PriceScale:    0.001,
		QuantityScale: 0.01,
		ActionTime:    "2024-03-19 18:00:00",
		//180(L)
		SubTime: 220,
	}

	s2 := &task.GrabbingTask{
		ProfitRate:    0.03,
		ApiKey:        "bg_8bb5b03eb0f08065b0442560441ed912",
		SecretKey:     "a43174b22a98f0195f6d4ac887707210b3168cca92c44ae4ac6449ed4566d505",
		Passphrase:    "Bsdk19901214123",
		Symbol:        "ZKUSDT",
		Amout:         "10",
		PriceScale:    0.001,
		QuantityScale: 0.01,
		ActionTime:    "2024-03-19 18:00:00",
		//200
		SubTime: 190,
	}
	//s3 := &task.GrabbingTask{
	//	ProfitRate:    0.06,
	//	ApiKey:        "bg_8bb5b03eb0f08065b0442560441ed912",
	//	SecretKey:     "a43174b22a98f0195f6d4ac887707210b3168cca92c44ae4ac6449ed4566d505",
	//	Passphrase:    "Bsdk19901214123",
	//	Symbol:        "SLERFUSDT",
	//	Amout:         "10",
	//	PriceScale:    0.0001,
	//	QuantityScale: 0.01,
	//	ActionTime:    "2024-03-18 18:00:00",
	//	//250
	//	SubTime: 100,
	//}
	//s4 := &task.GrabbingTask{
	//	ProfitRate:    0.06,
	//	ApiKey:        "bg_8bb5b03eb0f08065b0442560441ed912",
	//	SecretKey:     "a43174b22a98f0195f6d4ac887707210b3168cca92c44ae4ac6449ed4566d505",
	//	Passphrase:    "Bsdk19901214123",
	//	Symbol:        "SLERFUSDT",
	//	Amout:         "10",
	//	PriceScale:    0.0001,
	//	QuantityScale: 0.01,
	//	ActionTime:    "2024-03-18 18:00:00",
	//	//250
	//	SubTime: 50,
	//}
	//s5 := &task.GrabbingTask{
	//	ProfitRate:    0.06,
	//	ApiKey:        "bg_8bb5b03eb0f08065b0442560441ed912",
	//	SecretKey:     "a43174b22a98f0195f6d4ac887707210b3168cca92c44ae4ac6449ed4566d505",
	//	Passphrase:    "Bsdk19901214123",
	//	Symbol:        "SLERFUSDT",
	//	Amout:         "10",
	//	PriceScale:    0.0001,
	//	QuantityScale: 0.01,
	//	ActionTime:    "2024-03-18 18:00:00",
	//	//250
	//	SubTime: 20,
	//}

	tasks := make([]*task.GrabbingTask, 0)
	tasks = append(tasks, s1)
	tasks = append(tasks, s2)
	//tasks = append(tasks, s3)
	//tasks = append(tasks, s4)
	//tasks = append(tasks, s5)
	var wg sync.WaitGroup
	for _, t := range tasks {
		t.Init()
		wg.Add(1)
		go func(t *task.GrabbingTask) {
			defer wg.Done()
			t.Ready()
			var err error
			// buy order
			for i := 0; i < 5; i++ {
				err = t.MakeBuyOrder()
				if err == nil {
					break
				}
			}
			if err != nil {
				return
			}

			// get order
			for i := 0; i < 10; i++ {
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

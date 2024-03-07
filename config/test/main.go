package main

import (
	"CryptoCurrency/model"
	"fmt"
	"sync"
	"time"
)

var once sync.Once

func main() {
	//trade := config.Configuration.Trade
	//for k, v := range trade {
	//	fmt.Printf("k:%s,v:%s \n", k, v)
	//}
	//keys := maps.Keys(trade)
	//
	//slices.Sort(keys)
	//for _, k := range keys {
	//	fmt.Println(k)
	//}

	//binance := config.Configuration.Binance
	//for k, m := range binance {
	//	fmt.Printf("k:%s,v:", k)
	//	fmt.Println(m.CompletedOrderNumOfLatest30Day)
	//}

	//rule := config.Configuration.RiskControlRule
	//fmt.Println(rule.Fraud.StepOne.Binance["caseone"])

	//{map[caseone:{1 10 2}] map[casefive:{0 0} casefour:{0 5} caseone:{2 20} casethree:{1 10} casetwo:{1 10}]}
	//var rule Rule = func(p ...any) {
	//	fmt.Println("......", p[0])
	//
	//}
	//rule("hello")
	//
	//var i uint8 = 2
	//var ii interface{} = i
	//u := ii.(uint8)
	//fmt.Println(u)

	//onceBody := func() {
	//	fmt.Println("Only once")
	//}
	//foo := sync.OnceFunc(onceBody)
	//
	//for i := 0; i < 10; i++ {
	//	foo()
	//}
	//pp := Person{}
	//namePtr := A()
	//fmt.Println(*namePtr)
	//runtime.GC()
	//p := uintptr(unsafe.Pointer(namePtr)) + unsafe.Sizeof(pp.name)
	//fmt.Println(*namePtr)
	//
	//fmt.Println(*(*int8)(unsafe.Pointer(p)))
	//m1 := map[string]int{
	//	"Lee": 1,
	//}
	//m2 := m1
	//m2["Lee"] = 2
	//fmt.Println(m1["Lee"])
	//totalAmtStr := "1.00a"
	//totalAmt := 4.0
	//totalAmt, _ = strconv.ParseFloat(totalAmtStr, 64)
	//fmt.Println(totalAmt)
	//
	//aa := 1000
	//bb := float64(aa)
	//fmt.Println(bb)
	//m := map[string]interface{}{"Lee": 18}
	//i, ok := m["Lee"].(string)
	//if !ok {
	//	fmt.Println(i)
	//} else {
	//
	//	fmt.Println(i)
	//}

	//err1 := errors.New("err1...")
	//var err1 error
	//err2 := errors.New("err2...")
	//err := errors.Join(err1, err2)
	//fmt.Printf("%v", err)

	//var a uint8 = 10
	//b := (*int)(unsafe.Pointer(&a))
	//fmt.Println(*b)
	//TestDefter()

	//m := map[string]int{"Lee": 100}
	//v := m["AAA"]
	//fmt.Println("v:", v)
	//s := new(Stu)
	//s.name = "Lee"
	//once.Do(s.PrintName）
	now := time.Now()
	//add := now.Add(5 * time.Second)
	////add = add.Add(10 * time.Second)
	//sub := add.Sub(now)
	//s := time.Duration(sub.Seconds())
	//_ = s
	//fmt.Println(int64(s))
	//time.Sleep(s * time.Second)
	fmt.Println(now)

	c := model.GetAllOpenTradingCurrency("CLOSE")[0]
	fmt.Println(c.StatustartmissionDate)
	sub := c.StatustartmissionDate.Sub(now)
	fmt.Println(sub.Seconds())
}

type Stu struct {
	name string
}

func (s *Stu) PrintName() {
	fmt.Println("name:", s.name)
}
func TestDefter() {
	for _, s := range []int{1, 2, 3, 4} {
		defer func(i *int) { fmt.Println(*i) }(&s)
	}
}

type Rule func(p ...any)

type Person struct {
	name int8
	age  int8
}

func A() *int8 {
	p := Person{name: 28, age: 30}
	return &p.name
}

type Animal struct {
	name string
	age  int
	c    *Cat
}

type Cat struct {
	name string
	age  int
	a    *Animal
}

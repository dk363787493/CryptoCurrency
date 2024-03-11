package main

import (
	"bitget/pkg/client"
	v1 "bitget/pkg/client/v1"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
	"time"
)

type Order struct {
	Price    string
	Quantity string
}

var spotCli *v1.SpotOrderClient
var cli *client.BitgetApiClient

var symbol = "ORBKUSDT_SPBL"
var profitRate = 0.1

const layout = "2006-01-02 15:04:05"

// success orderId
var orderId = ""

func main() {

	TimeSleep()

	for i := 0; i < 10; i++ {
		err, _ := Start()
		if err == nil {
			break
		}
	}

}

func TimeSleep() {
	timeStr := "2024-03-11 17:30:00"

	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}

	localTime, err := time.ParseInLocation(layout, timeStr, location)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return
	}
	fmt.Println("Parsed time in UTC+8:", localTime)
	//100  300 millsec
	sleep := localTime.Sub(time.Now()).Milliseconds() - 300
	if sleep < 0 {
		return
	}
	fmt.Println("sleep(ms):", sleep)
	time.Sleep(time.Duration(sleep) * time.Millisecond)
}

func Start() (error, int) {
	//buy order
	var err error

	if orderId == "" {
		orderId, err = MakertBuyOrder()
	}

	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return err, 1
	}
	// get order info
	o, err := GetOrder(orderId)
	//o, err := MockGetOrder(order)
	for i := 0; i < 5; i++ {
		o, err = GetOrder(orderId)
		if err == nil {
			break
		}
	}
	if err != nil {
		return err, 2
	}
	fmt.Printf("sell order:%+v\n", o)
	// make a limit sell order

	err = LimitSellOrder(o)
	if err != nil {
		return err, 3
	}
	return nil, 4
}

func init() {
	spotCli = new(v1.SpotOrderClient).Init()
	spotCli.BitgetRestClient.ApiKey = "bg_8bb5b03eb0f08065b0442560441ed912"
	spotCli.BitgetRestClient.ApiSecretKey = "a43174b22a98f0195f6d4ac887707210b3168cca92c44ae4ac6449ed4566d505"
	spotCli.BitgetRestClient.Passphrase = "Bsdk19901214123"

	cli = new(client.BitgetApiClient).Init()
	cli.BitgetRestClient.ApiKey = "bg_8bb5b03eb0f08065b0442560441ed912"
	cli.BitgetRestClient.ApiSecretKey = "a43174b22a98f0195f6d4ac887707210b3168cca92c44ae4ac6449ed4566d505"
	cli.BitgetRestClient.Passphrase = "Bsdk19901214123"
}

func MakertBuyOrder() (string, error) {

	params := make(map[string]string)
	//params["symbol"] = "WENUSDT_SPBL"
	params["symbol"] = symbol
	params["side"] = "buy"
	params["orderType"] = "market"
	params["force"] = "normal"
	params["quantity"] = "6"

	resp, err := spotCli.PlaceOrder(params)
	if err != nil {
		println(err.Error())
		return "", err
	}
	if gjson.Parse(resp).Get("code").String() != "00000" {
		return "", fmt.Errorf("err when buy order,%s", gjson.Parse(resp).String())
	}
	order := gjson.Parse(resp).Get("data.orderId").String()
	fmt.Printf("success orderId:%s,resp:%s \n", order, resp)
	return order, err
}

func LimitSellOrder(o Order) error {
	priceStr := CaculteSellPrice(o.Price, profitRate+1)
	qantityeStr := CaculteSellQantity(o.Quantity)
	params := make(map[string]string)
	//params["symbol"] = "WENUSDT_SPBL"
	params["symbol"] = symbol
	params["side"] = "sell"
	params["orderType"] = "limit"
	params["force"] = "normal"
	params["price"] = priceStr
	params["quantity"] = qantityeStr

	resp, err := spotCli.PlaceOrder(params)
	if err != nil {
		println(err.Error())
		return err
	}
	fmt.Printf("sell order,priceStr:%s,qantityeStr:%s", priceStr, qantityeStr)
	fmt.Printf("success sell order:%s\n", resp)
	return nil
}

func GetOrder(orderId string) (Order, error) {
	if orderId == "" {
		return Order{}, errors.New("err: orderId can not be empty")
	}
	params := make(map[string]string)
	params["orderId"] = orderId

	resp, err := cli.Post("/api/spot/v1/trade/orderInfo", params)
	if err != nil {
		println(err.Error())
		return Order{}, err
	}
	r := gjson.Parse(resp)
	if r.Get("code").String() != "00000" {
		return Order{}, fmt.Errorf("err when query order,%s", r.String())
	}
	dataArr := r.Get("data").Array()
	if len(dataArr) <= 0 {
		return Order{}, fmt.Errorf("err: get order info,orderId:%s", orderId)
	}
	subR := r.Get("data").Array()[0]
	status := subR.Get("status").String()
	if status != "full_fill" {
		return Order{}, fmt.Errorf("err: get order info,orderId:%s,status:%s", orderId, status)
	}
	priceStr := subR.Get("fillPrice").String()
	quantityStr := subR.Get("fillQuantity").String()
	fmt.Printf("get order info:%s\n", r.String())

	return Order{
		Price:    priceStr,
		Quantity: quantityStr,
	}, nil
}

func ParseSale(f string) float64 {
	trim := strings.Trim(f, "0")
	split := strings.Split(trim, ".")
	if len(split) < 2 {
		return 0.00
	}
	s := split[1]
	var builder strings.Builder
	builder.WriteString("0.")
	for i, _ := range s {
		if i == len(s)-1 {
			builder.WriteString("1")
			break
		}
		builder.WriteString("0")
	}
	s2 := builder.String()
	ff, _ := strconv.ParseFloat(s2, 64)
	return ff
}

func CaculteSellPrice(price string, factor float64) string {
	scale := ParseSale(price)
	priceDecimal, _ := decimal.NewFromString(price)
	scaleDecimal := decimal.NewFromFloat(scale)
	pp := priceDecimal.Mul(decimal.NewFromFloat(factor))
	i := pp.Div(scaleDecimal).IntPart()
	dd := decimal.NewFromInt(i).Mul(scaleDecimal)
	return dd.String()
}

func CaculteSellQantity(qantity string) string {
	scale := ParseSale(qantity)
	qantityDecimal, _ := decimal.NewFromString(qantity)
	scaleDecimal := decimal.NewFromFloat(scale)
	// 扣除0.1%点
	pp := qantityDecimal.Mul(decimal.NewFromFloat(0.999))
	i := pp.Div(scaleDecimal).IntPart()
	dd := decimal.NewFromInt(i).Mul(scaleDecimal)
	return dd.String()
}

func MockGetOrder(orderId string) (Order, error) {

	return Order{
		Price:    "1.0002310010",
		Quantity: "0.000234",
	}, nil
}
func MockBuy() (string, error) {
	return "1111011", nil
}

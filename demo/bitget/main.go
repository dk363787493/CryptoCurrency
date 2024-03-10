package main

import (
	"bitget/pkg/client"
	v1 "bitget/pkg/client/v1"
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

var symbol = "BTCUSDT_SPBL"
var order = "1111111"

const layout = "2006-01-02 15:04:05"

func main() {

	//TimeSleep()

	for i := 0; i < 6; i++ {
		err, _ := Start()
		if err == nil {
			break
		}
	}

}

func TimeSleep() {
	timeStr := "2024-03-06 19:00:00"

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
	sleep := localTime.Sub(time.Now()).Seconds() - 0.01
	if sleep < 0 {
		return
	}
	time.Sleep(time.Duration(sleep) * time.Second)
}

func Start() (error, int) {
	//buy order
	var err error
	if order == "" {
		order, err = MakertBuyOrder()
	}
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return err, 1
	}
	// get order info
	o, err := GetOrder(order)
	for i := 0; i < 3; i++ {
		if err == nil {
			break
		}
		o, err = GetOrder(order)
	}
	if err != nil {
		fmt.Println("err:", err.Error())
		return err, 2
	}

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
	params["symbol"] = symbol
	params["side"] = "buy"
	params["orderType"] = "market"
	params["force"] = "normal"
	params["quantity"] = "6"

	resp, err := spotCli.PlaceOrder(params)
	if err != nil {
		fmt.Println("err:", err.Error())
		return "", err
	}
	fmt.Println("buy order:", resp)
	if gjson.Parse(resp).Get("code").String() != "00000" {
		return "", fmt.Errorf("err when buy order,code:%s", gjson.Parse(resp).Get("code").String())
	}
	order := gjson.Parse(resp).Get("data.orderId").String()
	return order, err
}

func LimitSellOrder(o Order) error {
	priceStr := CaculteSellPrice(o.Price, 1.5)
	params := make(map[string]string)
	//params["symbol"] = "WENUSDT_SPBL"
	params["symbol"] = symbol
	params["side"] = "sell"
	params["orderType"] = "limit"
	params["force"] = "normal"
	params["price"] = priceStr
	params["quantity"] = o.Quantity

	resp, err := spotCli.PlaceOrder(params)
	if err != nil {
		fmt.Println("err:", err.Error())
		return err
	}
	fmt.Printf("sell param:%+v\n", params)
	fmt.Println(resp)
	return nil
}

func GetOrder(orderId string) (Order, error) {

	params := make(map[string]string)
	params["orderId"] = orderId

	resp, err := cli.Post("/api/spot/v1/trade/orderInfo", params)
	if err != nil {
		fmt.Println(err.Error())
		return Order{}, err
	}
	r := gjson.Parse(resp)
	code := r.Get("code").String()

	if code != "00000" {
		fmt.Printf("order info:%s\n", r)
		err = fmt.Errorf("err when query order,code:%s", r.Get("code").String())
		return Order{}, err
	}
	datas := r.Get("data").Array()
	if len(datas) <= 0 {
		err = fmt.Errorf("no data when get data ,orderId:%s", orderId)
		return Order{}, err
	}
	subR := datas[0]
	status := subR.Get("status").String()
	if !(status == "filled" || status == "full_fill") {
		fmt.Printf("order info:%s\n", r)
		err = fmt.Errorf("err when query order,status:%s", status)
		return Order{}, err
	}

	priceStr := subR.Get("fillPrice").String()

	quantityStr := subR.Get("fillQuantity").String()
	fmt.Println("success getting order info:", r.String())
	return Order{
		Price:    priceStr,
		Quantity: quantityStr,
	}, nil
}

func ParseSale(f string) float64 {
	trim := strings.Trim(f, "0")
	split := strings.Split(trim, ".")
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

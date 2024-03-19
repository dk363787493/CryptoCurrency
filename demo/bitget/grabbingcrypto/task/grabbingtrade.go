package task

import (
	"CryptoCurrency/demo/bitget/utils"
	"bitget/pkg/client"
	v1 "bitget/pkg/client/v1"
	"bitget/uninternal/common"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"time"
)

const layout = "2006-01-02 15:04:05"

type Order struct {
	OrderId      string
	SellPrice    string
	SellQuantity string
}

type GrabbingTask struct {
	SpotCli       *v1.SpotOrderClient
	Cli           *client.BitgetApiClient
	Order         *Order
	ProfitRate    float64
	ApiKey        string
	SecretKey     string
	Passphrase    string
	Symbol        string
	Amout         string
	PriceScale    float64
	QuantityScale float64
	ActionTime    string
	SubTime       int64
}

func (c *GrabbingTask) Init() {
	c.SpotCli = new(v1.SpotOrderClient).Init()
	c.SpotCli.BitgetRestClient.ApiKey = c.ApiKey
	c.SpotCli.BitgetRestClient.ApiSecretKey = c.SecretKey
	c.SpotCli.BitgetRestClient.Passphrase = c.Passphrase
	c.SpotCli.BitgetRestClient.Signer = new(common.Signer).Init(c.SecretKey)

	c.Cli = new(client.BitgetApiClient).Init()
	c.Cli.BitgetRestClient.ApiKey = c.ApiKey
	c.Cli.BitgetRestClient.ApiSecretKey = c.SecretKey
	c.Cli.BitgetRestClient.Passphrase = c.Passphrase
	c.Cli.BitgetRestClient.Signer = new(common.Signer).Init(c.SecretKey)

}

func (c *GrabbingTask) Ready() {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}
	if c.ActionTime == "" {
		return
	}
	localTime, err := time.ParseInLocation(layout, c.ActionTime, location)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return
	}
	fmt.Println("Parsed time in UTC+8:", localTime)
	//100  300 millsec
	sleep := localTime.Sub(time.Now()).Milliseconds() - c.SubTime
	if sleep < 0 {
		return
	}
	fmt.Println("sleep(ms):", sleep)
	time.Sleep(time.Duration(sleep) * time.Millisecond)
}

func (c *GrabbingTask) MakeBuyOrder() error {
	params := make(map[string]string)
	//params["symbol"] = "WENUSDT_SPBL"
	params["symbol"] = fmt.Sprintf("%s_SPBL", c.Symbol)
	params["side"] = "buy"
	params["orderType"] = "market"
	params["force"] = "normal"
	params["quantity"] = c.Amout

	resp, err := c.SpotCli.PlaceOrder(params)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if gjson.Parse(resp).Get("code").String() != "00000" {
		err = fmt.Errorf("err when buy order,%s", gjson.Parse(resp).String())
		fmt.Println(err.Error())
		return err
	}
	orderId := gjson.Parse(resp).Get("data.orderId").String()
	fmt.Printf("success orderId:%s,resp:%s \n", orderId, resp)
	fmt.Println("===================================")
	c.Order = &Order{OrderId: orderId}
	return err
}

func (c *GrabbingTask) GetOrder() error {
	if c.Order == nil {
		return errors.New("err: orderId can not be empty")
	}
	params := make(map[string]string)
	params["orderId"] = c.Order.OrderId

	resp, err := c.Cli.Post("/api/spot/v1/trade/orderInfo", params)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Printf("get order respone:%s\n", resp)
	r := gjson.Parse(resp)
	if r.Get("code").String() != "00000" {
		err = fmt.Errorf("err when query order,%s", r.String())
		fmt.Printf("err:%s\n", err.Error())
		return err
	}
	dataArr := r.Get("data").Array()
	if len(dataArr) <= 0 {
		return fmt.Errorf("err: get order info,orderId:%s", c.Order.OrderId)
	}
	subR := r.Get("data").Array()[0]
	status := subR.Get("status").String()
	if status != "full_fill" {
		return fmt.Errorf("err: get order info,orderId:%s,status:%s", c.Order.OrderId, status)
	}
	priceStr := subR.Get("fillPrice").String()
	quantityStr := subR.Get("fillQuantity").String()
	fmt.Printf("get order info:%s\n", r.String())
	fmt.Println("===================================")
	c.Order.SellPrice = priceStr
	c.Order.SellQuantity = quantityStr

	return nil
}

func (c *GrabbingTask) MakeSellOrder() error {
	priceStr := utils.CaculteSellPrice(c.PriceScale, c.Order.SellPrice, c.ProfitRate+1)
	qantityeStr := utils.CaculteSellQantity(c.QuantityScale, c.Order.SellQuantity)
	params := make(map[string]string)
	//params["symbol"] = "WENUSDT_SPBL"
	params["symbol"] = fmt.Sprintf("%s_SPBL", c.Symbol)
	params["side"] = "sell"
	params["orderType"] = "limit"
	params["force"] = "normal"
	params["price"] = priceStr
	params["quantity"] = qantityeStr

	resp, err := c.SpotCli.PlaceOrder(params)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Printf("sell order,priceStr:%s,qantityeStr:%s\n", priceStr, qantityeStr)
	fmt.Printf("success sell order:%s\n", resp)
	fmt.Println("===================================")
	return nil
}

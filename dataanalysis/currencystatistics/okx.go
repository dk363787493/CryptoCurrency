package currency

import (
	"CryptoCurrency/common"
	"CryptoCurrency/model"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"sync"
)

type OkxCurrency struct {
	*BaseExchange
	ExchangeBit  uint16
	ExchangeName string
	Url          string
	WsURL        string
	dataCh       chan KlineData
	Cache        *sync.Map
}

func (*OkxCurrency) GetExchangeName() string {
	return "Okx"
}

func (okx *OkxCurrency) GetCurrencyInfoFromExchange() {
	list := okx.GetCurrencyList()
	for _, l := range list {
		c := model.CurrencyInfo{
			Currency:    l,
			ExchangeBit: common.OKxBit,
		}
		err := model.UpSertCurrencyInfo(c)
		if err != nil {
			slog.Error(err.Error())
			common.StackInfo(err)
		}
	}
}

func (okx *OkxCurrency) GetCurrencyList() []string {

	url := fmt.Sprint(okx.Url, "/public/instruments?instType=SPOT")
	resp, err := http.Get(url)
	if err != nil {
		slog.Error(fmt.Sprintf("err:%v", err))
		return nil
	}
	dataBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error(fmt.Sprintf("err:%v", err))
		return nil
	}
	data := gjson.ParseBytes(dataBytes)
	currencyList := data.Get("data").Array()
	result := make([]string, 0)

	for _, currency := range currencyList {
		if currency.Get("quoteCcy").String() == "USDT" {
			result = append(result, currency.Get("baseCcy").String())
			slog.Info(fmt.Sprintf("currency:%s", currency.Get("baseCcy").String()))
		}
	}
	return result
}

func (okx *OkxCurrency) GetMarketDataStreamFromExchange(currency []string) error {
	// Create a new WebSocket connection to the OKX server
	c, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s/business", okx.WsURL), nil)
	if err != nil {
		slog.Error(fmt.Sprintf("dial:%v", err))
		return err
	}

	agrs := make([]map[string]string, 0)
	for _, c := range currency {
		instId := c + "-USDT"
		agrs = append(agrs, map[string]string{"channel": "candle1s", "instId": instId})
	}
	// Subscribe to the BTC/USDT ticker channel
	subscribe := map[string]interface{}{
		"op": "subscribe",
		//"args": []map[string]string{{"channel": "candle1s", "instId": "BTC-USDT"}, {"channel": "candle1s", "instId": "ETH-USDT"}},
		"args": agrs,
	}
	if err := c.WriteJSON(subscribe); err != nil {
		log.Fatal("subscribe:", err)
	}

	done := make(chan struct{})

	// Start reading messages from the WebSocket connection
	go func() {
		defer close(done)
		defer c.Close()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			r := gjson.ParseBytes(message)
			data := r.Get("data").Array()
			if len(data) == 0 {
				continue
			}
			instId := r.Get("arg.instId").String()
			symbolIndex := strings.Index(instId, "-")
			symbol := instId[0:symbolIndex]
			price := data[0].Array()[1].Float()
			value, ok := okx.Cache.Load(symbol)
			if !ok || value.(float64) == price {
				okx.dataCh <- KlineData{
					Exchange: Binance,
					Symbol:   symbol,
					Price:    price,
				}
			}
		}
	}()
	return nil
}

func (okx *OkxCurrency) UpdateCache(symbol string, price float64) error {
	okx.Cache.Store(symbol, price)
	return nil
}

func (okx *OkxCurrency) AsyncUpdateSymbolPrice() {
	go func() {
		for d := range okx.dataCh {
			okx.UpdateCache(d.Symbol, d.Price)
		}
	}()
}
func (okx *OkxCurrency) GetPrice(symbol string) (float64, error) {
	price, ok := okx.Cache.Load(symbol)
	if !ok {
		return 0.0, fmt.Errorf("no symbol:%s", symbol)
	}

	return price.(float64), nil
}

func NewOkxCurrency(dataCh chan KlineData) Exchange {
	return &OkxCurrency{
		ExchangeBit:  common.OKxBit,
		ExchangeName: "Okx",
		Url:          "https://www.okx.com/api/v5",
		WsURL:        "wss://ws.okx.com:8443/ws/v5",
		dataCh:       dataCh,
		Cache:        new(sync.Map),
	}
}

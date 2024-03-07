package currency

import (
	"CryptoCurrency/common"
	"CryptoCurrency/config"
	"CryptoCurrency/model"
	"context"
	"fmt"
	binance_connector "github.com/binance/binance-connector-go"
	"github.com/tidwall/gjson"
	"log/slog"
	"strconv"
	"strings"
	"sync"
)

type BinanceCurrency struct {
	*BaseExchange
	ExchangeBit  uint16
	ExchangeName string
	Url          string
	dataCh       chan KlineData
	Cache        *sync.Map
}

func (*BinanceCurrency) GetExchangeName() string {
	return "Binance"
}

func (bc *BinanceCurrency) GetCurrencyInfoFromExchange() {
	list := bc.GetCurrencyList()
	for _, l := range list {
		c := model.CurrencyInfo{
			Currency:    l,
			ExchangeBit: common.BinanceBit,
		}
		fmt.Println("currency:", l)
		err := model.UpSertCurrencyInfo(c)
		if err != nil {
			slog.Error(err.Error())
			common.StackInfo(err)
		}
	}
}

func (bc *BinanceCurrency) GetCurrencyList() []string {
	client := binance_connector.NewClient(config.BinanceApiKey, config.BinanceSecretKey, bc.Url)
	exchangeInfo, err := client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	result := make([]string, 0)
	r := gjson.Parse(binance_connector.PrettyPrint(exchangeInfo))
	symbols := r.Get("symbols").Array()
	for _, s := range symbols {
		quoteAsset := s.Get("quoteAsset").String()
		if quoteAsset == "USDT" {
			baseAsset := s.Get("baseAsset").String()
			result = append(result, baseAsset)
		}
	}
	return result
}
func (bc *BinanceCurrency) GetMarketDataStreamFromExchange(currencyList []string) error {
	websocketStreamClient := binance_connector.NewWebsocketStreamClient(false)
	wsKlineHandler := func(event *binance_connector.WsKlineEvent) {
		symbol := event.Symbol
		priceStr := event.Kline.Open
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			slog.Error(fmt.Sprintf("symbol:%s,err:%v", symbol, err))
			return
		}
		symbol = symbol[:strings.Index(symbol, "USDT")]
		value, ok := bc.Cache.Load(symbol)
		if !ok || value.(float64) == price {
			bc.dataCh <- KlineData{
				Exchange: Binance,
				Symbol:   symbol,
				Price:    price,
			}
		}

	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	for _, c := range currencyList {
		go func(currency string) {
			symbol := currency + "USDT"
			slog.Info("ready to get kline data symbol:%s,exchange:binance")
			doneCh, _, err := websocketStreamClient.WsKlineServe(symbol, "1s", wsKlineHandler, errHandler)
			if err != nil {
				slog.Error(fmt.Sprintf("symbol:%s,err:%v", currency, err))
				return
			}
			slog.Info(fmt.Sprintf("getting kline data symbol:%s,exchange:binance", symbol))
			<-doneCh
			slog.Info(fmt.Sprintf("finish getting kline data symbol:%s", symbol))
		}(c)
	}

	return nil
}

func (bc *BinanceCurrency) UpdateCache(symbol string, price float64) error {
	bc.Cache.Store(symbol, price)
	return nil
}

func (bc *BinanceCurrency) GetPrice(symbol string) (float64, error) {
	price, ok := bc.Cache.Load(symbol)
	if !ok {
		return 0.0, fmt.Errorf("no symbol:%s", symbol)
	}

	return price.(float64), nil
}

func (bc *BinanceCurrency) AsyncUpdateSymbolPrice() {
	go func() {
		for d := range bc.dataCh {
			bc.UpdateCache(d.Symbol, d.Price)
		}
	}()
}

func NewBinanceCurrency(dataCh chan KlineData) Exchange {
	return &BinanceCurrency{
		ExchangeBit:  common.BinanceBit,
		ExchangeName: "Binance",
		Url:          "https://api.binance.com",
		dataCh:       dataCh,
		Cache:        new(sync.Map),
	}
}

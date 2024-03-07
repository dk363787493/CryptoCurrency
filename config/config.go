package config

import (
	"github.com/spf13/viper"
	"log"
)

const (
	BinanceApiKey    = "UHX1iqrKxSS40o7KZkzRAOc3ABpfFG3fzidEY5fTVkTdTtuQ6FLk0FTCLUFIQNEA"
	BinanceSecretKey = "PjmFJPSz19Rb5EooVkGE3RZolfJWCSfRCw0zwRpJy1vgD0HgCX3HIbTVeoU3tPML"
)
const (
	BinanceBit uint16 = 1
	OKxBit     uint16 = 1 << 1
)

var Configuration CryptoConfig

type CryptoConfig struct {
	Mysql            MysqlConfig
	BinanceApiKey    string
	BinanceSecretKey string
	BinanceBaseURL   string
	//Trade   map[string]string
	//Binance map[string]BinanceConfig
	//RiskControlRule RiskControlRule `yaml:"riskControlRule"`
}

//type RiskControlRule struct {
//	Aml   Aml
//	Fraud Fraud
//}

//type Aml struct {
//	BaseInfo  AmlBaseInfo  `yaml:"baseInfo"`
//	TradeInfo AmlTradeInfo `yaml:"tradeInfo"`
//}
//
//type AmlBaseInfo struct {
//	Age                       map[string]int `yaml:"age"`
//	Occupation                map[string]int `yaml:"occupation"`
//	ExpectedYearlyTradeVolume map[string]int `yaml:"expectedYearlyTradeVolume"`
//	Industry                  map[string]int `yaml:"industry"`
//	Citizenship               map[string]int `yaml:"citizenship"`
//	Pep                       map[string]int `yaml:"pep"`
//	Sanctions                 map[string]int `yaml:"sanctions"`
//	AdverseMedia              map[string]int `yaml:"adverseMedia"`
//	//sourceOfFunds
//}
//type AmlTradeInfo struct {
//	PendingOrders                map[int]int     `yaml:"pendingOrders"`
//	Total12mTradeVolume          map[int]int     `yaml:"total12mTradeVolume"`
//	LengthOfBusinessRelationship map[int]int     `yaml:"lengthOfBusinessRelationship"`
//	PayerBankCountry             int             `yaml:"payerBankCountry"`
//	scoreTradeLimit              map[int]float64 `yaml:"scoreTradeLimit"`
//}
//
//type Fraud struct {
//	StepOne   FraudStepOne   `yaml:"stepOne"`
//	StepTwo   FraudStepTwo   `yaml:"stepTwo"`
//	StepThree FraudStepThree `yaml:"stepThree"`
//	StepFour  FraudStepFour  `yaml:"stepFour"`
//}
//
//type FraudStepOne struct {
//	Binance map[string]BinanceCase `yaml:"binance"`
//	Okx     map[string]OkxCase     `yaml:"okx"`
//}
//
//type BinanceCase struct {
//	RegisterDay                    int `yaml:"registerDay"`
//	TraadedBefore30days            int `yaml:"traadedBefore30days"`
//	CompletedOrderNumOfLatest30day int `yaml:"completedOrderNumOfLatest30day"`
//}
//
//type OkxCase struct {
//	RegisterDay       float64 `yaml:"registerDay"`
//	CompletedOrderNum int     `yaml:"completedOrderNum"`
//}
//
//type FraudStepTwo struct {
//	Age map[string]int `yaml:"age"`
//}
//type FraudStepThree struct {
//	EarliestOrder map[string]int `yaml:"earliestOrder"`
//}
//type FraudStepFour struct {
//	IsGuided map[string]int `yaml:"isGuided"`
//}
//
//type BinanceConfig struct {
//	RegisterDay                    float64 `mapstructure:"registerDay"`
//	TradedBefore30Days             int     `mapstructure:"tradedBefore30days"`
//	CompletedOrderNumOfLatest30Day int     `mapstructure:"completedOrderNumOfLatest30day"`
//}

type RiskControlRule struct {
	Aml   Aml
	Fraud Fraud
}

type Aml struct {
	BaseInfo  AmlBaseInfo  `yaml:"baseInfo"`
	TradeInfo AmlTradeInfo `yaml:"tradeInfo"`
}

type AmlBaseInfo struct {
	Age                       map[string]int `yaml:"age"`
	Occupation                map[string]int `yaml:"occupation"`
	ExpectedYearlyTradeVolume map[string]int `yaml:"expectedYearlyTradeVolume"`
	Industry                  map[string]int `yaml:"industry"`
	Citizenship               map[string]int `yaml:"citizenship"`
	Pep                       map[string]int `yaml:"pep"`
	Sanctions                 map[string]int `yaml:"sanctions"`
	AdverseMedia              map[string]int `yaml:"adverseMedia"`
	//sourceOfFunds
}
type AmlTradeInfo struct {
	PendingOrders                map[int]int     `yaml:"pendingOrders"`
	Total12mTradeVolume          map[int]int     `yaml:"total12mTradeVolume"`
	LengthOfBusinessRelationship map[int]int     `yaml:"lengthOfBusinessRelationship"`
	PayerBankCountry             int             `yaml:"payerBankCountry"`
	scoreTradeLimit              map[int]float64 `yaml:"scoreTradeLimit"`
}

type Fraud struct {
	StepOne   FraudStepOne   `yaml:"stepOne"`
	StepTwo   FraudStepTwo   `yaml:"stepTwo"`
	StepThree FraudStepThree `yaml:"stepThree"`
	StepFour  FraudStepFour  `yaml:"stepFour"`
}

type FraudStepOne struct {
	Binance map[string]BinanceCase `yaml:"binance"`
	Okx     map[string]OkxCase     `yaml:"okx"`
}

type BinanceCase struct {
	RegisterDay                    float64 `yaml:"registerDay"`
	TraadedBefore30days            int     `yaml:"traadedBefore30days"`
	CompletedOrderNumOfLatest30day int     `yaml:"completedOrderNumOfLatest30day"`
}

type OkxCase struct {
	RegisterDay       float64 `yaml:"registerDay"`
	CompletedOrderNum int     `yaml:"completedOrderNum"`
}

type FraudStepTwo struct {
	Age map[string]int `yaml:"age"`
}
type FraudStepThree struct {
	EarliestOrder map[string]int `yaml:"earliestOrder"`
}
type FraudStepFour struct {
	IsGuided map[string]int `yaml:"isGuided"`
}

type MysqlConfig struct {
	Dsn      string
	Password string
}

func init() {
	// 初始化 viper
	viper.SetConfigName("prod")     // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")     // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./config") // 当前目录中查找配置文件
	// 读取配置数据
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&Configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}

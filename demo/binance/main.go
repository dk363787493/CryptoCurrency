package main

import (
	"context"
	"fmt"
	binance_connector "github.com/binance/binance-connector-go"
)

func main() {
	apiKey := "UHX1iqrKxSS40o7KZkzRAOc3ABpfFG3fzidEY5fTVkTdTtuQ6FLk0FTCLUFIQNEA"
	_ = apiKey
	secretKey := "PjmFJPSz19Rb5EooVkGE3RZolfJWCSfRCw0zwRpJy1vgD0HgCX3HIbTVeoU3tPML"
	_ = secretKey
	baseURL := "https://api.binance.com"

	// Initialise the client
	client := binance_connector.NewClient(apiKey, secretKey, baseURL)

	// FundingWalletService - /sapi/v1/asset/get-funding-asset
	//fundingWallet, err := client.NewFundingWalletService().Asset("USDT").
	//	Do(context.Background())
	//if err != nil {
	//	fmt.Println("err:", err.Error())
	//	return
	//}

	//accountSnapshot, err := client.NewGetAccountSnapshotService().MarketType("SPOT").
	//	Do(context.Background())
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(binance_connector.PrettyPrint(accountSnapshot))

	//WithdrawService - /sapi/v1/capital/withdraw/apply
	//withdraw, err := client.NewWithdrawService().Network("SOL").WalletType(1).Coin("USDT").Address("FAx2T4G4tQ1sZAYjj5jeSm2RL5Eo2v5ANHqoHvXeRX4R").
	//	Amount(2.00).Do(context.Background())
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(binance_connector.PrettyPrint(withdraw))

	// GetAllCoinsInfoService - /sapi/v1/capital/config/getall
	//allCoinsInfo, err := client.NewGetAllCoinsInfoService().Do(context.Background())
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(binance_connector.PrettyPrint(allCoinsInfo))

	// WithdrawHistoryService - /sapi/v1/capital/withdraw/history
	withdrawHistory, err := client.NewWithdrawHistoryService().
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(binance_connector.PrettyPrint(withdrawHistory))
}

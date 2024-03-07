package model

import (
	"CryptoCurrency/common"
	"fmt"
	"log/slog"
	"time"
)

type OpenTradingCurrency struct {
	ID                    uint      `gorm:"primarykey"`
	CrytoCurrency         string    `gorm:"column:cryto_currency"`
	BaseCurrency          string    `gorm:"column:base_currency"`
	QuoteOrderQty         float64   `gorm:"column:quote_order_qty"`
	ProfitRate            float64   `gorm:"column:profit_rate"`
	Status                string    `gorm:"column:status"`
	StatustartmissionDate time.Time `gorm:"column:startmission_date"`
	CreatedAt             time.Time `gorm:"column:created_at"`
}

func (*OpenTradingCurrency) TableName() string {
	return "opentradingcurrency"
}

func GetAllOpenTradingCurrency(status string) []OpenTradingCurrency {
	db := GetMysqlDB()
	var openTradingCurrency *OpenTradingCurrency
	whereSql := ""
	if status != "" {
		whereSql = fmt.Sprintf("where status='%s'", status)
	}
	sql := fmt.Sprintf("select cryto_currency,base_currency,profit_rate,quote_order_qty,startmission_date from %s %s ", openTradingCurrency.TableName(), whereSql)
	tx := db.Raw(sql)
	rows, err := tx.Rows()
	if err != nil {
		slog.Error(err.Error())
		common.StackInfo(err)
		return nil
	}
	defer rows.Close()
	result := make([]OpenTradingCurrency, 0)
	for rows.Next() {
		var r OpenTradingCurrency
		if err = rows.Scan(&r.CrytoCurrency, &r.BaseCurrency, &r.ProfitRate, &r.QuoteOrderQty, &r.StatustartmissionDate); err != nil {
			slog.Error(err.Error())
			common.StackInfo(err)
			break
		}
		result = append(result, r)
	}
	return result
}

func UpdadteStatus(cryto string, base string, status string) error {
	db := GetMysqlDB()
	var openTradingCurrency *OpenTradingCurrency
	sql := fmt.Sprintf("update %s set status=? where cryto_currency=? and base_currency=?", openTradingCurrency.TableName())
	result := db.Exec(sql, status, cryto, base)
	if result.Error != nil {
		// 错误处理
		fmt.Printf("Error occurred while updating: %v\n", result.Error)
		return result.Error
	}
	return nil
}

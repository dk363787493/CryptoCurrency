package model

import (
	"CryptoCurrency/common"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

type CurrencyInfo struct {
	ID          uint   `gorm:"primarykey"`
	Currency    string `gorm:"column:currency"`
	ExchangeBit uint16 `gorm:"column:exchange_bit"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (*CurrencyInfo) TableName() string {
	return "currency_info"
}

func GetAllCurrencyInfo(exchangeBit uint16) []CurrencyInfo {
	db := GetMysqlDB()
	var currencyInfo *CurrencyInfo
	whereSql := ""
	if exchangeBit != 0 {
		whereSql = fmt.Sprintf("where exchange_bit=%d", exchangeBit)
	}
	sql := fmt.Sprintf("select currency, exchange_bit from %s %s limit 10", currencyInfo.TableName(), whereSql)
	tx := db.Raw(sql)
	rows, err := tx.Rows()
	if err != nil {
		slog.Error(err.Error())
		common.StackInfo(err)
		return nil
	}
	defer rows.Close()
	result := make([]CurrencyInfo, 0)
	for rows.Next() {
		var r CurrencyInfo
		if err = rows.Scan(&r.Currency, &r.ExchangeBit); err != nil {
			slog.Error(err.Error())
			common.StackInfo(err)
			break
		}
		result = append(result, r)
	}
	return result
}

func UpSertCurrencyInfo(row CurrencyInfo) error {
	db := GetMysqlDB()
	row.Currency = strings.ToUpper(row.Currency)
	currencyInfoRow := new(CurrencyInfo)
	tx := db.First(&currencyInfoRow, "currency=?", row.Currency)
	if tx.Error != nil {
		//insert
		tx = db.Create(&row)
		affected := tx.RowsAffected
		if affected <= 0 {
			return fmt.Errorf("can not insert data to mysql db,raw row:%v", row)
		}
		return nil
	}
	currencyInfoRow.ExchangeBit = row.ExchangeBit | currencyInfoRow.ExchangeBit
	err := tx.Updates(currencyInfoRow).Error
	return err
}

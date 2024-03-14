package utils

import "github.com/shopspring/decimal"

func CaculteSellPrice(scale float64, price string, factor float64) string {
	priceDecimal, _ := decimal.NewFromString(price)
	scaleDecimal := decimal.NewFromFloat(scale)
	pp := priceDecimal.Mul(decimal.NewFromFloat(factor))
	i := pp.Div(scaleDecimal).IntPart()
	dd := decimal.NewFromInt(i).Mul(scaleDecimal)
	return dd.String()
}

func CaculteSellQantity(scale float64, qantity string) string {
	qantityDecimal, _ := decimal.NewFromString(qantity)
	scaleDecimal := decimal.NewFromFloat(scale)
	// 扣除0.1%点
	pp := qantityDecimal.Mul(decimal.NewFromFloat(0.999))
	i := pp.Div(scaleDecimal).IntPart()
	dd := decimal.NewFromInt(i).Mul(scaleDecimal)
	return dd.String()
}

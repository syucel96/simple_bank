package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/syucel96/simplebank/util"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}

var validAmount validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if amount, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsDecimal(amount)
	}
	return false
}

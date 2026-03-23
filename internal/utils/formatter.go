package utils

import (
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func FormatCurrency(amount float64, currencyCode string) string {
	locale := getLocale(currencyCode)

	tag := language.Make(locale)
	p := message.NewPrinter(tag)

	unit, err := currency.ParseISO(currencyCode)
	if err != nil {
		return ""
	}

	return p.Sprintf("%v", currency.Symbol(unit.Amount(amount)))

	// return p.Sprintf("%v", unit.Amount(amount))
}

func getLocale(currencyCode string) string {
	switch currencyCode {
	case "IDR":
		return "id-ID"
	default:
		return "en-US"
	}
}

package bizmandu

import (
	"fmt"
	"nepse-backend/utils"
	"net/http"
	"os"
)

const (
	All             = "tickers/all"
	Header          = "tearsheet/header"
	Health          = "overview/topGainers/?count=5"
	Summary         = "tearsheet/summary"
	Financial       = "tearsheet/financial/keyStats"
	BalanceSheet    = "tearsheet/financial/balanceSheet"
	PriceHistory    = "tradingView/history"
	IncomeStatement = "tearsheet/financial/incomeStatement"
	Dividend        = "tearsheet/dividend/"
	MutualFund      = "tearsheet/summaryMf/"
)

type BizmanduAPI struct {
	client *utils.Client
}

func NewBizmandu() (*BizmanduAPI, error) {
	client := utils.NewClient(nil, os.Getenv("BIZMANDU"))

	_, err := client.NewRequest(http.MethodGet, Health, nil)

	if err != nil {
		return nil, err
	}

	biz := &BizmanduAPI{
		client: client,
	}
	return biz, nil
}

func (b *BizmanduAPI) buildTickerSlug(urlPath, ticker string) string {
	return fmt.Sprintf("%s/?tkr=%s", urlPath, ticker)
}

func (b *BizmanduAPI) buildPriceHistorySlug(urlPath, ticker string, from, to int64) string {
	return fmt.Sprintf("%s?symbol=%s&resolution=D&from=%d&to=%d", urlPath, ticker, from, to)
}

//https://bizmandu.com/__stock/tradingView/history?symbol=BPCL&resolution=D&from=1626880096&to=1627744096

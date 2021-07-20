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
	IncomeStatement = "tearsheet/financial/incomeStatement"
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

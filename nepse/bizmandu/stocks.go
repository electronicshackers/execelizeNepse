package bizmandu

import (
	"context"
	"nepse-backend/nepse"
	"nepse-backend/nepse/neweb"
	"nepse-backend/utils"
	"net/http"
	"strings"
)

type Tickers struct {
	Response float64        `json:"response"`
	Error    string         `json:"error"`
	Message  []nepse.Ticker `json:"message"`
}

func (b *BizmanduAPI) GetStocks() ([]nepse.Ticker, error) {
	req, err := b.client.NewRequest(http.MethodGet, All, nil)
	if err != nil {
		return nil, err
	}

	res := &Tickers{}
	if _, err := b.client.Do(context.Background(), req, res); err != nil {
		return nil, err
	}

	var stocks []nepse.Ticker

	for _, ticker := range res.Message {
		stocks = append(stocks, nepse.Ticker{
			Ticker:      ticker.Ticker,
			Companyname: ticker.Companyname,
			Sector:      ticker.Sector,
		})
	}

	return stocks, nil
}

func (b *BizmanduAPI) GetSectorStock(sector string) ([]nepse.Ticker, error) {
	tickers, err := b.GetStocks()
	if err != nil {
		return nil, err
	}

	nep, err := neweb.Neweb()

	if err != nil {
		return nil, err
	}

	allStocks, err := nep.GetStocks()
	if err != nil {
		return nil, err
	}

	var stocks []nepse.Ticker

	for _, ticker := range tickers {
		if ticker.Sector == sector {
			if !strings.Contains(ticker.Companyname, "Promoter") {
				stocks = append(stocks, nepse.Ticker{
					Ticker:      ticker.Ticker,
					Companyname: ticker.Companyname,
					Sector:      ticker.Sector,
				})
			}
		}
	}

	return utils.SetIntersection(allStocks, stocks), nil
}

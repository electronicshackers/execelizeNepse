package bizmandu

import (
	"context"
	"nepse-backend/nepse"
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
			Ticker:          ticker.Ticker,
			Companyname:     ticker.Companyname,
			Sector:          ticker.Sector,
			Lasttradedprice: ticker.Lasttradedprice,
		})
	}

	return stocks, nil
}

func (b *BizmanduAPI) GetSectorStock(sector string, newWeb, biz []nepse.Ticker) ([]nepse.Ticker, error) {
	var stocks []nepse.Ticker

	for _, ticker := range biz {
		if ticker.Sector == sector {
			if !strings.Contains(ticker.Companyname, "Promoter") {
				stocks = append(stocks, nepse.Ticker{
					Ticker:          ticker.Ticker,
					Companyname:     ticker.Companyname,
					Sector:          ticker.Sector,
					Lasttradedprice: ticker.Lasttradedprice,
				})
			}
		}
	}
	return utils.SetIntersection(newWeb, stocks), nil
}

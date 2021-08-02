package neweb

import (
	"context"
	"nepse-backend/nepse"
	"net/http"
)

type ListedStocks []struct {
	Securityid          string      `json:"securityId"`
	Securityname        string      `json:"securityName"`
	Symbol              string      `json:"symbol"`
	Indexid             int         `json:"indexId"`
	Openprice           float64     `json:"openPrice"`
	Highprice           float64     `json:"highPrice"`
	Lowprice            float64     `json:"lowPrice"`
	Totaltradequantity  int         `json:"totalTradeQuantity"`
	Lasttradedprice     float64     `json:"lastTradedPrice"`
	Percentagechange    float64     `json:"percentageChange"`
	Lastupdateddatetime string      `json:"lastUpdatedDateTime"`
	Lasttradedvolume    interface{} `json:"lastTradedVolume"`
	Previousclose       float64     `json:"previousClose"`
}

func (n *NewebAPI) GetStocks() ([]nepse.Ticker, error) {
	req, err := n.client.NewRequest(http.MethodGet, All, nil)
	if err != nil {
		return nil, err
	}

	res := &ListedStocks{}
	if _, err := n.client.Do(context.Background(), req, res); err != nil {
		return nil, err
	}

	var stocks []nepse.Ticker

	for _, ticker := range *res {
		stocks = append(stocks, nepse.Ticker{
			Ticker:          ticker.Symbol,
			Companyname:     ticker.Securityname,
			Id:              ticker.Securityid,
			Lasttradedprice: ticker.Lasttradedprice,
		})
	}

	return stocks, nil
}

func (n *NewebAPI) GetMutualFundStock() ([]nepse.Ticker, error) {
	req, err := n.client.NewRequest(http.MethodGet, MutualFund, nil)
	if err != nil {
		return nil, err
	}

	res := &ListedStocks{}
	if _, err := n.client.Do(context.Background(), req, res); err != nil {
		return nil, err
	}

	var stocks []nepse.Ticker

	for _, ticker := range *res {
		stocks = append(stocks, nepse.Ticker{
			Ticker:      ticker.Symbol,
			Companyname: ticker.Securityname,
			Id:          ticker.Securityid,
		})
	}

	return stocks, nil
}

func (n *NewebAPI) GetSectorStock(sector string) ([]nepse.Ticker, error) {
	tickers, err := n.GetStocks()
	if err != nil {
		return nil, err
	}

	var stocks []nepse.Ticker

	for _, ticker := range tickers {
		if ticker.Sector == sector {
			stocks = append(stocks, nepse.Ticker{
				Ticker:          ticker.Ticker,
				Companyname:     ticker.Companyname,
				Id:              ticker.Id,
				Lasttradedprice: ticker.Lasttradedprice,
			})
		}
	}

	return stocks, nil
}

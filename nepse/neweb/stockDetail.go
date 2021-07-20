package neweb

import (
	"context"
	"nepse-backend/nepse"
	"net/http"
)

func (n *NewebAPI) GetCurrentPrice(ticker string) (*nepse.LastTradingDayStats, error) {
	req, err := n.client.NewRequest(http.MethodGet, All, nil)
	if err != nil {
		return nil, err
	}

	res := &ListedStocks{}
	if _, err := n.client.Do(context.Background(), req, res); err != nil {
		return nil, err
	}

	var currentPrice nepse.LastTradingDayStats

	for _, tickers := range *res {
		if tickers.Symbol == ticker {
			currentPrice = nepse.LastTradingDayStats{
				Ticker:              tickers.Symbol,
				Openprice:           tickers.Openprice,
				Lowprice:            tickers.Lowprice,
				Highprice:           tickers.Highprice,
				PointChanged:        tickers.Lasttradedprice - tickers.Previousclose,
				Totaltradequantity:  tickers.Totaltradequantity,
				Lasttradedprice:     tickers.Lasttradedprice,
				Percentagechange:    tickers.Percentagechange,
				Lastupdateddatetime: tickers.Lastupdateddatetime,
				Lasttradedvolume:    tickers.Lasttradedvolume,
				Previousclose:       tickers.Previousclose,
			}
		}
	}

	return &currentPrice, nil
}

func (n *NewebAPI) GetPriceHistory(ticker string) ([]nepse.PriceHistoryMinified, error) {
	url := n.buildHistorySlug(PriceHistory, ticker)
	req, err := n.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res := &nepse.PriceHistory{}

	if _, err := n.client.Do(context.Background(), req, res); err != nil {
		return nil, err
	}

	var prices []nepse.PriceHistoryMinified

	for _, content := range res.Content {
		prices = append(prices, nepse.PriceHistoryMinified{
			Date:         content.Businessdate,
			Price:        content.Closeprice,
			AveragePrice: content.Averagetradedprice,
			LowPrice:     content.Lowprice,
			HighPrice:    content.Highprice,
		})
	}

	return prices, nil
}

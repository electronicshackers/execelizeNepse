package smartwealth

import (
	"context"
	"fmt"
	"net/http"
)

type TechnicalResponse []struct {
	Stocksymbol       string      `json:"StockSymbol"`
	Companyname       string      `json:"CompanyName"`
	Ltp               float64     `json:"LTP"`
	Pricechange       float64     `json:"PriceChange"`
	Percentchange     float64     `json:"PercentChange"`
	Beta              float64     `json:"Beta"`
	Var5Percent30Days interface{} `json:"VAR5Percent30Days"`
	Macd              float64     `json:"MACD"`
	Rsi               float64     `json:"RSI"`
	Sma5              float64     `json:"SMA5"`
	Sma10             float64     `json:"SMA10"`
	Sma20             float64     `json:"SMA20"`
	Sma50             float64     `json:"SMA50"`
	Ema5              float64     `json:"EMA5"`
	Ema10             float64     `json:"EMA10"`
	Ema20             float64     `json:"EMA20"`
	Ema50             float64     `json:"EMA50"`
	Companyid         int         `json:"CompanyId"`
}

func (s *SmartWealthApi) GetTechnicalAnalysis(sectorId int) (*TechnicalResponse, error) {
	url := s.buildTechnicalSlug(Technical, sectorId)

	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res := &TechnicalResponse{}

	if _, err := s.client.Do(context.Background(), req, res); err != nil {
		fmt.Println("err", err)
		return nil, err
	}
	return res, nil
}

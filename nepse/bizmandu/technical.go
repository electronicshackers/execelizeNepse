package bizmandu

import (
	"context"
	"fmt"
	"nepse-backend/utils"
	"net/http"
	"time"
)

// https://bizmandu.com/__stock/tradingView/history?symbol=EIC&resolution=D&from=1601989046&to=1633093106

func (b *BizmanduAPI) GetTechnicalData(stock, resolution string) (*utils.TechnicalData, error) {
	now := time.Now()
	start := now.AddDate(-1, 0, -5)
	url := b.buildTickerSlugTechnicalURL(Technical, stock, resolution, start.Unix(), now.Unix())
	fmt.Println("url", url)
	req, err := b.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res := &utils.TechnicalData{}
	if _, err := b.client.Do(context.Background(), req, res); err != nil {
		return nil, err
	}

	return res, nil
}

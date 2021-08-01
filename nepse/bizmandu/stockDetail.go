package bizmandu

import (
	"context"
	"nepse-backend/nepse"
	"net/http"
	"time"
)

type CurrentPrice struct {
	Response int    `json:"response"`
	Error    string `json:"error"`
	Message  struct {
		Ticker           string    `json:"ticker"`
		Company          string    `json:"company"`
		Latestprice      float64   `json:"latestPrice"`
		Pointchange      float64   `json:"pointChange"`
		Percentagechange float64   `json:"percentageChange"`
		Timestamp        time.Time `json:"timestamp"`
		Wtavgprice       float64   `json:"wtAvgPrice"`
		Sharestraded     int       `json:"sharesTraded"`
		Volume           int       `json:"volume"`
		Mktcap           float64   `json:"mktCap"`
	} `json:"message"`
}

type PriceHistoryResponse struct {
	S string    `json:"s"`
	T []int     `json:"t"`
	C []float64 `json:"c"`
	O []float64 `json:"o"`
	H []float64 `json:"h"`
	L []float64 `json:"l"`
	V []int     `json:"v"`
}

type StockDetails struct {
	Response float64 `json:"response"`
	Error    string  `json:"error"`
	Message  struct {
		Keyfinancial struct {
			Ticker  string  `json:"ticker"`
			Year    string  `json:"year"`
			Quarter float64 `json:"quarter"`
			Data    []struct {
				Type         string  `json:"type"`
				Totalrevenue float64 `json:"totalRevenue"`
				Grossprofit  float64 `json:"grossProfit"`
				Netincome    float64 `json:"netIncome"`
				Eps          float64 `json:"eps"`
				Bvps         float64 `json:"bvps"`
				Roa          float64 `json:"roa"`
				Roe          float64 `json:"roe"`
			} `json:"data"`
		} `json:"keyFinancial"`
		Summary struct {
			Ticker           string  `json:"ticker"`
			Open             float64 `json:"open"`
			Avgvolume        float64 `json:"avgVolume"`
			Dayshigh         float64 `json:"daysHigh"`
			Dayslow          float64 `json:"daysLow"`
			Fiftytwoweekhigh float64 `json:"fiftyTwoWeekHigh"`
			Fiftytwoweeklow  float64 `json:"fiftyTwoWeekLow"`
			Listedshares     float64 `json:"listedShares"`
			Mktcap           float64 `json:"mktCap"`
			Epsdiluted       float64 `json:"epsDiluted"`
			Pediluted        float64 `json:"peDiluted"`
			Bvps             float64 `json:"bvps"`
			Beta             float64 `json:"beta"`
		} `json:"summary"`
	} `json:"message"`
}

type StockSummary struct {
	Response float64 `json:"response"`
	Error    string  `json:"error"`
	Message  struct {
		Keyfinancial struct {
			Ticker  string  `json:"ticker"`
			Year    string  `json:"year"`
			Quarter float64 `json:"quarter"`
			Data    []struct {
				Type         string  `json:"type"`
				Totalrevenue float64 `json:"totalRevenue"`
				Grossprofit  float64 `json:"grossProfit"`
				Netincome    float64 `json:"netIncome"`
				Eps          float64 `json:"eps"`
				Bvps         float64 `json:"bvps"`
				Roa          float64 `json:"roa"`
				Roe          float64 `json:"roe"`
			} `json:"data"`
		} `json:"keyFinancial"`
		Summary struct {
			Ticker           string  `json:"ticker"`
			Open             float64 `json:"open"`
			Avgvolume        float64 `json:"avgVolume"`
			Dayshigh         float64 `json:"daysHigh"`
			Dayslow          float64 `json:"daysLow"`
			Fiftytwoweekhigh float64 `json:"fiftyTwoWeekHigh"`
			Fiftytwoweeklow  float64 `json:"fiftyTwoWeekLow"`
			Listedshares     float64 `json:"listedShares"`
			Mktcap           float64 `json:"mktCap"`
			Epsdiluted       float64 `json:"epsDiluted"`
			Pediluted        float64 `json:"peDiluted"`
			Bvps             float64 `json:"bvps"`
			Beta             float64 `json:"beta"`
		} `json:"summary"`
	} `json:"message"`
}

type DividendHistory struct {
	Response int    `json:"response"`
	Error    string `json:"error"`
	Message  struct {
		Ticker   string `json:"ticker"`
		Dividend []struct {
			Year  string  `json:"year"`
			Cash  float64 `json:"cash"`
			Bonus float64 `json:"bonus"`
		} `json:"dividend"`
		Rights []struct {
			Date   int64   `json:"date"`
			Rights float64 `json:"rights"`
		} `json:"rights"`
	} `json:"message"`
}

func (b *BizmanduAPI) GetCurrentPrice(ticker string) (*nepse.LastTradingDayStats, error) {
	url := b.buildTickerSlug(Header, ticker)
	req, err := b.client.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	res := &CurrentPrice{}
	if _, err := b.client.Do(context.Background(), req, res); err != nil {
		return nil, err
	}

	currentPrice := nepse.LastTradingDayStats{
		Ticker:             res.Message.Ticker,
		Totaltradequantity: res.Message.Sharestraded,
		PointChanged:       res.Message.Pointchange,
		Lasttradedprice:    res.Message.Latestprice,
		Percentagechange:   res.Message.Percentagechange,
		Lasttradedvolume:   res.Message.Volume,
	}

	return &currentPrice, nil
}

func (b *BizmanduAPI) GetSummary(ticker string) (*StockSummary, error) {
	url := b.buildTickerSlug(Summary, ticker)

	req, err := b.client.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	res := &StockSummary{}
	if _, err := b.client.Do(context.Background(), req, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (b *BizmanduAPI) GetPriceHistory(ticker string, start, end int64) (*PriceHistoryResponse, error) {
	url := b.buildPriceHistorySlug(PriceHistory, ticker, start, end)

	req, err := b.client.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	res := &PriceHistoryResponse{}
	if _, err := b.client.Do(context.Background(), req, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (b *BizmanduAPI) GetDividendHistory(ticker string) (*DividendHistory, error) {
	url := b.buildTickerSlug(Dividend, ticker)
	req, err := b.client.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	res := &DividendHistory{}
	if _, err := b.client.Do(context.Background(), req, res); err != nil {
		return nil, err
	}

	return res, nil
}

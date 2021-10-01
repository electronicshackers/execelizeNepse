package controllers

import (
	"encoding/json"
	"io/ioutil"
	responses "nepse-backend/api/response"
	"nepse-backend/nepse/bizmandu"
	"nepse-backend/nepse/neweb"
	"net/http"
)

func (server *Server) GetTechnicalData(w http.ResponseWriter, r *http.Request) {
	var tickers neweb.ListedStocks
	stocks, err := ioutil.ReadFile("stocks.json")
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}

	err = json.Unmarshal([]byte(stocks), &tickers)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}

	biz, err := bizmandu.NewBizmandu()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rsiMap := make(map[string][]float64)

	macdMap := make(map[string][]float64)
	signalLineMap := make(map[string][]float64)
	histogramMap := make(map[string][]float64)

	for _, stock := range tickers[1:2] {
		data, err := biz.GetTechnicalData(stock.Symbol, "D")
		if err != nil {
			responses.ERROR(w, 400, err)
		}

		diff := data.Diff()
		gain := data.Gains(diff)
		loss := data.Losses(diff)

		averageGain := data.Average(gain, 14)
		averageLoss := data.Average(loss, 14)

		averageGains := data.MovingAverage(gain, 14, averageGain)
		averageLosses := data.MovingAverage(loss, 14, averageLoss)

		rs := data.RelativeStrength(averageLosses, averageGains)
		rsi := data.RelativeStrengthIndicator(rs)

		rsiMap[stock.Symbol] = rsi

		ema12 := data.ExponentialMovingAverage(data.C, 12, data.Average(data.C, 12), data.Multiplier(12))

		ema26 := data.ExponentialMovingAverage(data.C, 26, data.Average(data.C, 26), data.Multiplier(26))

		macd := data.MovingDifference(ema12, ema26, 14)
		signalLine := data.ExponentialMovingAverage(macd, 9, data.Average(macd, 9), data.Multiplier(9))
		histogram := data.MovingDifference(macd, signalLine, 8)
		macdMap[stock.Symbol] = macd
		signalLineMap[stock.Symbol] = signalLine
		histogramMap[stock.Symbol] = histogram
	}

	responses.JSON(w, 200, map[string]interface{}{"rsi": rsiMap, "macd": macdMap, "signalLine": signalLineMap, "histogram": histogramMap})
}

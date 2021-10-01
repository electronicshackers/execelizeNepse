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

	for _, stock := range tickers {
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
	}

	responses.JSON(w, 200, rsiMap)
}

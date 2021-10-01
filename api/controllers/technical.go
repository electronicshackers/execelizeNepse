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

		rsiMap[stock.Symbol] = data.RSI()
		macdMap[stock.Symbol], signalLineMap[stock.Symbol], histogramMap[stock.Symbol] = data.MACD()
	}

	responses.JSON(w, 200, map[string]interface{}{"rsi": rsiMap, "macd": macdMap, "signalLine": signalLineMap, "histogram": histogramMap})
}

package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	responses "nepse-backend/api/response"

	talib "github.com/markcheno/go-talib"

	"nepse-backend/nepse/bizmandu"
	"nepse-backend/nepse/neweb"
	"nepse-backend/utils"
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

	ema20Map := make(map[string][]float64)
	ema50Map := make(map[string][]float64)
	ema200Map := make(map[string][]float64)

	var keyLevels utils.KeyLevels

	for _, stock := range tickers[1:2] {
		data, err := biz.GetTechnicalData(stock.Symbol, "D")
		if err != nil {
			responses.ERROR(w, 400, err)
		}

		// rsi := talib.Rsi(data.C, 14)
		// fmt.Println("rsi", rsi)

		// macd, signalLine, histogram := talib.Macd(data.C, 12, 26, 9)
		// fmt.Println("macd", macd)
		// fmt.Println("signalLine", signalLine)
		// fmt.Println("histogram", histogram)

		// ema20, ema50, ema200 := talib.Ema(data.C, 20), talib.Ema(data.C, 50), talib.Ema(data.C, 200)
		// fmt.Println("ema20", ema20)
		// fmt.Println("ema50", ema50)
		// fmt.Println("ema200", ema200)

		// psar := talib.Sar(data.H, data.L, 0.02, 0.2)
		// fmt.Println("psar", psar)

		// willamsR := talib.WillR(data.H, data.L, data.C, 14)

		// fmt.Println("willamsR", willamsR)

		ad := talib.Ad(data.H, data.L, data.C, data.V)
		fmt.Println("ad", ad)

		rsiMap[stock.Symbol] = data.RSI()
		macdMap[stock.Symbol], signalLineMap[stock.Symbol], histogramMap[stock.Symbol] = data.MACD()
		ema20Map[stock.Symbol] = data.EMA(20)
		ema50Map[stock.Symbol] = data.EMA(50)
		ema200Map[stock.Symbol] = data.EMA(200)

		keyLevels = data.KeyLevels()

	}

	responses.JSON(w, 200, map[string]interface{}{
		"rsi":        rsiMap,
		"macd":       macdMap,
		"signalLine": signalLineMap,
		"histogram":  histogramMap,
		"ema20":      ema20Map,
		"ema50":      ema50Map,
		"ema200":     ema200Map,
		"keyLevels":  keyLevels,
	})
}

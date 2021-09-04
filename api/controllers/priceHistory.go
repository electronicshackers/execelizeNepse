package controllers

import (
	responses "nepse-backend/api/response"
	"nepse-backend/nepse"
	"nepse-backend/nepse/bizmandu"
	"nepse-backend/nepse/neweb"
	"nepse-backend/utils"
	"net/http"
	"os"
	"strings"
)

type Result struct {
	Ticker                  string  `json:"ticker"`
	CompanyName             string  `json:"companyName"`
	StartPrice              float64 `json:"startPrice"`
	EndPrice                float64 `json:"endPrice"`
	MaxPrice                float64 `json:"maxPrice"`
	MinPrice                float64 `json:"minPrice"`
	PercentageChange        float64 `json:"changed"`
	PointChange             float64 `json:"pointChange"`
	ExtremePointChange      float64 `json:"extremePointChange"`
	ExtremePercentageChange float64 `json:"extremePercentageChange"`
	MaxAveragePrice         float64 `json:"maxAveragePrice"`
	MinAveragePrice         float64 `json:"minAveragePrice"`
	AveragePointChange      float64 `json:"averagePointChange"`
	AveragePercentageChange float64 `json:"averagePercentageChange"`
}

func (server *Server) GetPriceHistory(w http.ResponseWriter, r *http.Request) {

	// Get query params from request
	params := r.URL.Query()
	start := params.Get("start")
	end := params.Get("end")
	sector := params.Get("sectors")

	if sector == "" {
		http.Error(w, "Sector is required", http.StatusBadRequest)
		return
	}

	// Comma separated string to a slice
	querySector := strings.Split(sector, ",")

	days, err := utils.GetDateRange(w, start, end)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	sectors := utils.MapColumns(querySector)

	nep, err := neweb.Neweb()

	if err != nil {
		http.Error(w, "Error in fetching data from Neweb", http.StatusBadRequest)
		return
	}

	biz, err := bizmandu.NewBizmandu()

	if err != nil {
		http.Error(w, "Error in fetching data from Bizmandu", http.StatusBadRequest)
		return
	}

	allStocks, err := biz.GetStocks()

	if err != nil {
		http.Error(w, "Error in getting stock list", http.StatusBadRequest)
		return
	}

	listedStocks, err := nep.GetStocks()

	if err != nil {
		http.Error(w, "Error in getting stock list", http.StatusBadRequest)
		return
	}

	var ticker []nepse.Ticker

	for _, bizStocks := range allStocks {
		for _, newStocks := range listedStocks {
			if bizStocks.Ticker == newStocks.Ticker {
				ticker = append(ticker, nepse.Ticker{
					Id:          newStocks.Id,
					Ticker:      newStocks.Ticker,
					Companyname: newStocks.Companyname,
					Sector:      bizStocks.Sector,
				})
			}
		}
	}

	for _, sec := range sectors {
		filtered := getFilteredTickers(ticker, sec)

		var resu []Result

		for _, fil := range filtered {
			re, err := nep.GetPriceHistory(fil.Id)

			if err != nil {
				http.Error(w, "Error in fetching data from Neweb", http.StatusBadRequest)
				return
			}

			result := preprocessChange(re, start, end, days)

			resu = append(resu, Result{
				Ticker:                  fil.Ticker,
				CompanyName:             fil.Companyname,
				StartPrice:              result.StartPrice,
				EndPrice:                result.EndPrice,
				PercentageChange:        result.PercentageChange,
				PointChange:             result.PointChange,
				MaxPrice:                result.MaxPrice,
				MinPrice:                result.MinPrice,
				MaxAveragePrice:         result.MaxAveragePrice,
				MinAveragePrice:         result.MinAveragePrice,
				ExtremePointChange:      result.ExtremePointChange,
				AveragePointChange:      result.AveragePointChange,
				AveragePercentageChange: result.AveragePercentageChange,
				ExtremePercentageChange: result.ExtremePercentageChange,
			})
		}

		categories := map[string]string{
			"A1": "Ticker",
			"B1": start, "C1": end, "D1": "Pt. Change", "E1": "% Change",
			"F1": "Max", "G1": "Min", "H1": "Max-Min", "I1": "%(Max-Min)",
			"J1": "Max. Avg", "K1": "Min. Avg", "L1": "Max-Min(avg)", "M1": "%Change(Avg)",
		}

		var excelVals []map[string]interface{}

		for k, v := range resu {
			excelVal := map[string]interface{}{
				utils.GetColumn("A", k): v.Ticker,
				utils.GetColumn("B", k): v.StartPrice, utils.GetColumn("C", k): v.EndPrice, utils.GetColumn("D", k): v.PointChange, utils.GetColumn("E", k): v.PercentageChange,
				utils.GetColumn("F", k): v.MaxPrice, utils.GetColumn("G", k): v.MinPrice, utils.GetColumn("H", k): v.ExtremePointChange, utils.GetColumn("I", k): v.ExtremePercentageChange,
				utils.GetColumn("J", k): v.MaxAveragePrice, utils.GetColumn("K", k): v.MinAveragePrice, utils.GetColumn("L", k): v.AveragePointChange, utils.GetColumn("M", k): v.AveragePercentageChange,
			}
			if v.Ticker != "" {
				excelVals = append(excelVals, excelVal)
			}
		}

		var folderName = start + end

		// Create a new folder if it doesn't exist
		if _, err := os.Stat(folderName); os.IsNotExist(err) {
			os.Mkdir(folderName, 0777)
		}

		utils.CreateExcelFile(folderName, sec, categories, excelVals)
	}

}

func getFilteredTickers(ticker []nepse.Ticker, sector string) []nepse.Ticker {
	var filteredStocks []nepse.Ticker
	for _, tick := range ticker {
		if tick.Sector == sector {
			filteredStocks = append(filteredStocks, tick)
		}
	}
	return filteredStocks
}

func preprocessChange(minHistory []nepse.PriceHistoryMinified, start, end string, days []string) Result {
	var closeStart, closeEnd, max, minP, maxAvg, minAvg, avg float64

	var maxPriceArray, minPriceArray, avgArray []float64

	for _, min := range minHistory {
		if min.Date == start {
			closeStart = min.Price
			avg = min.AveragePrice

		}

		if min.Date == end {
			closeEnd = min.Price
		}

		for _, day := range days {
			if day == min.Date {
				maxPriceArray = append(maxPriceArray, min.HighPrice)
				minPriceArray = append(minPriceArray, min.LowPrice)
				avgArray = append(avgArray, min.AveragePrice)
			}
		}
	}

	_, max = utils.MinMax(maxPriceArray)
	minP, _ = utils.MinMax(minPriceArray)
	minAvg, maxAvg = utils.MinMax(avgArray)

	changeClosed := closeEnd - closeStart

	changeExtreme := max - minP

	changeAverage := maxAvg - minAvg

	percentageChangeOnExtreme := (changeExtreme / closeStart) * 100

	percentageChangeOnAverage := (changeAverage / avg) * 100

	percentageChangedOnClosed := (changeClosed / closeStart) * 100

	return Result{
		StartPrice:              closeStart,
		EndPrice:                closeEnd,
		PercentageChange:        utils.ToFixed(percentageChangedOnClosed, 2),
		PointChange:             changeClosed,
		ExtremePointChange:      changeExtreme,
		AveragePointChange:      utils.ToFixed(changeAverage, 2),
		ExtremePercentageChange: utils.ToFixed(percentageChangeOnExtreme, 2),
		AveragePercentageChange: utils.ToFixed(percentageChangeOnAverage, 2),
		MaxPrice:                max,
		MinPrice:                minP,
		MaxAveragePrice:         maxAvg,
		MinAveragePrice:         minAvg,
	}

}

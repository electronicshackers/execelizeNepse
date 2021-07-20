package controllers

import (
	"fmt"
	"nepse-backend/nepse"
	"nepse-backend/nepse/bizmandu"
	"nepse-backend/nepse/neweb"
	"nepse-backend/utils"
	"net/http"
	"os"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
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

	sectors := utils.MapColumns(querySector)

	if start == "" || end == "" {
		http.Error(w, "Start and end date are required", http.StatusBadRequest)
		return
	}

	// Change String to Date
	startDate, err := utils.StringToTime(start)
	if err != nil {
		http.Error(w, "Start date is invalid", http.StatusBadRequest)
		return
	}

	startDay := startDate.Weekday().String()
	if startDay == "Friday" || startDay == "Saturday" {
		http.Error(w, "Start date should be a weekday", http.StatusBadRequest)
	}

	endDate, err := utils.StringToTime(end)
	if err != nil {
		http.Error(w, "End date is invalid", http.StatusBadRequest)
		return
	}

	endDay := endDate.Weekday().String()

	if endDay == "Friday" || endDay == "Saturday" {
		http.Error(w, "End date should be a weekday", http.StatusBadRequest)
	}

	// find the difference in days between start and end date
	diffDays := endDate.Sub(startDate).Hours() / 24
	if diffDays < 0 {
		http.Error(w, "Start date must be before end date", http.StatusBadRequest)
		return
	}
	if diffDays > 91 {
		http.Error(w, "Start date must be less than 65 Nepse Days before end date", http.StatusBadRequest)
		return
	}

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

			result := preprocessChange(re, start, end)

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
			"A1": "Company Name", "B1": "Ticker",
			"C1": "Start", "D1": "End", "E1": "Point Change", "F1": "Percentage Change",
			"G1": "Maximum Price", "H1": "Minimum Price", "I1": "Max-Min", "J1": "%Change(Extreme)",
			"K1": "Maximum Average", "L1": "Minimum Average", "M1": "Max-Min(avg)", "N1": "%Change(Avg)",
		}

		var excelVals []map[string]interface{}

		for k, v := range resu {
			excelVal := map[string]interface{}{
				utils.GetColumn("A", k): v.CompanyName, utils.GetColumn("B", k): v.Ticker,
				utils.GetColumn("C", k): v.StartPrice, utils.GetColumn("D", k): v.EndPrice, utils.GetColumn("E", k): v.PointChange, utils.GetColumn("F", k): v.PercentageChange,
				utils.GetColumn("G", k): v.MaxPrice, utils.GetColumn("H", k): v.MinPrice, utils.GetColumn("I", k): v.ExtremePointChange, utils.GetColumn("J", k): v.ExtremePercentageChange,
				utils.GetColumn("K", k): v.MaxAveragePrice, utils.GetColumn("L", k): v.MinAveragePrice, utils.GetColumn("M", k): v.AveragePointChange, utils.GetColumn("N", k): v.AveragePercentageChange,
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

		f := excelize.NewFile()
		for k, v := range categories {
			f.SetCellValue("Sheet1", k, v)
		}

		for _, vals := range excelVals {
			for k, v := range vals {
				f.SetCellValue("Sheet1", k, v)
			}
		}

		if err := f.SaveAs(fmt.Sprintf("%s/%s.xlsx", folderName, sec)); err != nil {
			fmt.Println(err)
		}

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

func preprocessChange(minHistory []nepse.PriceHistoryMinified, start, end string) Result {
	var closeStart, closeEnd float64
	var max = minHistory[0].HighPrice
	var minP = minHistory[0].LowPrice
	var maxAvg = minHistory[0].AveragePrice
	var minAvg = minHistory[0].AveragePrice

	for _, min := range minHistory {

		if min.Date == start {
			closeStart = min.Price
		}

		if min.Date == end {
			closeEnd = min.Price
		}

		if min.HighPrice > max {
			max = min.HighPrice
		}

		if min.AveragePrice > maxAvg {
			maxAvg = min.AveragePrice
		}

		if min.AveragePrice < minAvg {
			minAvg = min.AveragePrice
		}

		if min.LowPrice < minP {
			minP = min.LowPrice
		}

	}

	changeClosed := closeEnd - closeStart

	changeExtreme := max - minP

	changeAverage := maxAvg - minAvg

	percentageChangeOnExtreme := (changeExtreme / minP) * 100

	percentageChangeOnAverage := (changeAverage / minAvg) * 100

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

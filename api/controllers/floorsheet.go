package controllers

import (
	"fmt"
	responses "nepse-backend/api/response"
	"nepse-backend/nepse"
	"nepse-backend/nepse/neweb"
	"nepse-backend/utils"
	"net/http"
	"os"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type FloorsheetResult struct {
	Ticker                string             `json:"ticker"`
	BuyerQuantityMap      map[string]int64   `json:"buyerQuantityMap"`
	BuyerTurnOverMap      map[string]float64 `json:"buyerTurnOverMap"`
	BuyerAveragePriceMap  map[string]float64 `json:"buyerAveragePriceMap"`
	SellerQuantityMap     map[string]int64   `json:"sellerQuantityMap"`
	SellerTurnOverMap     map[string]float64 `json:"sellerTurnOverMap"`
	SellerAveragePriceMap map[string]float64 `json:"sellerAveragePriceMap"`
}

func (s *Server) GetFloorsheet(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	ticker := params.Get("ticker")
	start := params.Get("start")
	end := params.Get("end")
	randomId := params.Get("id")

	days, err := utils.GetDateRange(w, start, end)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	fmt.Println("days", days)

	var nepseBeta nepse.NepseInterface

	nepseBeta, err = neweb.Neweb()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nepseSectors, err := nepseBeta.GetStocks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var id string

	for _, stock := range nepseSectors {
		if ticker == stock.Ticker {
			id = stock.Id
		}
	}

	var result FloorsheetResult

	result.BuyerQuantityMap = make(map[string]int64)
	result.SellerQuantityMap = make(map[string]int64)
	result.BuyerTurnOverMap = make(map[string]float64)
	result.SellerTurnOverMap = make(map[string]float64)
	result.BuyerAveragePriceMap = make(map[string]float64)
	result.SellerAveragePriceMap = make(map[string]float64)

	var floorsheetContents []nepse.FloorsheetContent

	for _, day := range days {

		var count = 0
		for {
			time.Sleep(400 * time.Millisecond)
			floorsheet, err := nepseBeta.GetFloorsheet(id, day, randomId, count, 2000)

			if err != nil {
				responses.ERROR(w, 400, err)
				return
			}

			floorsheetContents = append(floorsheetContents, floorsheet.Floorsheets.Content...)

			isLastPage := floorsheet.Floorsheets.Last
			count++
			if isLastPage {
				break
			}

		}
	}

	for _, sheetData := range floorsheetContents {
		if sheetData.Buyermemberid != "" {
			result.BuyerQuantityMap[sheetData.Buyermemberid] += int64(sheetData.Contractquantity)
		}
		if sheetData.Sellermemberid != "" {
			result.SellerQuantityMap[sheetData.Sellermemberid] += int64(sheetData.Contractquantity)
		}
		if sheetData.Buyermemberid != "" {
			result.BuyerTurnOverMap[sheetData.Buyermemberid] += sheetData.Contractamount
		}
		if sheetData.Sellermemberid != "" {
			result.SellerTurnOverMap[sheetData.Sellermemberid] += sheetData.Contractamount
		}
	}

	for k, v := range result.BuyerQuantityMap {
		result.BuyerAveragePriceMap[k] = utils.ToFixed(float64(result.BuyerTurnOverMap[k])/float64(v), 2)
	}
	for k, v := range result.SellerQuantityMap {
		result.SellerAveragePriceMap[k] = utils.ToFixed(float64(result.SellerTurnOverMap[k])/float64(v), 2)
	}
	result.Ticker = ticker

	var allCharts []*charts.Bar
	folderName := "floorsheet-today"
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		os.Mkdir(folderName, 0777)
	}

	allCharts = append(allCharts, BarGraphFS(result.BuyerQuantityMap, result.SellerQuantityMap, "Top Buyers"))
	allCharts = append(allCharts, BarGraphFS(result.SellerQuantityMap, result.BuyerQuantityMap, "Top Sellers"))

	go CreateHTMLFS(allCharts, fmt.Sprintf("%s/%s", folderName, ticker))

	responses.JSON(w, http.StatusOK, floorsheetContents)
}

func BarGraphFS(aggregatedData, alterAggregateData map[string]int64, title string) *charts.Bar {
	bar := charts.NewBar()

	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: title,
	}))

	topSorted := SortMap(aggregatedData, false)

	var keys = make([]string, 0)

	for _, v := range topSorted {
		keys = append(keys, v.Key)
	}

	var alterSorted []kv
	for _, v := range topSorted {
		alterSorted = append(alterSorted, kv{v.Key, alterAggregateData[v.Key]})
	}

	// Put data into instance
	bar.SetXAxis(keys).
		AddSeries("Category A", generateBarItems(topSorted)).
		AddSeries("Category B", generateBarItems(alterSorted))

	return bar
}

func CreateHTMLFS(barCharts []*charts.Bar, fileName string) {
	page := components.NewPage()

	for _, v := range barCharts {
		page.AddCharts(v)
	}
	f, _ := os.Create(fmt.Sprintf("%s.html", fileName))
	page.Render(f)
}

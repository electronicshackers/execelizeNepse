package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	responses "nepse-backend/api/response"
	"nepse-backend/nepse"
	"nepse-backend/nepse/neweb"
	"nepse-backend/utils"
	"net/http"
	"os"
	"sort"
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

type TransactionData struct {
	Ticker                string
	Quantity              int64
	TotalQuantity         int64
	PercentageShare       float64
	BrokerPercentageShare float64
}

func (s *Server) FloorsheetAnalysis(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	start := params.Get("start")
	end := params.Get("end")
	randomId := params.Get("id")

	days, err := utils.GetDateRange(w, start, end)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

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

	var floorsheetContents []nepse.FloorsheetContent
	for _, day := range days {
		for _, v := range nepseSectors {
			var count = 0
			for {
				time.Sleep(800 * time.Millisecond)
				floorsheet, err := nepseBeta.GetFloorsheet(v.Id, day, randomId, count, 2000)

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
	}

	aggregatedDataBuy := make(map[string][]TransactionData)
	aggregatedDataSell := make(map[string][]TransactionData)

	brokerBuyVolume := make(map[string]int64)
	brokerSellVolume := make(map[string]int64)

	stockMap := make(map[string]int64)

	for _, v := range floorsheetContents {
		stockMap[v.Stocksymbol] += int64(v.Contractquantity)
		if v.Buyermemberid != "" {
			brokerBuyVolume[v.Buyermemberid] += int64(v.Contractquantity)
			aggregatedDataBuy[v.Buyermemberid] = append(aggregatedDataBuy[v.Buyermemberid], TransactionData{Ticker: v.Stocksymbol, Quantity: int64(v.Contractquantity)})
		}

		if v.Sellermemberid != "" {
			brokerSellVolume[v.Sellermemberid] += int64(v.Contractquantity)
			aggregatedDataSell[v.Sellermemberid] = append(aggregatedDataSell[v.Sellermemberid], TransactionData{Ticker: v.Stocksymbol, Quantity: int64(v.Contractquantity)})
		}
	}

	finalDataBuy := make(map[string][]TransactionData)
	finalDataSell := make(map[string][]TransactionData)

	for k, v := range aggregatedDataBuy {
		sumDataBuy := make(map[string]int64)
		for _, w := range v {
			sumDataBuy[w.Ticker] += w.Quantity
		}
		totalBrokerVolume := brokerBuyVolume[k]

		for key, value := range sumDataBuy {
			total := stockMap[key]
			brokerPercentage := float64(value) / float64(totalBrokerVolume) * 100
			percentage := float64(value) / float64(total) * 100
			finalDataBuy[k] = append(finalDataBuy[k], TransactionData{Ticker: key, Quantity: value, TotalQuantity: total, PercentageShare: percentage, BrokerPercentageShare: brokerPercentage})
		}
	}

	for k, v := range aggregatedDataSell {
		sumDataSell := make(map[string]int64)
		for _, w := range v {
			sumDataSell[w.Ticker] += w.Quantity
		}
		totalBrokerVolume := brokerSellVolume[k]

		for key, value := range sumDataSell {
			total := stockMap[key]
			percentage := float64(value) / float64(total) * 100
			brokerPercentage := float64(value) / float64(totalBrokerVolume) * 100
			finalDataSell[k] = append(finalDataSell[k], TransactionData{Ticker: key, Quantity: value, TotalQuantity: total, PercentageShare: percentage, BrokerPercentageShare: brokerPercentage})
		}
	}

	var pumpedStockDataBroker = make(map[string][]TransactionData)
	for k, v := range finalDataBuy {

		for _, stockTransaction := range v {
			if stockTransaction.PercentageShare > 25 && stockTransaction.BrokerPercentageShare > 5 {
				pumpedStockDataBroker[k] = append(pumpedStockDataBroker[k], stockTransaction)
			}
		}
	}

	var dumpedStockDataBroker = make(map[string][]TransactionData)

	for k, v := range finalDataSell {
		for _, stockTransaction := range v {
			if stockTransaction.PercentageShare > 25 && stockTransaction.BrokerPercentageShare > 5 {
				dumpedStockDataBroker[k] = append(dumpedStockDataBroker[k], stockTransaction)
			}
		}
	}

	var brokerBuySorted []kv
	for k, v := range brokerBuyVolume {
		brokerBuySorted = append(brokerBuySorted, kv{k, v})
	}

	sort.Slice(brokerBuySorted, func(i, j int) bool {
		return brokerBuySorted[i].Value > brokerBuySorted[j].Value
	})

	for k, v := range finalDataBuy {
		sort.Slice(v, func(i, j int) bool {
			return v[i].Quantity > v[j].Quantity
		})
		if len(v) > 10 {
			finalDataBuy[k] = v[0:10]
		}
	}

	var brokerSellSorted []kv
	for k, v := range brokerSellVolume {
		brokerSellSorted = append(brokerSellSorted, kv{k, v})
	}
	sort.Slice(brokerSellSorted, func(i, j int) bool {
		return brokerSellSorted[i].Value > brokerSellSorted[j].Value
	})

	for k, v := range finalDataSell {
		sort.Slice(v, func(i, j int) bool {
			return v[i].Quantity > v[j].Quantity
		})
		if len(v) > 10 {
			finalDataSell[k] = v[0:10]
		}
	}

	var buyCharts []*charts.Bar

	var topBrokerSellData = make(map[string][]TransactionData)
	for k, v := range finalDataSell {
		for _, bs := range brokerSellSorted[0:10] {
			if k == bs.Key {
				topBrokerSellData[k] = v
				buyCharts = append(buyCharts, BarGraphAgg(finalDataSell[k], fmt.Sprintf("Top Sell of Broker Number %s", k)))
			}
		}
	}

	var topBrokerBuyData = make(map[string][]TransactionData)
	for k, v := range finalDataBuy {
		for _, bs := range brokerBuySorted[0:10] {
			if k == bs.Key {
				topBrokerBuyData[k] = v
				buyCharts = append(buyCharts, BarGraphAgg(finalDataBuy[k], fmt.Sprintf("Top Buy of Broker Number %s", k)))
			}
		}
	}

	for k, v := range pumpedStockDataBroker {
		sort.Slice(v, func(i, j int) bool {
			return v[i].Quantity > v[j].Quantity
		})
		buyCharts = append(buyCharts, BarGraphAgg(v, fmt.Sprintf("Top Pumped Stock of Broker Number %s", k)))
	}

	for k, v := range dumpedStockDataBroker {
		sort.Slice(v, func(i, j int) bool {
			return v[i].Quantity > v[j].Quantity
		})
		buyCharts = append(buyCharts, BarGraphAgg(v, fmt.Sprintf("Top Dumped Stock of Broker Number %s", k)))
	}

	CreateHTMLAgg(buyCharts, "analysis")

}

func (s *Server) GetFloorSheetAggregated(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	start := params.Get("start")
	end := params.Get("end")
	randomId := params.Get("id")

	days, err := utils.GetDateRange(w, start, end)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

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

	var floorsheetContents []nepse.FloorsheetContent
	for _, day := range days {
		for _, v := range nepseSectors {
			var count = 0
			for {
				time.Sleep(800 * time.Millisecond)
				floorsheet, err := nepseBeta.GetFloorsheet(v.Id, day, randomId, count, 2000)

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
	}
	document, _ := json.MarshalIndent(floorsheetContents, "", " ")
	_ = ioutil.WriteFile("./floorsheet.json", document, 0644)
	aggregatedDataBuy := make(map[string][]TransactionData)
	aggregatedDataSell := make(map[string][]TransactionData)

	for _, v := range floorsheetContents {
		if v.Buyermemberid != "" {
			aggregatedDataBuy[v.Buyermemberid] = append(aggregatedDataBuy[v.Buyermemberid], TransactionData{Ticker: v.Stocksymbol, Quantity: int64(v.Contractquantity)})
		}

		if v.Sellermemberid != "" {
			aggregatedDataSell[v.Sellermemberid] = append(aggregatedDataSell[v.Sellermemberid], TransactionData{Ticker: v.Stocksymbol, Quantity: int64(v.Contractquantity)})
		}

	}

	finalDataBuy := make(map[string][]TransactionData)
	finalDataSell := make(map[string][]TransactionData)

	for k, v := range aggregatedDataBuy {
		sumDataBuy := make(map[string]int64)
		for _, w := range v {
			sumDataBuy[w.Ticker] += w.Quantity
		}

		for key, value := range sumDataBuy {
			finalDataBuy[k] = append(finalDataBuy[k], TransactionData{Ticker: key, Quantity: value})
		}
	}

	for k, v := range aggregatedDataSell {
		sumDataSell := make(map[string]int64)
		for _, w := range v {
			sumDataSell[w.Ticker] += w.Quantity
		}

		for key, value := range sumDataSell {
			finalDataSell[k] = append(finalDataSell[k], TransactionData{Ticker: key, Quantity: value})
		}
	}

	var buyCharts []*charts.Bar
	folderName := "floorsheetAgg"
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		os.Mkdir(folderName, 0777)
	}
	for k, v := range finalDataBuy {
		sort.Slice(v, func(i, j int) bool {
			return v[i].Quantity > v[j].Quantity
		})
		finalDataBuy[k] = v[0:10]
		buyCharts = append(buyCharts, BarGraphAgg(finalDataBuy[k], fmt.Sprintf("Top Buy of Broker Number %s", k)))
	}
	CreateHTMLAgg(buyCharts, "weekBuy")

	var sellCharts []*charts.Bar
	for k, v := range finalDataSell {
		sort.Slice(v, func(i, j int) bool {
			return v[i].Quantity > v[j].Quantity
		})
		finalDataSell[k] = v[0:10]
		sellCharts = append(sellCharts, BarGraphAgg(finalDataSell[k], fmt.Sprintf("Top Sell of Broker Number %s", k)))
	}
	CreateHTMLAgg(sellCharts, "weekSell")
	responses.JSON(w, 200, finalDataBuy)
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
			time.Sleep(1000 * time.Millisecond)
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

	document, _ := json.MarshalIndent(floorsheetContents, "", " ")
	_ = ioutil.WriteFile(fmt.Sprintf("./%s.json", ticker), document, 0644)

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
	folderName := "for-youtube"
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		os.Mkdir(folderName, 0777)
	}

	allCharts = append(allCharts, BarGraphFS(result.BuyerQuantityMap, result.SellerQuantityMap, "Top Buyers"))
	allCharts = append(allCharts, BarGraphFS(result.SellerQuantityMap, result.BuyerQuantityMap, "Top Sellers"))

	go CreateHTMLFS(allCharts, fmt.Sprintf("%s/%s", folderName, ticker))

	responses.JSON(w, http.StatusOK, result)
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

func BarGraphAgg(data []TransactionData, title string) *charts.Bar {
	bar := charts.NewBar()

	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: title,
	}))

	var keys = make([]string, 0)

	for _, v := range data {
		keys = append(keys, v.Ticker)
	}

	// Put data into instance
	bar.SetXAxis(keys).
		AddSeries("Category A", generateBarItemsAgg(data))

	return bar
}

func generateBarItemsAgg(data []TransactionData) []opts.BarData {
	items := make([]opts.BarData, 0)

	for _, v := range data {
		items = append(items, opts.BarData{Value: v.Quantity})
	}
	return items
}

func CreateHTMLAgg(barCharts []*charts.Bar, fileName string) {
	page := components.NewPage()

	for _, v := range barCharts {
		page.AddCharts(v)
	}
	f, _ := os.Create(fmt.Sprintf("%s.html", fileName))
	page.Render(f)
}

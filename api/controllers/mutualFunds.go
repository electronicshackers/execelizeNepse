package controllers

import (
	responses "nepse-backend/api/response"
	"nepse-backend/nepse/bizmandu"
	"nepse-backend/nepse/neweb"
	"nepse-backend/utils"
	"net/http"
	"os"
	"sort"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var (
	options     = []string{"whole", "sector", "topHolding", "topSold", "topBought", "netBought"}
	MutualFunds = []string{"CMF1", "CMF2", "GIMES1", "KEF", "LEMF", "LUK", "NBF2", "NEF", "NIBLPF", "NIBSF1",
		"NICBF", "NICGF", "NMB50", "NMBHF1", "PSF", "SAEF", "SBCF", "SEF", "SFMF", "SIGS2", "SLCF"}
)

type kv struct {
	Key   string
	Value int64
}

type MutualFund struct {
	MutualFundKeyMetrics []MutualFundKeyMetrics `json:"mutualFundKeyMetrics"`
	SectorMap            map[string]float64     `json:"sectorMap"`
	TopHoldingMap        map[string]int64       `json:"topHoldingMap"`
	TopstockboughtMap    map[string]int64       `json:"topStockBoughtMap"`
	TopstocksoldMap      map[string]int64       `json:"topStockSoldMap"`
	NetstockboughtMap    map[string]int64       `json:"netStockBoughtMap"`
}

type MutualFundKeyMetrics struct {
	Ticker               string                      `json:"ticker"`
	WeeklyNav            float64                     `json:"weeklyNav"`
	MonthlyNav           float64                     `json:"monthlyNav"`
	PriceVsNav           float64                     `json:"priceVsNav"`
	MarketCapatilization float64                     `json:"marketCapatilization"`
	TotalSector          int                         `json:"totalSector"`
	TotalCompanies       int                         `json:"totalCompanies"`
	LastTradedPrice      float64                     `json:"lastTradedPrice"`
	Sector               []bizmandu.Sector           `json:"sector"`
	Topstockholdings     []bizmandu.Topstockholdings `json:"topStockHoldings"`
	Topstockbought       []bizmandu.Topstock         `json:"topStockBought"`
	Topstocksold         []bizmandu.Topstock         `json:"topStockSold"`
}

func (server *Server) GetMutualFundsInfo(w http.ResponseWriter, r *http.Request) {
	biz, err := bizmandu.NewBizmandu()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nepseBeta, err := neweb.Neweb()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mfs, err := nepseBeta.GetMutualFundStock()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var mfsInfo MutualFund

	for _, mf := range mfs {
		mutualFundData, err := biz.GetMutualFundData(mf.Ticker)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		mfsInfo.MutualFundKeyMetrics = append(mfsInfo.MutualFundKeyMetrics, MutualFundKeyMetrics{
			Ticker:               mf.Ticker,
			WeeklyNav:            mutualFundData.Message.Summary.Weeklynav,
			MonthlyNav:           mutualFundData.Message.Summary.Monthlynav,
			LastTradedPrice:      mf.Lasttradedprice,
			PriceVsNav:           mutualFundData.Message.Summary.Pricevsnav,
			MarketCapatilization: mutualFundData.Message.Summary.Aum,
			TotalSector:          int(mutualFundData.Message.Summary.Totalsectorsinvested),
			TotalCompanies:       int(mutualFundData.Message.Summary.Totalcompaniesheld),
			Sector:               mutualFundData.Message.Summary.Sector,
			Topstockholdings:     mutualFundData.Message.Summary.Topstockholdings,
			Topstockbought:       mutualFundData.Message.Summary.Topstockbought,
			Topstocksold:         mutualFundData.Message.Summary.Topstocksold,
		})
	}

	mfsInfo.SectorMap = make(map[string]float64)
	mfsInfo.TopHoldingMap = make(map[string]int64)
	mfsInfo.TopstockboughtMap = make(map[string]int64)
	mfsInfo.TopstocksoldMap = make(map[string]int64)
	mfsInfo.NetstockboughtMap = make(map[string]int64)

	for _, mf := range mfsInfo.MutualFundKeyMetrics {
		for _, sector := range mf.Sector {
			mfsInfo.SectorMap[sector.Label] += sector.Value
		}

		for _, topHolding := range mf.Topstockholdings {
			mfsInfo.TopHoldingMap[topHolding.Ticker] += topHolding.Qty
		}
		for _, topBought := range mf.Topstockbought {
			mfsInfo.TopstockboughtMap[topBought.Ticker] += topBought.Noofstocks
		}
		for _, topSold := range mf.Topstocksold {
			mfsInfo.TopstocksoldMap[topSold.Ticker] += topSold.Noofstocks
		}
	}

	for label, value := range mfsInfo.SectorMap {
		mfsInfo.SectorMap[label] = utils.ToFixed((value / float64(len(mfsInfo.MutualFundKeyMetrics)) * 100), 2)
	}

	for k, _ := range mfsInfo.TopstockboughtMap {
		mfsInfo.NetstockboughtMap[k] = mfsInfo.TopstockboughtMap[k] - mfsInfo.TopstocksoldMap[k]
	}

	folderName := "mutualFund"
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		os.Mkdir(folderName, 0777)
	}

	categories := GetMutualFundHeaders()

	var excelVals []map[string]interface{}

	for k, v := range mfsInfo.MutualFundKeyMetrics {
		excelVal := GetMutualFundValues(v, k)
		if v.Ticker != "" {
			excelVals = append(excelVals, excelVal)
		}
	}

	go utils.CreateExcelFile(folderName, "whole", categories, excelVals)

	var allCharts []*charts.Bar
	var pie *charts.Pie

	for _, option := range options {
		aggregatedHeaders := GetAggregatedMutualFundHeaders(option)

		var excelVals []map[string]interface{}

		if option == "sector" {
			pie = PieChart(mfsInfo.SectorMap, "Sector Distribution of Mutual Funds")
			var count = 0
			for k, v := range mfsInfo.SectorMap {
				excelVal := GetAggregatedMutualFundValues(k, v, count)
				excelVals = append(excelVals, excelVal)
				count++
			}
		}

		if option == "topHolding" {
			allCharts = append(allCharts, BarGraph(mfsInfo.TopHoldingMap, "Total Number of Holdings in Mutual Fund's Portfolio"))
			var count = 0
			for k, v := range mfsInfo.TopHoldingMap {
				excelVal := GetAggregatedMutualFundValues(k, v, count)
				excelVals = append(excelVals, excelVal)
				count++
			}
		}
		if option == "topBought" {
			allCharts = append(allCharts, BarGraph(mfsInfo.TopstockboughtMap, "Top Stock Bought within One Month"))
			var count = 0
			for k, v := range mfsInfo.TopstockboughtMap {
				excelVal := GetAggregatedMutualFundValues(k, v, count)
				excelVals = append(excelVals, excelVal)
				count++
			}
		}
		if option == "topSold" {
			allCharts = append(allCharts, BarGraph(mfsInfo.TopstocksoldMap, "Top Stock Sold within One Month"))
			var count = 0
			for k, v := range mfsInfo.TopstocksoldMap {
				excelVal := GetAggregatedMutualFundValues(k, v, count)
				excelVals = append(excelVals, excelVal)
				count++
			}
		}

		if option == "netBought" {
			allCharts = append(allCharts, BarGraph(mfsInfo.NetstockboughtMap, "Net Stock Bought Of Mutual Fund For One Month"))
			var count = 0
			for k, v := range mfsInfo.NetstockboughtMap {
				excelVal := GetAggregatedMutualFundValues(k, v, count)
				excelVals = append(excelVals, excelVal)
				count++
			}
		}

		go utils.CreateExcelFile(folderName, option, aggregatedHeaders, excelVals)
	}

	go CreateHTML(allCharts, pie)

	responses.JSON(w, http.StatusOK, mfsInfo)
}

func SortMap(m map[string]int64) []kv {

	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	return ss[0:10]
}

func generateBarItems(data []kv) []opts.BarData {
	items := make([]opts.BarData, 0)

	for _, v := range data {
		items = append(items, opts.BarData{Value: v.Value})
	}
	return items
}

func generatePieItems(sector map[string]float64) []opts.PieData {
	items := make([]opts.PieData, 0)
	for k, v := range sector {
		items = append(items, opts.PieData{Value: v, Name: k})
	}
	return items
}

func PieChart(aggregateDate map[string]float64, title string) *charts.Pie {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: title}),
	)

	pie.AddSeries("pie", generatePieItems(aggregateDate)).
		SetSeriesOptions(charts.WithLabelOpts(
			opts.Label{
				Show:      true,
				Formatter: "{b}: {c}",
			}),
		)
	return pie
}

func BarGraph(aggregatedData map[string]int64, title string) *charts.Bar {
	bar := charts.NewBar()

	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: title,
	}))

	topSorted := SortMap(aggregatedData)

	var keys = make([]string, 0)

	for _, v := range topSorted {
		keys = append(keys, v.Key)
	}

	// Put data into instance
	bar.SetXAxis(keys).
		AddSeries("Category A", generateBarItems(topSorted))

	return bar
}

func CreateHTML(barCharts []*charts.Bar, pieChart *charts.Pie) {
	page := components.NewPage()

	page.AddCharts(pieChart)

	for _, v := range barCharts {
		page.AddCharts(v)
	}
	f, _ := os.Create("index.html")
	page.Render(f)
}

func GetAggregatedMutualFundValues(key string, value interface{}, k int) map[string]interface{} {
	excelVals := map[string]interface{}{
		utils.GetColumn("A", k): key, utils.GetColumn("B", k): value,
	}

	return excelVals
}

func GetAggregatedMutualFundHeaders(option string) map[string]string {
	var headers = make(map[string]string)

	if option == "sector" {
		headers = map[string]string{
			"A1": "Sector",
			"B1": "Distribution",
		}
	} else {
		headers = map[string]string{
			"A1": "Ticker",
			"B1": "TotalStock",
		}
	}

	return headers

}

func GetMutualFundHeaders() map[string]string {
	headers := map[string]string{
		"A1": "Ticker", "B1": "LastTradedPrice", "C1": "WeeklyNav",
		"D1": "MonthlyNav", "E1": "PriceVsNav", "F1": "AUM",
		"G1": "TotalSector", "H1": "TotalCompanies",
	}
	return headers
}

func GetMutualFundValues(data MutualFundKeyMetrics, k int) map[string]interface{} {
	excelVals := map[string]interface{}{
		utils.GetColumn("A", k): data.Ticker, utils.GetColumn("B", k): data.LastTradedPrice, utils.GetColumn("C", k): data.WeeklyNav,
		utils.GetColumn("D", k): data.MonthlyNav, utils.GetColumn("E", k): (data.LastTradedPrice - data.WeeklyNav), utils.GetColumn("F", k): data.MarketCapatilization,
		utils.GetColumn("G", k): data.TotalSector, utils.GetColumn("H", k): data.TotalCompanies,
	}

	return excelVals
}

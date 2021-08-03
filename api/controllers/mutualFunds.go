package controllers

import (
	responses "nepse-backend/api/response"
	"nepse-backend/nepse/bizmandu"
	"nepse-backend/nepse/neweb"
	"nepse-backend/utils"
	"net/http"
	"os"
)

var options = []string{"whole", "sector", "topHolding", "topSold", "topBought"}

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

type AggregatedMutualFundMetrics struct {
	SectorMap         map[string]float64 `json:"sectorMap"`
	TopHoldingMap     map[string]int64   `json:"topHoldingMap"`
	TopstockboughtMap map[string]int64   `json:"topStockBoughtMap"`
	TopstocksoldMap   map[string]int64   `json:"topStockSoldMap"`
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

	var mfsInfo []MutualFundKeyMetrics

	for _, mf := range mfs {
		mutualFundData, err := biz.GetMutualFundData(mf.Ticker)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		mfsInfo = append(mfsInfo, MutualFundKeyMetrics{
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

	var aggregatedMetrics AggregatedMutualFundMetrics

	var sectorMap = make(map[string]float64)
	var topHoldingMap = make(map[string]int64)
	var topBoughtMap = make(map[string]int64)
	var topSoldMap = make(map[string]int64)

	for _, mf := range mfsInfo {
		for _, sector := range mf.Sector {
			sectorMap[sector.Label] += sector.Value
		}

		for _, topHolding := range mf.Topstockholdings {
			topHoldingMap[topHolding.Ticker] += topHolding.Qty
		}
		for _, topBought := range mf.Topstockbought {
			topBoughtMap[topBought.Ticker] += topBought.Noofstocks
		}
		for _, topSold := range mf.Topstocksold {
			topSoldMap[topSold.Ticker] += topSold.Noofstocks
		}
	}

	for label, value := range sectorMap {
		sectorMap[label] = utils.ToFixed(value/float64(len(mfsInfo)), 2) * 100
	}

	aggregatedMetrics.SectorMap = sectorMap
	aggregatedMetrics.TopHoldingMap = topHoldingMap
	aggregatedMetrics.TopstockboughtMap = topBoughtMap
	aggregatedMetrics.TopstocksoldMap = topSoldMap

	folderName := "mutualFund"
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		os.Mkdir(folderName, 0777)
	}

	for _, option := range options {
		categories := GetMutualFundHeaders(option)

		var excelVals []map[string]interface{}

		for k, v := range mfsInfo {
			excelVal := GetMutualFundValues(option, v, k)
			if v.Ticker != "" {
				excelVals = append(excelVals, excelVal)
			}
		}

		go utils.CreateExcelFile(folderName, option, categories, excelVals)
	}

	responses.JSON(w, http.StatusOK, mfsInfo)
}

func GetMutualFundHeaders(option string) map[string]string {
	var headers = make(map[string]string)
	if option == "whole" {
		headers = map[string]string{
			"A1": "Ticker", "B1": "LastTradedPrice", "C1": "WeeklyNav",
			"D1": "MonthlyNav", "E1": "PriceVsNav", "F1": "AUM",
			"G1": "TotalSector", "H1": "TotalCompanies",
		}
	}
	if option == "sector" {
		headers = map[string]string{
			"A1": "Sector", "B1": "Percentage",
		}
	}
	if option == "topstock" {
		headers = map[string]string{
			"A1": "Ticker", "B1": "Qty",
		}
	}
	return headers
}

func GetMutualFundValues(option string, data MutualFundKeyMetrics, k int) map[string]interface{} {
	var excelVals = make(map[string]interface{})
	if option == "whole" {
		excelVals = map[string]interface{}{
			utils.GetColumn("A", k): data.Ticker, utils.GetColumn("B", k): data.LastTradedPrice, utils.GetColumn("C", k): data.WeeklyNav,
			utils.GetColumn("D", k): data.MonthlyNav, utils.GetColumn("E", k): (data.LastTradedPrice - data.WeeklyNav), utils.GetColumn("F", k): data.MarketCapatilization,
			utils.GetColumn("G", k): data.TotalSector, utils.GetColumn("H", k): data.TotalCompanies,
		}
	}
	if option == "sector" {
		excelVals = map[string]interface{}{
			utils.GetColumn("A", k): data.Sector, utils.GetColumn("B", k): data.MarketCapatilization,
		}
	}
	return excelVals
}

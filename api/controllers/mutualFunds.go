package controllers

import (
	responses "nepse-backend/api/response"
	"nepse-backend/nepse/bizmandu"
	"nepse-backend/nepse/neweb"
	"nepse-backend/utils"
	"net/http"
	"os"
)

type MutualFundKeyMetrics struct {
	Ticker               string  `json:"ticker"`
	WeeklyNav            float64 `json:"weeklyNav"`
	MonthlyNav           float64 `json:"monthlyNav"`
	PriceVsNav           float64 `json:"priceVsNav"`
	MarketCapatilization float64 `json:"marketCapatilization"`
	TotalSector          int     `json:"totalSector"`
	TotalCompanies       int     `json:"totalCompanies"`
	LastTradedPrice      float64 `json:"lastTradedPrice"`
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
		})
	}

	folderName := "mutualFund"
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		os.Mkdir(folderName, 0777)
	}
	responses.JSON(w, http.StatusOK, mfsInfo)

	categories := GetMutualFundHeaders()

	var excelVals []map[string]interface{}

	for k, v := range mfsInfo {
		excelVal := GetMutualFundValues(v, k)
		if v.Ticker != "" {
			excelVals = append(excelVals, excelVal)
		}
	}

	go utils.CreateExcelFile(folderName, "whole", categories, excelVals)
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
	excelVal := map[string]interface{}{
		utils.GetColumn("A", k): data.Ticker, utils.GetColumn("B", k): data.LastTradedPrice, utils.GetColumn("C", k): data.WeeklyNav,
		utils.GetColumn("D", k): data.MonthlyNav, utils.GetColumn("E", k): (data.LastTradedPrice - data.WeeklyNav), utils.GetColumn("F", k): data.MarketCapatilization,
		utils.GetColumn("G", k): data.TotalSector, utils.GetColumn("H", k): data.TotalCompanies,
	}
	return excelVal
}

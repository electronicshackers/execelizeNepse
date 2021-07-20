package controllers

import (
	"fmt"
	"nepse-backend/nepse/bizmandu"
	"nepse-backend/utils"
	"net/http"
	"os"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type KeyFinancialMetrics struct {
	Ticker                     string  `json:"ticker"`
	LTP                        float64 `json:"ltp"`
	DiversionFromFair          float64 `json:"divesionFromFair"`
	PE                         float64 `json:"pe"`
	Eps                        float64 `json:"eps"`
	FairValue                  float64 `json:"fairValue"`
	Bvps                       float64 `json:"bvps"`
	Pbv                        float64 `json:"pbv"`
	Roa                        float64 `json:"roa"`
	Roe                        float64 `json:"roe"`
	NPL                        float64 `json:"npl"`
	Listedshares               float64 `json:"listedShares"`
	Reserves                   float64 `json:"reserves"`
	Mktcap                     float64 `json:"mktCap"`
	DistributableProfit        float64 `json:"distributableProfit"`
	DistibutableProfitPerShare float64 `json:"distributableProfitPerShare"`
	PaidUpCapital              float64 `json:"paidUpCapital"`
	DividendCapacity           float64 `json:"dividendCapacity"`
	RetentionRatio             float64 `json:"retentionRatio"`
}

func (server *Server) GetFundamentalSectorwise(w http.ResponseWriter, r *http.Request) {
	sector := r.URL.Query().Get("sector")
	if sector == "" {
		http.Error(w, "sector is required", http.StatusBadRequest)
		return
	}

	querySector := strings.Split(sector, ",")

	sectors := utils.MapColumns(querySector)

	biz, err := bizmandu.NewBizmandu()

	if err != nil {
		fmt.Println("err", err)
	}

	if err != nil {
		fmt.Println("err", err)
	}

	for _, sector := range sectors {

		sectorStocks, err := biz.GetSectorStock(sector)
		fmt.Println("sectorStocks", sectorStocks)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var keys []KeyFinancialMetrics

		for _, ticker := range sectorStocks {
			var key KeyFinancialMetrics

			detail, err := biz.GetSummary(ticker.Ticker)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			financial, err := biz.GetFinancial(ticker.Ticker)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			balancesheet, err := biz.GetBalanceSheet(ticker.Ticker)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			incomeStatement, err := biz.GetIncomeStatement(ticker.Ticker)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			price, err := biz.GetCurrentPrice(ticker.Ticker)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			key.Eps = utils.ToFixed(detail.Message.Summary.Epsdiluted, 2)
			key.PE = utils.ToFixed(detail.Message.Summary.Pediluted, 2)
			key.LTP = detail.Message.Summary.Open
			key.Bvps = utils.ToFixed(detail.Message.Summary.Bvps, 2)
			key.Ticker = detail.Message.Keyfinancial.Ticker
			key.Listedshares = detail.Message.Summary.Listedshares
			key.LTP = price.Lasttradedprice
			key.Mktcap = detail.Message.Summary.Mktcap
			key.Pbv = utils.ToFixed(key.LTP/key.Bvps, 2)

			for _, quarter := range detail.Message.Keyfinancial.Data {
				if quarter.Type == "CURRENT" {
					key.Roa = utils.ToFixed(quarter.Roa*100, 2)
					key.Roe = utils.ToFixed(quarter.Roe*100, 2)
				}
			}

			if len(financial.Message.Data) != 0 {
				key.NPL = utils.ToFixed(financial.Message.Data[0].Nonperformingloannpltototalloan*100, 2)
			}

			if len(balancesheet.Message.Data) != 0 {
				key.PaidUpCapital = float64(balancesheet.Message.Data[0].Paidupcapital)
				key.Reserves = float64(balancesheet.Message.Data[0].Reserves)
			}

			if len(incomeStatement.Message.Data) != 0 {
				key.DistributableProfit = float64(incomeStatement.Message.Data[0].Freeprofit)
				key.DistibutableProfitPerShare = utils.ToFixed((key.DistributableProfit/key.Listedshares)*100, 2)
			}

			if len(balancesheet.Message.Data) != 0 && len(incomeStatement.Message.Data) != 0 {
				key.DividendCapacity = utils.ToFixed((key.DistributableProfit/key.PaidUpCapital)*100, 2)
				key.RetentionRatio = utils.ToFixed((balancesheet.Message.Data[0].Retainedearnings/incomeStatement.Message.Data[0].Netopincome)*100, 2)
			}
			key.FairValue = utils.ToFixed(utils.CalculateGrahamValue(key.Eps, key.Bvps), 2)
			key.DiversionFromFair = utils.ToFixed(((key.LTP-key.FairValue)/(key.FairValue))*100, 2)

			keys = append(keys, key)
		}

		categories := map[string]string{
			"A1": "Ticker", "B1": "LTP", "C1": "%Fair", "D1": "P/E", "E1": "EPS", "F1": "FairValue",
			"G1": "BookValue", "H1": "ROA", "I1": "ROE", "J1": "NPL", "K1": "TotalShare", "L1": "Reserve",
			"M1": "MarketCap", "N1": "DisProfit", "O1": "paidUp", "P1": "ExepectedDividend", "Q1": "PBV",
			"R1": "RetentionRatio", "S1": "Profit/Share",
		}
		folderName := "fundamental"
		if _, err := os.Stat(folderName); os.IsNotExist(err) {
			os.Mkdir(folderName, 0777)
		}
		var excelVals []map[string]interface{}

		for k, v := range keys {
			excelVal := map[string]interface{}{
				utils.GetColumn("A", k): v.Ticker, utils.GetColumn("B", k): v.LTP, utils.GetColumn("C", k): v.DiversionFromFair,
				utils.GetColumn("D", k): v.PE, utils.GetColumn("E", k): v.Eps, utils.GetColumn("F", k): v.FairValue, utils.GetColumn("G", k): v.Bvps,
				utils.GetColumn("H", k): v.Roa, utils.GetColumn("I", k): v.Roe, utils.GetColumn("J", k): v.NPL,
				utils.GetColumn("K", k): v.Listedshares, utils.GetColumn("L", k): v.Reserves, utils.GetColumn("M", k): v.Mktcap,
				utils.GetColumn("N", k): v.DistributableProfit, utils.GetColumn("O", k): v.PaidUpCapital, utils.GetColumn("P", k): v.DividendCapacity,
				utils.GetColumn("Q", k): v.Pbv, utils.GetColumn("R", k): v.RetentionRatio, utils.GetColumn("S", k): v.DistibutableProfitPerShare,
			}
			if v.Ticker != "" {
				excelVals = append(excelVals, excelVal)
			}
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

		if err := f.SaveAs(fmt.Sprintf("%s/%s.xlsx", folderName, sector)); err != nil {
			fmt.Println(err)
		}

	}
}

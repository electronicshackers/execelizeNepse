package controllers

import (
	"fmt"
	"nepse-backend/nepse/bizmandu"
	"nepse-backend/utils"
	"net/http"
	"os"
	"strings"
)

type KeyFinancialMetrics struct {
	Ticker                     string          `json:"ticker"`
	LTP                        float64         `json:"ltp"`
	DiversionFromFair          float64         `json:"divesionFromFair"`
	PE                         float64         `json:"pe"`
	Eps                        float64         `json:"eps"`
	FairValue                  float64         `json:"fairValue"`
	Bvps                       float64         `json:"bvps"`
	Pbv                        float64         `json:"pbv"`
	Roa                        float64         `json:"roa"`
	Roe                        float64         `json:"roe"`
	NPL                        float64         `json:"npl"`
	Listedshares               float64         `json:"listedShares"`
	Reserves                   float64         `json:"reserves"`
	Mktcap                     float64         `json:"mktCap"`
	DistributableProfit        float64         `json:"distributableProfit"`
	DistibutableProfitPerShare float64         `json:"distributableProfitPerShare"`
	PaidUpCapital              float64         `json:"paidUpCapital"`
	DividendCapacity           float64         `json:"dividendCapacity"`
	RetentionRatio             float64         `json:"retentionRatio"`
	Hydro                      HydroKeyMetrics `json:"hydro"`
}

type HydroKeyMetrics struct {
	NetIncome              float64 `json:"netIncome"`
	IncomeFromSaleOfEnergy float64 `json:"incomeFromSaleOfEnergy"`
	CostOfProduction       float64 `json:"costOfProduction"`
	Investements           float64 `json:"investements"`
	WorkInProgress         float64 `json:"workInProgress"`
	CashInHand             float64 `json:"cashInHand"`
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, sector := range sectors {

		sectorStocks, err := biz.GetSectorStock(sector)
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

			if sector == "Hydro Power" {
				if len(incomeStatement.Message.Data) > 0 {
					key.Hydro.CostOfProduction = incomeStatement.Message.Data[0].Costofproduction
					key.Hydro.IncomeFromSaleOfEnergy = incomeStatement.Message.Data[0].Energysales
				}

				if len(financial.Message.Data) != 0 {
					key.Hydro.NetIncome = financial.Message.Data[0].Netincome
				}

				if len(balancesheet.Message.Data) != 0 {
					key.Hydro.CashInHand = balancesheet.Message.Data[0].Cash
					key.Hydro.Investements = balancesheet.Message.Data[0].Investments
					key.Hydro.WorkInProgress = balancesheet.Message.Data[0].Workinprogress
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

		folderName := "fundamental"
		if _, err := os.Stat(folderName); os.IsNotExist(err) {
			os.Mkdir(folderName, 0777)
		}

		categories := GetHeaders(sector)

		var excelVals []map[string]interface{}

		for k, v := range keys {
			excelVal := GetValues(sector, v, k)
			if v.Ticker != "" {
				excelVals = append(excelVals, excelVal)
			}
		}
		utils.CreateExcelFile(folderName, sector, categories, excelVals)
	}
}

func GetHeaders(sector string) map[string]string {
	headers := map[string]string{
		"A1": "Ticker", "B1": "LTP", "C1": "EPS", "D1": "P/E", "E1": "Book Value",
		"F1": "PBV", "G1": "Fair Value", "H1": "ROA", "I1": "ROE", "J1": "Total Share", "K1": "Reserve",
	}

	if sector == "Hydro Power" {
		headers["L1"] = "Net Income"
		headers["M1"] = "Energy Sale"
		headers["N1"] = "Energy Production Cost"
		headers["O1"] = "Work in Progress"
		headers["P1"] = "Cash in Hand"
	}
	return headers
}

func GetValues(sector string, data KeyFinancialMetrics, k int) map[string]interface{} {
	excelVal := map[string]interface{}{
		utils.GetColumn("A", k): data.Ticker, utils.GetColumn("B", k): data.LTP, utils.GetColumn("C", k): data.Eps,
		utils.GetColumn("D", k): data.PE, utils.GetColumn("E", k): data.Bvps, utils.GetColumn("F", k): data.Pbv, utils.GetColumn("G", k): data.FairValue,
		utils.GetColumn("H", k): data.Roa, utils.GetColumn("I", k): data.Roe, utils.GetColumn("J", k): data.Listedshares, utils.GetColumn("K", k): data.Reserves,
	}
	if sector == "Hydro Power" {
		excelVal[fmt.Sprintf("L%d", k+2)] = data.Hydro.NetIncome
		excelVal[fmt.Sprintf("M%d", k+2)] = data.Hydro.IncomeFromSaleOfEnergy
		excelVal[fmt.Sprintf("N%d", k+2)] = data.Hydro.CostOfProduction
		excelVal[fmt.Sprintf("O%d", k+2)] = data.Hydro.WorkInProgress
		excelVal[fmt.Sprintf("P%d", k+2)] = data.Hydro.CashInHand
	}
	return excelVal
}

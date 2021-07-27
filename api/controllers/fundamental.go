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
	Ticker                     string                     `json:"ticker"`
	LTP                        float64                    `json:"ltp"`
	DiversionFromFair          float64                    `json:"divesionFromFair"`
	PE                         float64                    `json:"pe"`
	Eps                        float64                    `json:"eps"`
	FairValue                  float64                    `json:"fairValue"`
	Bvps                       float64                    `json:"bvps"`
	Pbv                        float64                    `json:"pbv"`
	Roa                        float64                    `json:"roa"`
	Roe                        float64                    `json:"roe"`
	NPL                        float64                    `json:"npl"`
	Mktcap                     float64                    `json:"mktCap"`
	DistributableProfit        float64                    `json:"distributableProfit"`
	DistibutableProfitPerShare float64                    `json:"distributableProfitPerShare"`
	PaidUpCapital              float64                    `json:"paidUpCapital"`
	DividendCapacity           float64                    `json:"dividendCapacity"`
	RetentionRatio             float64                    `json:"retentionRatio"`
	Hydro                      HydroKeyMetrics            `json:"hydro"`
	Hotel                      HotelKeyMetrics            `json:"hotel"`
	LifeInsurance              LifeInsuranceKeyMetrics    `json:"lifeInsurance"`
	NonLifeInsurance           NonLifeInsuranceKeyMetrics `json:"nonLifeInsurance"`
	Manufacturing              ManufacturingKeyMetrics    `json:"manufacturing"`
}

type HydroKeyMetrics struct {
	NetIncome              float64 `json:"netIncome"`
	IncomeFromSaleOfEnergy float64 `json:"incomeFromSaleOfEnergy"`
	CostOfProduction       float64 `json:"costOfProduction"`
	Investements           float64 `json:"investements"`
	WorkInProgress         float64 `json:"workInProgress"`
	CashInHand             float64 `json:"cashInHand"`
}

type HotelKeyMetrics struct {
	TotalIncome             float64 `json:"totalIncome"`
	TotalExpenditure        float64 `json:"totalExpenditure"`
	NetIncome               float64 `json:"netIncome"`
	TotalCurrentAssests     float64 `json:"totalCurrentAssests"`
	TotalCurrentLiabilities float64 `json:"totalCurrentLiabilities"`
	NetCurrentAssests       float64 `json:"netCurrentAssests"`
	ReserveAndSurplus       float64 `json:"reserveAndSurplus"`
}

type LifeInsuranceKeyMetrics struct {
	Income             float64 `json:"income"`
	Expenditure        float64 `json:"expenditure"`
	NetIncome          float64 `json:"netIncome"`
	LifeInsuranceFund  float64 `json:"lifeInsuranceFund"`
	ReserveAndSurplus  float64 `json:"reserveAndSurplus"`
	TotalRevenue       float64 `json:"totalRevenue"`
	GrossProfit        float64 `json:"grossProfit"`
	CatastropheReserve float64 `json:"catastropheReserve"`
}

type NonLifeInsuranceKeyMetrics struct {
	Income             float64 `json:"income"`
	Expenditure        float64 `json:"expenditure"`
	NetIncome          float64 `json:"netIncome"`
	InsuranceFund      float64 `json:"InsuranceFund"`
	ReserveAndSurplus  float64 `json:"reserveAndSurplus"`
	TotalRevenue       float64 `json:"totalRevenue"`
	GrossProfit        float64 `json:"grossProfit"`
	CatastropheReserve float64 `json:"catastropheReserve"`
}

type ManufacturingKeyMetrics struct {
	TotalIncome               float64 `json:"totalIncome"`
	TotalExpenditure          float64 `json:"totalExpenditure"`
	NetIncome                 float64 `json:"netIncome"`
	TotalEquityAndLiabilities float64 `json:"totalEquityAndLiabilities"`
	TotalRevenue              float64 `json:"totalRevenue"`
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

			if sector == "Hotels" {
				if len(incomeStatement.Message.Data) > 0 {
					key.Hotel.TotalIncome = incomeStatement.Message.Data[0].Totalincome
					key.Hotel.TotalExpenditure = incomeStatement.Message.Data[0].Totalexpenditure
				}

				if len(balancesheet.Message.Data) != 0 {
					key.Hotel.TotalCurrentAssests = balancesheet.Message.Data[0].Totalcurrentassets
					key.Hotel.TotalCurrentLiabilities = balancesheet.Message.Data[0].Totalcurrentliabilities
					key.Hotel.NetCurrentAssests = balancesheet.Message.Data[0].Netcurrentassets
					key.Hotel.ReserveAndSurplus = balancesheet.Message.Data[0].Reservesurplus
				}

				if len(financial.Message.Data) != 0 {
					key.Hotel.NetIncome = financial.Message.Data[0].Netincome
				}
			}

			if sector == "Life Insurance" {

				if len(incomeStatement.Message.Data) > 0 {
					key.LifeInsurance.Income = incomeStatement.Message.Data[0].Income
					key.LifeInsurance.Expenditure = incomeStatement.Message.Data[0].Expenses
				}

				if len(balancesheet.Message.Data) != 0 {
					key.LifeInsurance.LifeInsuranceFund = balancesheet.Message.Data[0].Lifeinsurancefund
					key.LifeInsurance.CatastropheReserve = balancesheet.Message.Data[0].Catastrophereserve
					key.LifeInsurance.ReserveAndSurplus = balancesheet.Message.Data[0].Reservesurplus
				}

				if len(financial.Message.Data) != 0 {
					key.LifeInsurance.NetIncome = financial.Message.Data[0].Netincome
					key.LifeInsurance.GrossProfit = financial.Message.Data[0].Grossprofit
					key.LifeInsurance.TotalRevenue = financial.Message.Data[0].Totalrevenue
				}
			}

			if sector == "Non Life Insurance" {
				if len(incomeStatement.Message.Data) > 0 {
					key.NonLifeInsurance.Income = incomeStatement.Message.Data[0].Income
					key.NonLifeInsurance.Expenditure = incomeStatement.Message.Data[0].Expenses
				}
				if len(balancesheet.Message.Data) != 0 {
					key.NonLifeInsurance.InsuranceFund = balancesheet.Message.Data[0].Insurancefund
					key.NonLifeInsurance.ReserveAndSurplus = balancesheet.Message.Data[0].Reservesurplus
					key.NonLifeInsurance.CatastropheReserve = balancesheet.Message.Data[0].Catastrophereserve
				}
				if len(financial.Message.Data) != 0 {
					key.NonLifeInsurance.NetIncome = financial.Message.Data[0].Netincome
					key.NonLifeInsurance.GrossProfit = financial.Message.Data[0].Grossprofit
					key.NonLifeInsurance.TotalRevenue = financial.Message.Data[0].Totalrevenue
				}
			}

			if sector == "Manufacturing And Processing" {
				if len(incomeStatement.Message.Data) > 0 {
					key.Manufacturing.TotalIncome = incomeStatement.Message.Data[0].Totalincome
					key.Manufacturing.TotalExpenditure = incomeStatement.Message.Data[0].Totalexpenditure
					key.Manufacturing.NetIncome = incomeStatement.Message.Data[0].Netprofit
				}
				if len(balancesheet.Message.Data) != 0 {
					key.Manufacturing.TotalEquityAndLiabilities = balancesheet.Message.Data[0].Totaleqli
				}
				if len(financial.Message.Data) != 0 {
					key.Manufacturing.TotalRevenue = financial.Message.Data[0].Totalrevenue
				}
			}

			if len(financial.Message.Data) != 0 {
				key.NPL = utils.ToFixed(financial.Message.Data[0].Nonperformingloannpltototalloan*100, 2)
			}

			if len(balancesheet.Message.Data) != 0 {
				key.PaidUpCapital = float64(balancesheet.Message.Data[0].Paidupcapital)
			}

			// if len(incomeStatement.Message.Data) != 0 {
			// 	key.DistributableProfit = float64(incomeStatement.Message.Data[0].Freeprofit)
			// 	key.DistibutableProfitPerShare = utils.ToFixed((key.DistributableProfit/key.Listedshares)*100, 2)
			// }

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
		"F1": "PBV", "G1": "Fair Value", "H1": "ROA", "I1": "ROE",
	}

	if sector == "Hydro Power" {
		headers["J1"] = "Net Income"
		headers["K1"] = "Energy Sale"
		headers["L1"] = "Energy Production Cost"
		headers["M1"] = "Work in Progress"
		headers["N1"] = "Cash in Hand"
	}

	if sector == "Hotels" {
		headers["J1"] = "Total Income"
		headers["K1"] = "Total Expenditure"
		headers["L1"] = "Net Income"
		headers["M1"] = "Total Current Assests"
		headers["N1"] = "Total Current Liabilities"
		headers["O1"] = "Net Current Assets"
		headers["P1"] = "Reserve & Surplus"
	}

	if sector == "Life Insurance" {
		headers["J1"] = "Income"
		headers["K1"] = "Expenditure"
		headers["L1"] = "Net Income"
		headers["M1"] = "Life Insurance Fund"
		headers["N1"] = "Catastrophe Reserve"
		headers["O1"] = "Total Revenue"
		headers["P1"] = "Gross Profit"
		headers["Q1"] = "Reserve"
	}

	if sector == "Non Life Insurance" {
		headers["J1"] = "Income"
		headers["K1"] = "Expenditure"
		headers["L1"] = "Net Income"
		headers["M1"] = "Insurance Fund"
		headers["N1"] = "Catastrophe Reserve"
		headers["O1"] = "Total Revenue"
		headers["P1"] = "Gross Profit"
		headers["Q1"] = "Reserve"
	}

	if sector == "Manufacturing And Processing" {
		headers["J1"] = "Income"
		headers["K1"] = "Expenditure"
		headers["L1"] = "Net Profit"
		headers["M1"] = "Total Revenue"
		headers["N1"] = "Total Equity"
	}

	return headers
}

func GetValues(sector string, data KeyFinancialMetrics, k int) map[string]interface{} {
	excelVal := map[string]interface{}{
		utils.GetColumn("A", k): data.Ticker, utils.GetColumn("B", k): data.LTP, utils.GetColumn("C", k): data.Eps,
		utils.GetColumn("D", k): data.PE, utils.GetColumn("E", k): data.Bvps, utils.GetColumn("F", k): data.Pbv, utils.GetColumn("G", k): data.FairValue,
		utils.GetColumn("H", k): data.Roa, utils.GetColumn("I", k): data.Roe,
	}
	if sector == "Hydro Power" {
		excelVal[fmt.Sprintf("J%d", k+2)] = data.Hydro.NetIncome
		excelVal[fmt.Sprintf("K%d", k+2)] = data.Hydro.IncomeFromSaleOfEnergy
		excelVal[fmt.Sprintf("L%d", k+2)] = data.Hydro.CostOfProduction
		excelVal[fmt.Sprintf("M%d", k+2)] = data.Hydro.WorkInProgress
		excelVal[fmt.Sprintf("N%d", k+2)] = data.Hydro.CashInHand
	}
	if sector == "Hotels" {
		excelVal[fmt.Sprintf("J%d", k+2)] = data.Hotel.TotalIncome
		excelVal[fmt.Sprintf("K%d", k+2)] = data.Hotel.TotalExpenditure
		excelVal[fmt.Sprintf("L%d", k+2)] = data.Hotel.NetIncome
		excelVal[fmt.Sprintf("M%d", k+2)] = data.Hotel.TotalCurrentAssests
		excelVal[fmt.Sprintf("N%d", k+2)] = data.Hotel.TotalCurrentLiabilities
		excelVal[fmt.Sprintf("O%d", k+2)] = data.Hotel.NetCurrentAssests
		excelVal[fmt.Sprintf("P%d", k+2)] = data.Hotel.ReserveAndSurplus
	}

	if sector == "Life Insurance" {
		excelVal[fmt.Sprintf("J%d", k+2)] = data.LifeInsurance.Income
		excelVal[fmt.Sprintf("K%d", k+2)] = data.LifeInsurance.Expenditure
		excelVal[fmt.Sprintf("L%d", k+2)] = data.LifeInsurance.NetIncome
		excelVal[fmt.Sprintf("M%d", k+2)] = data.LifeInsurance.LifeInsuranceFund
		excelVal[fmt.Sprintf("N%d", k+2)] = data.LifeInsurance.CatastropheReserve
		excelVal[fmt.Sprintf("O%d", k+2)] = data.LifeInsurance.TotalRevenue
		excelVal[fmt.Sprintf("P%d", k+2)] = data.LifeInsurance.GrossProfit
		excelVal[fmt.Sprintf("Q%d", k+2)] = data.LifeInsurance.ReserveAndSurplus
	}
	if sector == "Non Life Insurance" {
		excelVal[fmt.Sprintf("J%d", k+2)] = data.NonLifeInsurance.Income
		excelVal[fmt.Sprintf("K%d", k+2)] = data.NonLifeInsurance.Expenditure
		excelVal[fmt.Sprintf("L%d", k+2)] = data.NonLifeInsurance.NetIncome
		excelVal[fmt.Sprintf("M%d", k+2)] = data.NonLifeInsurance.InsuranceFund
		excelVal[fmt.Sprintf("N%d", k+2)] = data.NonLifeInsurance.CatastropheReserve
		excelVal[fmt.Sprintf("O%d", k+2)] = data.NonLifeInsurance.TotalRevenue
		excelVal[fmt.Sprintf("P%d", k+2)] = data.NonLifeInsurance.GrossProfit
		excelVal[fmt.Sprintf("Q%d", k+2)] = data.NonLifeInsurance.ReserveAndSurplus
	}
	if sector == "Manufacturing And Processing" {
		excelVal[fmt.Sprintf("J%d", k+2)] = data.Manufacturing.TotalIncome
		excelVal[fmt.Sprintf("K%d", k+2)] = data.Manufacturing.TotalExpenditure
		excelVal[fmt.Sprintf("L%d", k+2)] = data.Manufacturing.NetIncome
		excelVal[fmt.Sprintf("M%d", k+2)] = data.Manufacturing.TotalRevenue
		excelVal[fmt.Sprintf("N%d", k+2)] = data.Manufacturing.TotalEquityAndLiabilities
	}

	return excelVal
}

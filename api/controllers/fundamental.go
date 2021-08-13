package controllers

import (
	"fmt"
	"log"
	"nepse-backend/nepse"
	"nepse-backend/nepse/bizmandu"
	"nepse-backend/nepse/neweb"
	"nepse-backend/utils"
	"net/http"
	"os"
	"runtime/pprof"
	"strings"
	"time"
)

const (
	HydroPower                 = "Hydro Power"
	OrganizedFund              = "Organized Fund"
	LifeInsurance              = "Life Insurance"
	Microcredit                = "Microcredit"
	DevelopmentBank            = "Development Bank"
	Hotels                     = "Hotels"
	NonLifeInsurance           = "Non Life Insurance"
	Finance                    = "Finance"
	CommercialBanks            = "Commercial Banks"
	Trading                    = "Trading"
	ManufacturingAndProcessing = "Manufacturing And Processing"
	Telecom                    = "Telecom"
)

type KeyFinancialMetrics struct {
	Ticker            string                     `json:"ticker"`
	LTP               float64                    `json:"ltp"`
	DiversionFromFair float64                    `json:"divesionFromFair"`
	PE                float64                    `json:"pe"`
	Eps               float64                    `json:"eps"`
	FairValue         float64                    `json:"fairValue"`
	Bvps              float64                    `json:"bvps"`
	Pbv               float64                    `json:"pbv"`
	Roa               float64                    `json:"roa"`
	Roe               float64                    `json:"roe"`
	Mktcap            float64                    `json:"mktCap"`
	Quarter           float64                    `json:"quarter"`
	PaidUpCapital     float64                    `json:"paidUpCapital"`
	Hydro             HydroKeyMetrics            `json:"hydro"`
	Hotel             HotelKeyMetrics            `json:"hotel"`
	LifeInsurance     LifeInsuranceKeyMetrics    `json:"lifeInsurance"`
	NonLifeInsurance  NonLifeInsuranceKeyMetrics `json:"nonLifeInsurance"`
	Manufacturing     ManufacturingKeyMetrics    `json:"manufacturing"`
	BFI               BFIKeyMetrics              `json:"bfi"`
	Microcredit       MicrocreditKeyMetrics      `json:"microcredit"`
}

type MicrocreditKeyMetrics struct {
	Assests                   float64 `json:"assests"`
	NetInterestIncome         float64 `json:"netInterestIncome"`
	NetIncome                 float64 `json:"netIncome"`
	Reserves                  float64 `json:"reserves"`
	NetInterestIncomePerShare float64 `json:"netInterestIncomePerShare"`
}

type BFIKeyMetrics struct {
	NPL                        float64 `json:"npl"`
	DistributableProfit        float64 `json:"distributableProfit"`
	DistibutableProfitPerShare float64 `json:"distributableProfitPerShare"`
	DividendCapacity           float64 `json:"dividendCapacity"`
	Reserve                    float64 `json:"reserve"`
}

type HydroKeyMetrics struct {
	NetIncome              float64 `json:"netIncome"`
	IncomeFromSaleOfEnergy float64 `json:"incomeFromSaleOfEnergy"`
	CostOfProduction       float64 `json:"costOfProduction"`
	IncomeByCost           float64 `json:"incomeByCost"`
	Investements           float64 `json:"investements"`
	Reserves               float64 `json:"reserves"`
	NetInterestIncome      float64 `json:"netInterestIncome"`
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
	TotalPremium        float64 `json:"totalPremium"`
	TotalNumberOfPolicy float64 `json:"totalNumberOfPolicy"`
	NetIncome           float64 `json:"netIncome"`
	LifeInsuranceFund   float64 `json:"lifeInsuranceFund"`
	ReserveAndSurplus   float64 `json:"reserveAndSurplus"`
	TotalRevenue        float64 `json:"totalRevenue"`
	TotalInvestment     float64 `json:"totalInvestment"`
}

type NonLifeInsuranceKeyMetrics struct {
	TotalPolicy       float64 `json:"totalPolicy"`
	RenewedPolicy     float64 `json:"renewedPolicy"`
	TotalByRenewed    float64 `json:"totalByRenewed"`
	NetIncome         float64 `json:"netIncome"`
	InsuranceFund     float64 `json:"InsuranceFund"`
	ReserveAndSurplus float64 `json:"reserveAndSurplus"`
	TotalRevenue      float64 `json:"totalRevenue"`
	TotalInvestment   float64 `json:"totalInvestment"`
}

type ManufacturingKeyMetrics struct {
	TotalIncome               float64 `json:"totalIncome"`
	TotalExpenditure          float64 `json:"totalExpenditure"`
	NetIncome                 float64 `json:"netIncome"`
	TotalEquityAndLiabilities float64 `json:"totalEquityAndLiabilities"`
	TotalRevenue              float64 `json:"totalRevenue"`
}

func (server *Server) GetFundamentalSectorwise(w http.ResponseWriter, r *http.Request) {
	f, err := os.Create("Fundamental.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	start := time.Now()
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

	stockMap := make(map[string][]nepse.Ticker)

	for stockIndex := range nepseSectors {
		nepseStocks := nepseSectors[stockIndex]
		stockMap[nepseStocks.Ticker] = append(stockMap[nepseStocks.Ticker], nepse.Ticker{
			Ticker:          nepseStocks.Ticker,
			Id:              nepseStocks.Id,
			Companyname:     nepseStocks.Companyname,
			Sector:          nepseStocks.Sector,
			Lasttradedprice: nepseStocks.Lasttradedprice,
		})
	}

	bizSectors, err := biz.GetStocks()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, sector := range sectors {

		sectorStocks, err := biz.GetSectorStock(sector, nepseSectors, bizSectors)
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

			key.Eps = utils.ToFixed(detail.Message.Summary.Epsdiluted, 2)
			key.PE = utils.ToFixed(detail.Message.Summary.Pediluted, 2)
			key.LTP = stockMap[ticker.Ticker][0].Lasttradedprice
			key.Bvps = utils.ToFixed(detail.Message.Summary.Bvps, 2)
			key.Ticker = detail.Message.Keyfinancial.Ticker
			key.Mktcap = detail.Message.Summary.Mktcap
			key.Pbv = utils.ToFixed(key.LTP/key.Bvps, 2)

			if len(incomeStatement.Message.Data) != 0 {
				key.Quarter = incomeStatement.Message.Data[0].Quarter
			}

			for _, quarter := range detail.Message.Keyfinancial.Data {
				if quarter.Type == "CURRENT" {
					key.Roa = utils.ToFixed(quarter.Roa*100, 2)
					key.Roe = utils.ToFixed(quarter.Roe*100, 2)
				}
			}

			if sector == Microcredit {
				if len(balancesheet.Message.Data) != 0 {
					key.Microcredit.Assests = utils.ToFixed(balancesheet.Message.Data[0].Totalassets, 2)
					key.Microcredit.Reserves = utils.ToFixed(balancesheet.Message.Data[0].Reservesandsurplus, 2)
				}

				if len(incomeStatement.Message.Data) != 0 {
					key.Microcredit.NetInterestIncome = utils.ToFixed(incomeStatement.Message.Data[0].NetInterestincome, 2)
					key.Microcredit.NetIncome = utils.ToFixed(incomeStatement.Message.Data[0].Netprofitorloss, 2)

				}

				if len(financial.Message.Data) != 0 {
					totalShare := financial.Message.Data[0].Outstandingshares
					key.Microcredit.NetInterestIncomePerShare = utils.ToFixed((key.Microcredit.NetInterestIncome / totalShare), 2)
				}
			}

			if sector == HydroPower {
				if len(incomeStatement.Message.Data) > 0 {
					key.Hydro.CostOfProduction = incomeStatement.Message.Data[0].Costofproduction
					key.Hydro.IncomeFromSaleOfEnergy = incomeStatement.Message.Data[0].Energysales
					key.Hydro.IncomeByCost = utils.ToFixed((key.Hydro.IncomeFromSaleOfEnergy / key.Hydro.CostOfProduction), 2)
					key.Hydro.NetInterestIncome = utils.ToFixed(incomeStatement.Message.Data[0].Interestincomeexpense, 2)
				}

				if len(financial.Message.Data) != 0 {
					key.Hydro.NetIncome = financial.Message.Data[0].Netincome
				}

				if len(balancesheet.Message.Data) != 0 {
					key.Hydro.Investements = balancesheet.Message.Data[0].Investments
					key.Hydro.Reserves = balancesheet.Message.Data[0].Reserves
				}
			}

			if sector == Hotels {
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

			if sector == LifeInsurance {
				if len(balancesheet.Message.Data) != 0 {
					key.LifeInsurance.LifeInsuranceFund = balancesheet.Message.Data[0].Lifeinsurancefund
					key.LifeInsurance.TotalInvestment = balancesheet.Message.Data[0].Shortterminvestmentsandloans + balancesheet.Message.Data[0].Longterminvestmentsloan
					key.LifeInsurance.ReserveAndSurplus = balancesheet.Message.Data[0].Reservesurplus
				}

				if len(financial.Message.Data) != 0 {
					key.LifeInsurance.NetIncome = financial.Message.Data[0].Netincome
					key.LifeInsurance.TotalRevenue = financial.Message.Data[0].Totalrevenue
					key.LifeInsurance.TotalPremium = financial.Message.Data[0].Totalpremium
					key.LifeInsurance.TotalNumberOfPolicy = financial.Message.Data[0].Totalnoofpolicies
				}
			}

			if sector == NonLifeInsurance {
				if len(balancesheet.Message.Data) != 0 {
					key.NonLifeInsurance.InsuranceFund = balancesheet.Message.Data[0].Insurancefund
					key.NonLifeInsurance.ReserveAndSurplus = balancesheet.Message.Data[0].Reservesurplus
					key.NonLifeInsurance.TotalInvestment = balancesheet.Message.Data[0].Shortterminvestmentsandloans + balancesheet.Message.Data[0].Longterminvestmentsloan
				}
				if len(financial.Message.Data) != 0 {
					key.NonLifeInsurance.NetIncome = financial.Message.Data[0].Netincome
					key.NonLifeInsurance.TotalPolicy = financial.Message.Data[0].Totalnoofpolicies
					key.NonLifeInsurance.RenewedPolicy = financial.Message.Data[0].Totalrenewedpolicies
					key.NonLifeInsurance.TotalRevenue = financial.Message.Data[0].Totalrevenue
					key.NonLifeInsurance.TotalByRenewed = utils.ToFixed((key.NonLifeInsurance.TotalPolicy / key.NonLifeInsurance.RenewedPolicy), 2)
				}
			}

			if sector == ManufacturingAndProcessing {
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

			if sector == CommercialBanks || sector == DevelopmentBank || sector == Finance {
				listedShares := detail.Message.Summary.Listedshares
				if len(incomeStatement.Message.Data) != 0 {
					key.BFI.DistributableProfit = float64(incomeStatement.Message.Data[0].Freeprofit)
					key.BFI.DistibutableProfitPerShare = utils.ToFixed((key.BFI.DistributableProfit/listedShares)*100, 2)
				}

				if len(financial.Message.Data) != 0 {
					key.BFI.NPL = utils.ToFixed(financial.Message.Data[0].Nonperformingloannpltototalloan*100, 2)
					key.PaidUpCapital = utils.ToFixed(financial.Message.Data[0].Outstandingshares, 2)
				}

				if len(balancesheet.Message.Data) != 0 && len(incomeStatement.Message.Data) != 0 {
					key.BFI.Reserve = float64(balancesheet.Message.Data[0].Reserves)
					paidUp := float64(balancesheet.Message.Data[0].Paidupcapital)
					key.BFI.DividendCapacity = utils.ToFixed((key.BFI.DistributableProfit/paidUp)*100, 2)
				}

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
		go utils.CreateExcelFile(folderName, sector, categories, excelVals)
	}
	duration := time.Since(start)
	fmt.Println(duration)
	fmt.Println(duration.Nanoseconds())
}

func GetHeaders(sector string) map[string]string {
	headers := map[string]string{
		"A1": "Ticker", "B1": "LTP", "C1": "EPS", "D1": "P/E", "E1": "Book Value",
		"F1": "PBV", "G1": "Fair Value", "H1": "ROA", "I1": "ROE",
	}

	if sector == Microcredit {
		headers["J1"] = "Assests"
		headers["K1"] = "NetInterestIncome"
		headers["L1"] = "NetIncome"
		headers["M1"] = "Reserves"
		headers["N1"] = "NetInterestIncomePerShare"
		headers["O1"] = "Quarter"
	}

	if sector == CommercialBanks || sector == DevelopmentBank || sector == Finance {
		headers["J1"] = "NPL"
		headers["K1"] = "Reserve"
		headers["L1"] = "Distributable Profit"
		headers["M1"] = "Dividend Capacity"
		headers["N1"] = "Profit/Share"
		headers["O1"] = "Quarter"
		headers["P1"] = "Outstanding Shares"
	}

	if sector == HydroPower {
		headers["J1"] = "Net Income"
		headers["K1"] = "Energy Sale"
		headers["L1"] = "Production Cost"
		headers["M1"] = "Sale/Production"
		headers["N1"] = "Investement"
		headers["O1"] = "Reserves"
		headers["P1"] = "Interest Income"
		headers["Q1"] = "Quarter"
	}

	if sector == Hotels {
		headers["J1"] = "Total Income"
		headers["K1"] = "Total Expenditure"
		headers["L1"] = "Net Income"
		headers["M1"] = "Total Current Assests"
		headers["N1"] = "Total Current Liabilities"
		headers["O1"] = "Net Current Assets"
		headers["P1"] = "Reserve & Surplus"
		headers["Q1"] = "Quarter"
	}

	if sector == LifeInsurance {
		headers["J1"] = "Total Premium"
		headers["K1"] = "Total Policy"
		headers["L1"] = "Net Income"
		headers["M1"] = "Life Insurance Fund"
		headers["N1"] = "Total Investment"
		headers["O1"] = "Total Revenue"
		headers["P1"] = "Reserve"
		headers["Q1"] = "Quarter"
	}

	if sector == NonLifeInsurance {
		headers["J1"] = "TotalPolicy"
		headers["K1"] = "TotalRenewedPolicy"
		headers["L1"] = "Total/Renewed"
		headers["M1"] = "Life Insurance Fund"
		headers["N1"] = "Total Investment"
		headers["O1"] = "Total Revenue"
		headers["P1"] = "Net Income"
		headers["Q1"] = "Reserve"
		headers["R1"] = "Quarter"
	}

	if sector == ManufacturingAndProcessing {
		headers["J1"] = "Income"
		headers["K1"] = "Expenditure"
		headers["L1"] = "Net Profit"
		headers["M1"] = "Total Revenue"
		headers["N1"] = "Total Equity"
		headers["O1"] = "Quarter"
	}

	return headers
}

func GetValues(sector string, data KeyFinancialMetrics, k int) map[string]interface{} {
	excelVal := map[string]interface{}{
		utils.GetColumn("A", k): data.Ticker, utils.GetColumn("B", k): data.LTP, utils.GetColumn("C", k): data.Eps,
		utils.GetColumn("D", k): data.PE, utils.GetColumn("E", k): data.Bvps, utils.GetColumn("F", k): data.Pbv, utils.GetColumn("G", k): data.FairValue,
		utils.GetColumn("H", k): data.Roa, utils.GetColumn("I", k): data.Roe,
	}
	if sector == CommercialBanks || sector == DevelopmentBank || sector == Finance {
		excelVal[utils.GetColumn("J", k)] = data.BFI.NPL
		excelVal[utils.GetColumn("K", k)] = data.BFI.Reserve
		excelVal[utils.GetColumn("L", k)] = data.BFI.DistributableProfit
		excelVal[utils.GetColumn("M", k)] = data.BFI.DividendCapacity
		excelVal[utils.GetColumn("N", k)] = data.BFI.DistibutableProfitPerShare
		excelVal[utils.GetColumn("O", k)] = data.Quarter
		excelVal[utils.GetColumn("P", k)] = data.PaidUpCapital
	}

	if sector == HydroPower {
		excelVal[fmt.Sprintf("J%d", k+2)] = data.Hydro.NetIncome
		excelVal[fmt.Sprintf("K%d", k+2)] = data.Hydro.IncomeFromSaleOfEnergy
		excelVal[fmt.Sprintf("L%d", k+2)] = data.Hydro.CostOfProduction
		excelVal[fmt.Sprintf("M%d", k+2)] = data.Hydro.IncomeByCost
		excelVal[fmt.Sprintf("N%d", k+2)] = data.Hydro.Investements
		excelVal[fmt.Sprintf("O%d", k+2)] = data.Hydro.Reserves
		excelVal[fmt.Sprintf("P%d", k+2)] = data.Hydro.NetInterestIncome
		excelVal[fmt.Sprintf("Q%d", k+2)] = data.Quarter
	}
	if sector == Hotels {
		excelVal[fmt.Sprintf("J%d", k+2)] = data.Hotel.TotalIncome
		excelVal[fmt.Sprintf("K%d", k+2)] = data.Hotel.TotalExpenditure
		excelVal[fmt.Sprintf("L%d", k+2)] = data.Hotel.NetIncome
		excelVal[fmt.Sprintf("M%d", k+2)] = data.Hotel.TotalCurrentAssests
		excelVal[fmt.Sprintf("N%d", k+2)] = data.Hotel.TotalCurrentLiabilities
		excelVal[fmt.Sprintf("O%d", k+2)] = data.Hotel.NetCurrentAssests
		excelVal[fmt.Sprintf("P%d", k+2)] = data.Hotel.ReserveAndSurplus
		excelVal[fmt.Sprintf("Q%d", k+2)] = data.Quarter
	}

	if sector == LifeInsurance {
		excelVal[fmt.Sprintf("J%d", k+2)] = data.LifeInsurance.TotalPremium
		excelVal[fmt.Sprintf("K%d", k+2)] = data.LifeInsurance.TotalNumberOfPolicy
		excelVal[fmt.Sprintf("L%d", k+2)] = data.LifeInsurance.NetIncome
		excelVal[fmt.Sprintf("M%d", k+2)] = data.LifeInsurance.LifeInsuranceFund
		excelVal[fmt.Sprintf("N%d", k+2)] = data.LifeInsurance.TotalInvestment
		excelVal[fmt.Sprintf("O%d", k+2)] = data.LifeInsurance.TotalRevenue
		excelVal[fmt.Sprintf("P%d", k+2)] = data.LifeInsurance.ReserveAndSurplus
		excelVal[fmt.Sprintf("Q%d", k+2)] = data.Quarter
	}

	if sector == NonLifeInsurance {
		excelVal[fmt.Sprintf("J%d", k+2)] = data.NonLifeInsurance.TotalPolicy
		excelVal[fmt.Sprintf("K%d", k+2)] = data.NonLifeInsurance.RenewedPolicy
		excelVal[fmt.Sprintf("L%d", k+2)] = data.NonLifeInsurance.TotalByRenewed
		excelVal[fmt.Sprintf("M%d", k+2)] = data.NonLifeInsurance.InsuranceFund
		excelVal[fmt.Sprintf("N%d", k+2)] = data.NonLifeInsurance.TotalInvestment
		excelVal[fmt.Sprintf("O%d", k+2)] = data.NonLifeInsurance.TotalRevenue
		excelVal[fmt.Sprintf("P%d", k+2)] = data.NonLifeInsurance.NetIncome
		excelVal[fmt.Sprintf("Q%d", k+2)] = data.NonLifeInsurance.ReserveAndSurplus
		excelVal[fmt.Sprintf("R%d", k+2)] = data.Quarter
	}

	if sector == ManufacturingAndProcessing {
		excelVal[fmt.Sprintf("J%d", k+2)] = data.Manufacturing.TotalIncome
		excelVal[fmt.Sprintf("K%d", k+2)] = data.Manufacturing.TotalExpenditure
		excelVal[fmt.Sprintf("L%d", k+2)] = data.Manufacturing.NetIncome
		excelVal[fmt.Sprintf("M%d", k+2)] = data.Manufacturing.TotalRevenue
		excelVal[fmt.Sprintf("N%d", k+2)] = data.Manufacturing.TotalEquityAndLiabilities
		excelVal[fmt.Sprintf("O%d", k+2)] = data.Quarter
	}

	if sector == Microcredit {
		excelVal[fmt.Sprintf("J%d", k+2)] = data.Microcredit.Assests
		excelVal[fmt.Sprintf("K%d", k+2)] = data.Microcredit.NetInterestIncome
		excelVal[fmt.Sprintf("L%d", k+2)] = data.Microcredit.NetIncome
		excelVal[fmt.Sprintf("M%d", k+2)] = data.Microcredit.Reserves
		excelVal[fmt.Sprintf("N%d", k+2)] = data.Microcredit.NetInterestIncomePerShare
		excelVal[fmt.Sprintf("O%d", k+2)] = data.Quarter
	}
	return excelVal
}

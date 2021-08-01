package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	responses "nepse-backend/api/response"
	"nepse-backend/nepse/bizmandu"
	"nepse-backend/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type WhatIfRequest struct {
	Stock    string    `json:"stock"`
	Amount   float64   `json:"amount"`
	Quantity int       `json:"quantity"`
	BuyDate  string    `json:"buyDate"`
	SellDate time.Time `json:"sellDate"`
}

func (server *Server) WhatIf(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	whatIf := &WhatIfRequest{}
	err = json.Unmarshal(body, whatIf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buyDate, err := utils.StringToTime(whatIf.BuyDate)
	if err != nil {
		http.Error(w, "Start date is invalid", http.StatusBadRequest)
		return
	}

	buyTime := buyDate.Unix()

	timeNow := time.Now().Unix()

	biz, err := bizmandu.NewBizmandu()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dividend, err := biz.GetDividendHistory(whatIf.Stock)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prices, err := biz.GetPriceHistory(whatIf.Stock, buyTime, timeNow)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	indexOfTime := closest(prices.T, int(buyTime))

	priceAtTime := prices.C[indexOfTime]

	if whatIf.Amount != 0 {
		whatIf.Quantity = int(math.Floor(whatIf.Amount / priceAtTime))
	}

	var yield DividendYield

	year := time.Unix(buyTime, 0).Year()
	yield = calculateDividendYield(dividend, year, whatIf.Quantity)

	indexOfSellingTime := closest(prices.T, int(timeNow))
	sellingPrice := prices.C[indexOfSellingTime]

	yield.InitialCost = yield.InitialQuantity * priceAtTime
	yield.TotalValue = sellingPrice*yield.TotalQuantity + yield.TotalCashYield
	yield.Profit = yield.TotalValue - yield.InitialCost

	yield.ProfitPercentage = fmt.Sprintf("%.2f%%", utils.ToFixed((yield.Profit/yield.InitialCost)*100, 2))

	responses.JSON(w, 200, yield)
}

func closest(array []int, num int) int {
	current := array[0]
	for _, numbers := range array {
		if math.Abs(float64(num-numbers)) < math.Abs(float64(num-current)) {
			current = numbers
		}
	}
	return indexOf(current, array)
}

func indexOf(element int, data []int) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}

type DividendYield struct {
	TotalCashYield   float64               `json:"totalCashYield"`
	TotalQuantity    float64               `json:"totalQuantity"`
	InitialQuantity  float64               `json:"initialQuantity"`
	InitialCost      float64               `json:"initialCost"`
	TotalValue       float64               `json:"totalValue"`
	Profit           float64               `json:"profit"`
	ProfitPercentage string                `json:"profitPercentage"`
	YearlyYield      []YearlyDividendYield `json:"yearlyYield"`
}

type YearlyDividendYield struct {
	Year        string  `json:"year"`
	Cash        string  `json:"cash"`
	Bonus       string  `json:"bonus"`
	CashAfter   float64 `json:"cashAfter"`
	CashBefore  float64 `json:"cashBefore"`
	ShareAfter  float64 `json:"shareAfter"`
	ShareBefore float64 `json:"shareBefore"`
}

func calculateDividendYield(dividend *bizmandu.DividendHistory, year int, quanity int) DividendYield {
	var dividendYield DividendYield
	dividendYield.TotalQuantity = float64(quanity)
	dividendYield.InitialQuantity = float64(quanity)
	for _, div := range dividend.Message.Dividend {
		if div.Bonus != 0 || div.Cash != 0 {
			divYear := strings.Split(div.Year, "/")[0]
			divYearInt, _ := strconv.Atoi(divYear)

			if divYearInt > year {

				dividendYield.YearlyYield = append(dividendYield.YearlyYield, YearlyDividendYield{
					Year:        div.Year,
					Cash:        fmt.Sprintf("%.2f%%", div.Cash*100),
					Bonus:       fmt.Sprintf("%.2f%%", div.Bonus*100),
					CashBefore:  utils.ToFixed(dividendYield.TotalCashYield, 2),
					CashAfter:   utils.ToFixed(dividendYield.TotalCashYield+div.Cash*100*dividendYield.TotalQuantity, 2),
					ShareBefore: math.Floor(dividendYield.TotalQuantity),
					ShareAfter:  math.Floor(dividendYield.TotalQuantity + div.Bonus*dividendYield.TotalQuantity),
				})
				dividendYield.TotalCashYield = utils.ToFixed(dividendYield.TotalCashYield+div.Cash*100*dividendYield.TotalQuantity, 2)
				dividendYield.TotalQuantity = math.Floor(dividendYield.TotalQuantity + div.Bonus*dividendYield.TotalQuantity)

				fmt.Println("dividendYield.TotalCashYield", dividendYield.TotalCashYield, "Year", divYearInt)
			}
		}

	}
	return dividendYield
}

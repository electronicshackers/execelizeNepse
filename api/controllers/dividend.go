package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	responses "nepse-backend/api/response"
	"nepse-backend/nepse/bizmandu"
	"nepse-backend/utils"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Rights struct {
	Data []CompanyRights `json:"data"`
}

type Dividends struct {
	Data []CompanyDividend `json:"data"`
}

type CompanyDividend struct {
	Sn                  int       `json:"sn"`
	Dividendhistoryid   int       `json:"dividendHistoryId"`
	Companyname         string    `json:"companyName"`
	Stocksymbol         string    `json:"stockSymbol"`
	Sectorname          string    `json:"sectorName"`
	Filename            string    `json:"fileName"`
	Bonus               string    `json:"bonus"`
	Cash                string    `json:"cash"`
	Totaldividend       string    `json:"totalDividend"`
	Distributiondatead  string    `json:"distributionDateAD"`
	Distributiondatebs  string    `json:"distributionDateBS"`
	Bookclosuredatead   string    `json:"bookClosureDateAD"`
	Bookclosuredatetime time.Time `json:"bookClosureDateTime"`
	Bookclosuredatebs   string    `json:"bookClosureDateBS"`
	Fiscalyearad        string    `json:"fiscalYearAD"`
	Fiscalyearbs        string    `json:"fiscalYearBS"`
	Companyid           int       `json:"companyId"`
	Status              string    `json:"status"`
}

type CompanyRights struct {
	ID                  int       `json:"id"`
	Companyid           int       `json:"companyId"`
	Sn                  int       `json:"sn"`
	Rightshareid        int       `json:"rightShareId"`
	Companyname         string    `json:"companyName"`
	Stocksymbol         string    `json:"stockSymbol"`
	Issuemanager        string    `json:"issueManager"`
	Sectorname          string    `json:"sectorName"`
	Filename            string    `json:"fileName"`
	Ratio               string    `json:"ratio"`
	Priceperunit        string    `json:"pricePerUnit"`
	Units               string    `json:"units"`
	Openingdatead       string    `json:"openingDateAD"`
	Openingdatebs       string    `json:"openingDateBS"`
	Closingdatead       string    `json:"closingDateAD"`
	Closingdatebs       string    `json:"closingDateBS"`
	Extendeddatead      string    `json:"extendedDateAD"`
	Extendeddatebs      string    `json:"extendedDateBS"`
	Bookclosuredatead   string    `json:"bookClosureDateAD"`
	Bookclosuredatetime time.Time `json:"bookClosureDateTime"`
	Bookclosuredatebs   string    `json:"bookClosureDateBS"`
	Fiscalyearad        string    `json:"fiscalYearAD"`
	Fiscalyearbs        string    `json:"fiscalYearBS"`
	Isactive            bool      `json:"isActive"`
	Status              string    `json:"status"`
}

func (s *Server) GetDividends(w http.ResponseWriter, r *http.Request) {
	var div Dividends
	var rights Rights
	stock := r.URL.Query().Get("stock")
	if stock == "" {
		responses.ERROR(w, http.StatusBadRequest, errors.New("stock is required"))
		return
	}
	quantity := r.URL.Query().Get("quantity")
	if quantity == "" {
		responses.ERROR(w, http.StatusBadRequest, errors.New("quantity is required"))
		return
	}

	dividendMap := make(map[string][]CompanyDividend)
	rightsMap := make(map[string][]CompanyRights)
	dividendFile, err := ioutil.ReadFile("dividend.json")
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}
	rightFile, err := ioutil.ReadFile("right.json")
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}

	err = json.Unmarshal([]byte(dividendFile), &div)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}
	err = json.Unmarshal([]byte(rightFile), &rights)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}
	for _, v := range div.Data {
		dividendMap[v.Stocksymbol] = append(dividendMap[v.Stocksymbol], v)
	}
	for _, v := range rights.Data {
		rightsMap[v.Stocksymbol] = append(rightsMap[v.Stocksymbol], v)
	}

	divSorted, rightSorted, err := SortData(dividendMap[stock], rightsMap[stock])
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}

	re, err := CalculateReturn(quantity, divSorted, rightSorted)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}

	responses.JSON(w, http.StatusOK, map[string]interface{}{
		"dividends": divSorted,
		"rights":    rightSorted,
		"return":    re,
	})

}

func CalculateReturn(quantity string, dividend []CompanyDividend, right []CompanyRights) (DividendYield, error) {
	var dividendYield DividendYield
	if len(dividend) == 0 && len(right) == 0 {
		fmt.Println("Noting to return!!")
		return dividendYield, nil
	}

	if len(dividend) == 0 {
		fmt.Println("No dividend to return")
	}

	if len(right) == 0 {
		fmt.Println("No right to return")
	}

	quantityInt, err := strconv.Atoi(quantity)
	if err != nil {
		fmt.Println("Error in converting quantity to int")
	}

	dividendYield.TotalQuantity = float64(quantityInt)
	dividendYield.InitialQuantity = float64(quantityInt)

	for _, v := range dividend {
		cashInt, err := strconv.ParseFloat(v.Cash, 64)
		if err != nil {
			fmt.Println("Error in converting cash to int")
			return dividendYield, err
		}
		bonusInt, err := strconv.ParseFloat(v.Bonus, 64)
		if err != nil {
			fmt.Println("Error in converting bonus to int")
			return dividendYield, err
		}

		dividendYield.YearlyYield = append(dividendYield.YearlyYield, YearlyDividendYield{
			Year:        v.Fiscalyearad,
			Cash:        fmt.Sprintf("%.2f%%", cashInt),
			Bonus:       fmt.Sprintf("%.2f%%", bonusInt),
			CashBefore:  utils.ToFixed(dividendYield.TotalCashYield, 2),
			CashAfter:   utils.ToFixed(dividendYield.TotalCashYield+cashInt*dividendYield.TotalQuantity, 2),
			ShareBefore: math.Floor(dividendYield.TotalQuantity),
			ShareAfter:  math.Floor(dividendYield.TotalQuantity + (bonusInt/100)*dividendYield.TotalQuantity),
		})
		dividendYield.TotalCashYield = utils.ToFixed(dividendYield.TotalCashYield+cashInt*dividendYield.TotalQuantity, 2)
		dividendYield.TotalQuantity = math.Floor(dividendYield.TotalQuantity + (bonusInt/100)*dividendYield.TotalQuantity)
	}
	return dividendYield, nil
}

func caelculateDividendYield(dividend *bizmandu.DividendHistory, year int, quanity int) DividendYield {
	var dividendYield DividendYield
	dividendYield.TotalQuantity = float64(quanity)
	dividendYield.InitialQuantity = float64(quanity)
	for _, div := range dividend.Message.Dividend {
		if div.Bonus != 0 || div.Cash != 0 {
			divYear := strings.Split(div.Year, "/")[0]
			divYearInt, _ := strconv.Atoi(divYear)

			if divYearInt > year {

				dividendYield.YearlyYield = append(dividendYield.YearlyYield, YearlyDividendYield{
					Year:        divYearInt,
					Cash:        fmt.Sprintf("%.2f%%", div.Cash*100),
					Bonus:       fmt.Sprintf("%.2f%%", div.Bonus*100),
					CashBefore:  utils.ToFixed(dividendYield.TotalCashYield, 2),
					CashAfter:   utils.ToFixed(dividendYield.TotalCashYield+div.Cash*100*dividendYield.TotalQuantity, 2),
					ShareBefore: math.Floor(dividendYield.TotalQuantity),
					ShareAfter:  math.Floor(dividendYield.TotalQuantity + div.Bonus*dividendYield.TotalQuantity),
				})
				dividendYield.TotalCashYield = utils.ToFixed(dividendYield.TotalCashYield+div.Cash*100*dividendYield.TotalQuantity, 2)
				dividendYield.TotalQuantity = math.Floor(dividendYield.TotalQuantity + div.Bonus*dividendYield.TotalQuantity)
			}
		}

	}
	return dividendYield
}

func SortData(dividend []CompanyDividend, right []CompanyRights) ([]CompanyDividend, []CompanyRights, error) {
	var dateDiv []CompanyDividend
	var dateRight []CompanyRights

	for _, v := range dividend {
		date, err := utils.StringToTime(v.Bookclosuredatead)
		if err != nil {
			return nil, nil, err
		}
		v.Bookclosuredatetime = date
		dateDiv = append(dateDiv, v)
	}

	for _, v := range right {
		date, err := utils.StringToTime(v.Bookclosuredatead)
		if err != nil {
			return nil, nil, err
		}
		v.Bookclosuredatetime = date
		dateRight = append(dateRight, v)
	}

	sort.Slice(dateDiv, func(i, j int) bool {
		return dateDiv[i].Bookclosuredatetime.Before(dateDiv[j].Bookclosuredatetime)
	})

	sort.Slice(dateRight, func(i, j int) bool {
		return dateRight[i].Bookclosuredatetime.Before(dateRight[j].Bookclosuredatetime)
	})

	return dateDiv, dateRight, nil
}

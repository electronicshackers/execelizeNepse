package controllers

import (
	"encoding/json"
	"io/ioutil"
	responses "nepse-backend/api/response"
	"net/http"
)

type Dividends struct {
	Data []CompanyDividend `json:"data"`
}

type CompanyDividend struct {
	Sn                 int    `json:"sn"`
	Dividendhistoryid  int    `json:"dividendHistoryId"`
	Companyname        string `json:"companyName"`
	Stocksymbol        string `json:"stockSymbol"`
	Sectorname         string `json:"sectorName"`
	Filename           string `json:"fileName"`
	Bonus              string `json:"bonus"`
	Cash               string `json:"cash"`
	Totaldividend      string `json:"totalDividend"`
	Distributiondatead string `json:"distributionDateAD"`
	Distributiondatebs string `json:"distributionDateBS"`
	Bookclosuredatead  string `json:"bookClosureDateAD"`
	Bookclosuredatebs  string `json:"bookClosureDateBS"`
	Fiscalyearad       string `json:"fiscalYearAD"`
	Fiscalyearbs       string `json:"fiscalYearBS"`
	Companyid          int    `json:"companyId"`
	Status             string `json:"status"`
}

func (s *Server) GetDividends(w http.ResponseWriter, r *http.Request) {
	var div Dividends
	dividendMap := make(map[string][]CompanyDividend)
	dividendFile, err := ioutil.ReadFile("dividend.json")
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}

	err = json.Unmarshal([]byte(dividendFile), &div)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}
	for _, v := range div.Data {
		dividendMap[v.Stocksymbol] = append(dividendMap[v.Stocksymbol], v)
	}
	responses.JSON(w, http.StatusOK, dividendMap)
}

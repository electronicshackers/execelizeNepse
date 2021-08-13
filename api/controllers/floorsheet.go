package controllers

import (
	"fmt"
	responses "nepse-backend/api/response"
	"nepse-backend/nepse"
	"nepse-backend/nepse/neweb"
	"nepse-backend/utils"
	"net/http"
)

type FloorsheetResult struct {
	Ticker                string             `json:"ticker"`
	BuyerQuantityMap      map[string]int     `json:"buyerQuantityMap"`
	BuyerTurnOverMap      map[string]float64 `json:"buyerTurnOverMap"`
	BuyerAveragePriceMap  map[string]float64 `json:"buyerAveragePriceMap"`
	SellerQuantityMap     map[string]int     `json:"sellerQuantityMap"`
	SellerTurnOverMap     map[string]float64 `json:"sellerTurnOverMap"`
	SellerAveragePriceMap map[string]float64 `json:"sellerAveragePriceMap"`
}

func (s *Server) GetFloorsheet(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	ticker := params.Get("ticker")
	start := params.Get("start")
	end := params.Get("end")
	randomId := params.Get("id")

	days, err := utils.GetDateRange(w, start, end)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	fmt.Println("days", days)

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

	var id string

	for _, stock := range nepseSectors {
		if ticker == stock.Ticker {
			id = stock.Id
		}
	}

	var result FloorsheetResult

	result.BuyerQuantityMap = make(map[string]int)
	result.SellerQuantityMap = make(map[string]int)
	result.BuyerTurnOverMap = make(map[string]float64)
	result.SellerTurnOverMap = make(map[string]float64)
	result.BuyerAveragePriceMap = make(map[string]float64)
	result.SellerAveragePriceMap = make(map[string]float64)

	for _, day := range days {

		floorsheetInfo, err := nepseBeta.GetFloorsheet(id, day, randomId, 0)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		totalElement := floorsheetInfo.Totaltrades

		floorsheetInfoAgg, err := nepseBeta.GetFloorsheet(id, day, randomId, totalElement)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("err", err)
			return
		}

		for _, sheetData := range floorsheetInfoAgg.Floorsheets.Content {
			if sheetData.Buyermemberid != "" {
				result.BuyerQuantityMap[sheetData.Buyermemberid] += sheetData.Contractquantity
			}
			if sheetData.Sellermemberid != "" {
				result.SellerQuantityMap[sheetData.Sellermemberid] += sheetData.Contractquantity
			}
			if sheetData.Buyermemberid != "" {
				result.BuyerTurnOverMap[sheetData.Buyermemberid] += sheetData.Contractamount
			}
			if sheetData.Sellermemberid != "" {
				result.SellerTurnOverMap[sheetData.Sellermemberid] += sheetData.Contractamount
			}
		}

		for k, v := range result.BuyerQuantityMap {
			result.BuyerAveragePriceMap[k] = utils.ToFixed(float64(result.BuyerTurnOverMap[k])/float64(v), 2)
		}
		for k, v := range result.SellerQuantityMap {
			result.SellerAveragePriceMap[k] = utils.ToFixed(float64(result.SellerTurnOverMap[k])/float64(v), 2)
		}
	}
	result.Ticker = ticker

	responses.JSON(w, http.StatusOK, result)
}

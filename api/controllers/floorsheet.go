package controllers

import (
	"fmt"
	responses "nepse-backend/api/response"
	"nepse-backend/nepse"
	"nepse-backend/nepse/neweb"
	"net/http"
)

func (s *Server) GetFloorsheet(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	start := r.URL.Query().Get("start")

	randomId := r.URL.Query().Get("id")

	fmt.Println("start", start, "randomId", randomId)
	var nepseBeta nepse.NepseInterface

	nepseBeta, err := neweb.Neweb()

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

	floorsheetInfo, err := nepseBeta.GetFloorsheet(id, start, randomId, 0)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalElement := floorsheetInfo.Totaltrades

	floorsheetInfoAgg, err := nepseBeta.GetFloorsheet(id, start, randomId, totalElement)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses.JSON(w, http.StatusOK, floorsheetInfoAgg)
}

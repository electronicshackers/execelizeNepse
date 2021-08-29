package controllers

import (
	responses "nepse-backend/api/response"
	"nepse-backend/nepse"
	"nepse-backend/nepse/neweb"
	"net/http"
)

func (s *Server) GetStocks(w http.ResponseWriter, r *http.Request) {
	var nepseBeta nepse.NepseInterface

	nepseBeta, err := neweb.Neweb()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stocks, err := nepseBeta.GetStocks()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responses.JSON(w, http.StatusOK, stocks)
}

package controllers

import (
	"fmt"
	responses "nepse-backend/api/response"
	smartwealth "nepse-backend/nepse/smartWealth"
	"net/http"
)

func (server *Server) GetTechnicalData(w http.ResponseWriter, r *http.Request) {
	smart, err := smartwealth.SmartWealth()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	technical, err := smart.GetTechnicalAnalysis(7)

	fmt.Println("technical", technical)

	if err != nil {
		responses.ERROR(w, 400, err)
	}
	responses.JSON(w, 200, technical)

}

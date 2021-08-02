package controllers

import "net/http"

func (server *Server) GetMutualFundsInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mutualFunds":[]}`))
}

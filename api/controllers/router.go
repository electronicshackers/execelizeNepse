package controllers

import (
	"log"
	"nepse-backend/api/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
}

func (server *Server) setJSON(path string, next func(http.ResponseWriter, *http.Request), method string) {
	server.Router.Use(middlewares.CORS)
	server.Router.HandleFunc(path, middlewares.SetMiddlewareJSON(next)).Methods(method, "OPTIONS")
}

func (server *Server) InitRoutes() {
	server.setJSON("/api/v1/health", server.Health, "GET")
	server.setJSON("/api/v1/pricehistory", server.GetPriceHistory, "GET")
	server.setJSON("/api/v1/fundamental", server.GetFundamentalSectorwise, "GET")
	server.setJSON("/api/v1/whatif", server.WhatIf, "POST")
	server.setJSON("/api/v1/mutualfund", server.GetMutualFundsInfo, "GET")
	server.setJSON("/api/v1/floorsheet", server.GetFloorsheet, "GET")

	server.setJSON("/api/v1/floorsheet/bulk", server.GetFloorSheetAggregated, "GET")
	server.setJSON("/api/v1/floorsheet/analysis", server.FloorsheetAnalysis, "GET")

	server.setJSON("/api/v1/technical", server.GetTechnicalData, "GET")
	server.setJSON("/api/v1/stocks", server.GetStocks, "GET")
	server.setJSON("/api/v1/dividend", server.GetDividends, "GET")
}

func (server *Server) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

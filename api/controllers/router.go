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
	server.Router.HandleFunc(path, middlewares.SetMiddlewareJSON(next)).Methods(method, "OPTIONS")
}

func (server *Server) InitRoutes() {
	server.setJSON("/api/v1/health", server.Health, "GET")
	server.setJSON("/api/v1/pricehistory", server.GetPriceHistory, "GET")
	server.setJSON("/api/v1/fundamental", server.GetFundamentalSectorwise, "GET")
}

func (server *Server) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

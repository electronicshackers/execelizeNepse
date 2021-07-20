package api

import (
	"fmt"
	"nepse-backend/api/controllers"

	"github.com/gorilla/mux"
)

func New() *controllers.Server {
	s := &controllers.Server{}
	s.Router = mux.NewRouter()
	return s
}

func Run() {
	s := New()
	fmt.Println("Server is running on Port 8080")
	s.InitRoutes()
	s.Run(":8080")
}

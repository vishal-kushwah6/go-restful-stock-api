package routes

import (
	"go-postgress/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/stocks/{id}", middleware.GetAstock).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/stock", middleware.GetAllstocks).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/CreateStock", middleware.Createstock).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/stocks/{id}", middleware.Updatestocks).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deletestocks/{id}", middleware.Deletestocks).Methods("DELETE", "OPTIONS")

	return router

}

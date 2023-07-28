package routes

import (
	"github.com/gorilla/mux"
	"github.com/thiago1623/go-microservice/api/handlers"
)

// SetupRoutes Configure microservice routes.
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	router.HandleFunc("/consumption", handlers.GetConsumptionHandler).Methods("GET")

	return router
}

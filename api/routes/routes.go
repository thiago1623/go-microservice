package routes

import (
	"github.com/gorilla/mux"
	"github.com/thiago1623/go-microservice/api/handlers"
)

// SetupRoutes configura las rutas para el microservicio.
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", handlers.HomeHandler).Methods("GET")

	return router
}
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gitlab.com/baugoncalves/goclass-rest-api/middlewares"
	"gitlab.com/baugoncalves/goclass-rest-api/routes"
)

func majorRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func setRoutes(router *mux.Router) {
	router.HandleFunc("/", majorRoute)
	router.HandleFunc("/games", routes.GetGames).Methods("GET")
	router.HandleFunc("/games/{id}", routes.GetGameByID).Methods("GET")
	router.HandleFunc("/games", routes.NewGame).Methods("POST")
	router.HandleFunc("/games/{id}", routes.UpdateGame).Methods("PUT")
	router.HandleFunc("/games/{id}", routes.DeleteGame).Methods("DELETE")
}

func main() {
	var router *mux.Router

	log.Printf("Server is working here today on http://localhost:1616")

	router = mux.NewRouter().StrictSlash(true)

	router.Use(middlewares.JsonMiddleware)

	setRoutes(router)

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost", "http://localhost:8080"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		// AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	}).Handler(router)

	err := http.ListenAndServe(":1616", handler)

	if err != nil {
		fmt.Println("Error", err)
	}
}

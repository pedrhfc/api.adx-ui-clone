package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pedro9128/api_adx-ui-clone/middlewares"
	"github.com/pedro9128/api_adx-ui-clone/routes/units"
	"github.com/pedro9128/api_adx-ui-clone/routes/users"
	"github.com/rs/cors"
)

func majorRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func setRoutes(router *mux.Router) {
	router.HandleFunc("/", majorRoute)

	router.HandleFunc("/session", users.Session).Methods("POST")

	router.HandleFunc("/users", users.Index).Methods("GET")
	router.HandleFunc("/users/{id}", users.Show).Methods("GET")
	router.HandleFunc("/users", users.Store).Methods("POST")
	router.HandleFunc("/users/{id}", users.Update).Methods("PUT")
	router.HandleFunc("/users/{id}", users.Destroy).Methods("DELETE")

	router.HandleFunc("/units", units.Index).Methods("GET")
	router.HandleFunc("/units/{id}", units.Show).Methods("GET")
	router.HandleFunc("/units", units.Store).Methods("POST")
	router.HandleFunc("/units/{id}", units.Update).Methods("PUT")
	router.HandleFunc("/units/{id}", units.Destroy).Methods("DELETE")

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
		Debug:          true,
	}).Handler(router)

	err := http.ListenAndServe(":1616", handler)

	if err != nil {
		fmt.Println("Error", err)
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/baugoncalves/goclass-rest-api/middlewares"
	"gitlab.com/baugoncalves/goclass-rest-api/routes"
)

func majorRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func setRoutes(router *mux.Router) {
	router.HandleFunc("/", majorRoute)
	router.HandleFunc("/games", routes.GetGames)
	router.HandleFunc("/games/{gameID}", routes.GetGameByID)
	router.HandleFunc("/newgame", routes.NewGame)
	router.HandleFunc("/lookForGames", routes.LookForGame)
	router.HandleFunc("/updateGame", routes.UpdateGame)
}

func main() {
	var router *mux.Router

	log.Printf("Server is working here agai today on http://localhost:1602")

	router = mux.NewRouter()

	router.Use(middlewares.JsonMiddleware)

	setRoutes(router)

	err := http.ListenAndServe(":1616", router)
	if err != nil {
		fmt.Println("Error", err)
	}
}

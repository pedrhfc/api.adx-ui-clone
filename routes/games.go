package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/baugoncalves/goclass-rest-api/database"
)

type Game struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func GetGames(w http.ResponseWriter, r *http.Request) {
	var games []Game

	conn := database.SetConnection()
	defer conn.Close()

	selDB, err := conn.Query("SELECT * FROM games")

	if err != nil {
		fmt.Println("Error to fetch", err)
	}

	for selDB.Next() {
		var game Game

		err = selDB.Scan(&game.ID, &game.Name, &game.Price)
		games = append(games, game)
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(games)
}

func NewGame(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	conn := database.SetConnection()
	stmt, err := conn.Prepare("INSERT INTO games(name,price) VALUES(?,?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	name := keyVal["name"]
	price := keyVal["price"]

	_, err = stmt.Exec(name, price)
	if err != nil {
		panic(err.Error())
	}
}

func GetGameByID(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	conn := database.SetConnection()
	defer conn.Close()

	vars := mux.Vars(r)
	id := vars["id"]
	game := Game{}
	encoder := json.NewEncoder(w)

	err := conn.QueryRow("SELECT * FROM games WHERE id=?", id).Scan(&game.ID, &game.Name, &game.Price)

	if err != nil {
		fmt.Println("Error to fetch", err)
	}

	encoder.Encode(game)
}

func UpdateGame(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	conn := database.SetConnection()

	stmt, err := conn.Prepare("UPDATE games SET name=?, price=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	var register Game

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, "Bad Request")
	}
	json.Unmarshal(body, &register)

	_, err = stmt.Exec(register.Name, register.Price, register.ID)
	if err != nil {
		panic(err.Error())
	}
}

func DeleteGame(w http.ResponseWriter, r *http.Request) {
	conn := database.SetConnection()

	params := mux.Vars(r)
	stmt, err := conn.Prepare("DELETE FROM games WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
}

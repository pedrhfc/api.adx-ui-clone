package units

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pedro9128/api_adx-ui-clone/database"
)

type Unit struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	var units []Unit

	conn := database.SetConnection()
	defer conn.Close()

	selDB, err := conn.Query("SELECT * FROM units WHERE id<>1")

	if err != nil {
		fmt.Println("Error to fetch", err)
	}

	for selDB.Next() {
		var unit Unit

		err = selDB.Scan(&unit.ID, &unit.Name)
		units = append(units, unit)
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(units)
}

func Store(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	conn := database.SetConnection()
	stmt, err := conn.Prepare("INSERT INTO units(name,price) VALUES(?,?)")
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

func Show(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	conn := database.SetConnection()
	defer conn.Close()

	vars := mux.Vars(r)
	id := vars["id"]
	unit := Unit{}
	encoder := json.NewEncoder(w)

	err := conn.QueryRow("SELECT * FROM units WHERE id=?", id).Scan(&unit.ID, &unit.Name)

	if err != nil {
		fmt.Println("Error to fetch", err)
	}

	encoder.Encode(unit)
}

func Update(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	conn := database.SetConnection()

	stmt, err := conn.Prepare("UPDATE units SET name=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	var register Unit

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, "Bad Request")
	}
	json.Unmarshal(body, &register)

	_, err = stmt.Exec(register.Name, register.ID)
	if err != nil {
		panic(err.Error())
	}
}

func Destroy(w http.ResponseWriter, r *http.Request) {
	conn := database.SetConnection()

	params := mux.Vars(r)
	stmt, err := conn.Prepare("DELETE FROM units WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
}

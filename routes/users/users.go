package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pedro9128/api_adx-ui-clone/database"
)

type User struct {
	ID           int    `json:"id,omitempty"`
	UnitID       int    `json:"unit_id,omitempty"`
	Name         string `json:"name,omitempty"`
	Unit         string `json:"unit,omitempty"`
	Email        string `json:"email,omitempty"`
	Password     string `json:"password,omitempty"`
	Registration string `json:"registration,omitempty"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	var users []User

	conn := database.SetConnection()
	defer conn.Close()
	selDB, err := conn.Query("SELECT users.id,users.name,users.email,users.registration,units.name as unit FROM users INNER JOIN users_has_units ON users_has_units.user_id=users.id INNER JOIN units ON units.id=users_has_units.unit_id")

	if err != nil {
		fmt.Println("Error to fetch", err)
	}

	for selDB.Next() {
		var user User
		err = selDB.Scan(&user.ID, &user.Name, &user.Email, &user.Registration, &user.Unit)
		users = append(users, user)
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(users)
}

func setUnit(name string,
	email string,
	password string,
	registration string) {
	conn := database.SetConnection()
	user := User{}
	_err := conn.QueryRow("SELECT users.id FROM users WHERE name=? AND email=? AND password=? AND registration=?", name, email, password, registration).Scan(&user.ID)

	if _err != nil {
		fmt.Println("Error to fetch", _err)
	}

	stmt, err := conn.Prepare("INSERT INTO users_has_units(user_id,unit_id) VALUES(?,?)")

	_, err = stmt.Exec(user.ID, 1)
	if err != nil {
		panic(err.Error())
	}

}

func Store(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	conn := database.SetConnection()
	stmt, err := conn.Prepare("INSERT INTO users(name,email,password,registration) VALUES(?,?,?,?)")

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
	email := keyVal["email"]
	password := keyVal["password"]
	registration := keyVal["registration"]

	_, err = stmt.Exec(name, email, password, registration)
	if err != nil {
		panic(err.Error())
	}

	setUnit(name, email, password, registration)
}

func Show(w http.ResponseWriter, r *http.Request) {
	conn := database.SetConnection()
	defer conn.Close()

	vars := mux.Vars(r)
	id := vars["id"]
	user := User{}

	err := conn.QueryRow("SELECT users.id,users.name,users.email, users.registration, units.id as unit_id FROM users INNER JOIN users_has_units ON users_has_units.user_id = users.id INNER JOIN units ON units.id = users_has_units.unit_id WHERE users.id=?", id).Scan(&user.ID, &user.Name, &user.Email, &user.Registration, &user.UnitID)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println("Error to fetch", err)
	}

	defer r.Body.Close()

	json.NewEncoder(w).Encode(user)
}

func updateUserUnit(unitID int, userID int) {
	conn := database.SetConnection()

	stmt, err := conn.Prepare("UPDATE users_has_units SET unit_id=? WHERE user_id=?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(unitID, userID)
	if err != nil {
		panic(err.Error())
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	conn := database.SetConnection()

	stmt, err := conn.Prepare("UPDATE users SET name=?, email=?, registration=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	var register User

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, "Bad Request")
	}
	json.Unmarshal(body, &register)

	_, err = stmt.Exec(register.Name, register.Email, register.Registration, register.ID)
	if err != nil {
		panic(err.Error())
	}

	updateUserUnit(register.UnitID, register.ID)
}

func Destroy(w http.ResponseWriter, r *http.Request) {
	conn := database.SetConnection()

	params := mux.Vars(r)
	stmt, err := conn.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
}

func Session(w http.ResponseWriter, r *http.Request) {
	conn := database.SetConnection()

	var tempUser User
	var user User

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&tempUser)

	err := conn.QueryRow("SELECT users.id,users.name,users.email, users.registration, units.name FROM users INNER JOIN users_has_units ON users_has_units.user_id = users.id INNER JOIN units ON units.id = users_has_units.unit_id WHERE password=? AND registration=? AND units.id=?", tempUser.Password, tempUser.Registration, tempUser.UnitID).Scan(&user.ID, &user.Name, &user.Email, &user.Registration, &user.Unit)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println("Error to fetch", err)
	}

	defer r.Body.Close()

	json.NewEncoder(w).Encode(user)

}

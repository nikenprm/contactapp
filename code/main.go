package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"net/http"


	"github.com/gorilla/mux"
	"encoding/json"
	"log"
	"strconv"
	"./contact"
//	"errors"
)

type ErrorMessage struct {
    Err string
    Code int
}

var db *sql.DB


func GetAllContacts(w http.ResponseWriter, req *http.Request) {


	rows, err := db.Query("SELECT * FROM user")
	checkErr(err)

	defer rows.Close()

	for rows.Next() {
		var id string
		var name string
		var phoneNum string
		var address string
		err = rows.Scan(&id, &name, &phoneNum, &address)
		checkErr(err)
		contact := &contact.Contact{id, name, phoneNum, address}
		json.NewEncoder(w).Encode(contact)
	}


}

func GetOneContact(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	p := params["id"]

	id, _ := strconv.Atoi(p)

	maxID := getMaxID()


	if id<=maxID {

		rows, err := db.Query("SELECT * FROM user WHERE id = ?", id)
		checkErr(err)

		defer rows.Close()

		for rows.Next() {
			var id string
			var name string
			var phoneNum string
			var address string

			err = rows.Scan(&id, &name, &phoneNum, &address)
			checkErr(err)
			contact := &contact.Contact{id, name, phoneNum, address}
			json.NewEncoder(w).Encode(contact)

		}
	} else {

		sendErrorMessage(w, "There is no user with that ID")
	}


}

func CreateContact(w http.ResponseWriter, req *http.Request) {
	name := req.PostFormValue("name")
	phoneNum := req.PostFormValue("phoneNum")
	address:= req.PostFormValue("address")


	if _, err := strconv.Atoi(name); err != nil {
		  sendErrorMessage(w, "Cannot contain number")
	} else {

		query, err := db.Prepare("INSERT INTO user(name, phoneNum, address) VALUES(?,?,?)")
		checkErr(err)

		_, error := query.Exec(name, phoneNum, address)
		checkErr(error)

	}
	
}

func EditContact(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id := params["id"]
	name := req.PostFormValue("name")
	phoneNum := req.PostFormValue("phoneNum")
	address:= req.PostFormValue("address")

	query, err := db.Prepare("UPDATE user SET name=?, phoneNum=?, address=? WHERE id=?")
	checkErr(err)                                                                      
	                                                                                   
	res, err := query.Exec(name,phoneNum,address,id)
	checkErr(err)                                                                      
	                                                                                   
	if res != nil {                                                                    
		fmt.Println("Success")                                                     
	}                                                                                  


}

func DeleteContact(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	p := params["id"]

	id, _ := strconv.Atoi(p)
	maxID := getMaxID()

	fmt.Println("Max ID: ", maxID)
	fmt.Println("ID: ", id)

	if id<=maxID {
		query, err := db.Prepare("DELETE FROM user WHERE id=?")
		checkErr(err)

		_, error := query.Exec(id)
		checkErr(error)


	} else {

		//errortext := ErrorMessage{"There is no user with that ID", 400}
		//json.NewEncoder(w).Encode(errortext)
		sendErrorMessage(w, "There is no user with that ID")
		//w.WriteHeader(http.StatusBadRequest)
	}

}


func main() {

	var err error

	db, err = sql.Open("mysql", "root:123456@/contactDB")
	checkErr(err)

	router := mux.NewRouter()
	router.HandleFunc("/contacts", GetAllContacts).Methods("GET")
	router.HandleFunc("/contacts/{id}", GetOneContact).Methods("GET")
	router.HandleFunc("/contacts", CreateContact).Methods("POST")
	router.HandleFunc("/contacts/{id}", EditContact).Methods("PUT")
	router.HandleFunc("/contacts/{id}", DeleteContact).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":12345", router))

	defer db.Close()


}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}


func getMaxID() int {
	var maxID int

	rows, err := db.Query("SELECT MAX(id) FROM user")
	checkErr(err)

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&maxID)
		checkErr(err)
	}

	return maxID

}

func sendErrorMessage(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusBadRequest)
}

func (e ErrorMessage) sendErrorMessage(w http.ResponseWriter) {
	http.Error(w, e.Err, e.Code)
}
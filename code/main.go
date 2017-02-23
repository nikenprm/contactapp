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
//	"errors"
)

type Contact struct {
	Id string `json:"id"`
	Name string `json:"name,"`
	PhoneNum string `json:"phoneNum"`
	Address string `json:"address"`
}

type ErrorMessage struct {
    Err string
    Code int
}

func GetAllContacts(w http.ResponseWriter, req *http.Request) {
	db, err := sql.Open("mysql", "root:123456@/contactDB")
	checkErr(err)

	rows, err := db.Query("SELECT * FROM user")
	checkErr(err)

	for rows.Next() {
		var id string
		var name string
		var phoneNum string
		var address string
		err = rows.Scan(&id, &name, &phoneNum, &address)
		checkErr(err)
		contact := &Contact{id, name, phoneNum, address}
		json.NewEncoder(w).Encode(contact)
	}

	db.Close()

}

func GetOneContact(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	db, err := sql.Open("mysql", "root:123456@/contactDB")
	checkErr(err)

	rows, err := db.Query("SELECT * FROM user")
	checkErr(err)

	for rows.Next() {
		var id string
		var name string
		var phoneNum string
		var address string

		err = rows.Scan(&id, &name, &phoneNum, &address)
		checkErr(err)
		if id == params["id"]   {

			 contact := &Contact{id, name, phoneNum, address}
		         json.NewEncoder(w).Encode(contact)
		}
	}

	db.Close()
}


func CreateContact(w http.ResponseWriter, req *http.Request) {
	name := req.PostFormValue("name")
	phoneNum := req.PostFormValue("phoneNum")
	address:= req.PostFormValue("address")

	//fmt.Println(name)
	//fmt.Println(phoneNum)
	//fmt.Println(address)

	db, err := sql.Open("mysql", "root:123456@/contactDB")
	checkErr(err)

	if _, err := strconv.Atoi(name); err != nil {
		  sendErrorMessage(w, "Cannot contain number")
	} else {

		query, err := db.Prepare("INSERT INTO user(name, phoneNum, address) VALUES(?,?,?)")
		checkErr(err)

		_, error := query.Exec(name, phoneNum, address)
		checkErr(error)

	}

	db.Close()
	
	
}

func EditContact(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id := params["id"]
	name := req.PostFormValue("name")
	phoneNum := req.PostFormValue("phoneNum")
	address:= req.PostFormValue("address")

	db, err := sql.Open("mysql", "root:123456@/contactDB")                             
	checkErr(err)                                                                      
	                                                                                   
	query, err := db.Prepare("UPDATE user SET name=?, phoneNum=?, address=? WHERE id=?")
	checkErr(err)                                                                      
	                                                                                   
	res, err := query.Exec(name,phoneNum,address,id)
	checkErr(err)                                                                      
	                                                                                   
	if res != nil {                                                                    
		fmt.Println("Success")                                                     
	}                                                                                  
	                                                                                   
	db.Close()                                                                         

}

func DeleteContact(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	p := params["id"]

	id, err := strconv.Atoi(p)
	var maxID int

	db, err := sql.Open("mysql", "root:123456@/contactDB")
	checkErr(err)

	row, err := db.Query("SELECT MAX(id) FROM user")
	checkErr(err)

	err = row.Scan(&maxID)

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

	db.Close()
}

func sendErrorMessage(w http.ResponseWriter, message string) {
	      http.Error(w, message, http.StatusBadRequest)
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/contacts", GetAllContacts).Methods("GET")
	router.HandleFunc("/contacts/{id}", GetOneContact).Methods("GET")
	router.HandleFunc("/contacts", CreateContact).Methods("POST")
	router.HandleFunc("/contacts/{id}", EditContact).Methods("PUT")
	router.HandleFunc("/contacts/{id}", DeleteContact).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":12345", router))


}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
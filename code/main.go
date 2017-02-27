package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"net/http"


	"github.com/gorilla/mux"
//	"encoding/json"
	"log"
//	"strconv"
//	"./contact"
	"./endpoints"
//	"errors"
)

type ErrorMessage struct {
    Err string
    Code int
}

var db *sql.DB
func main() {

	var err error

	db, err = sql.Open("mysql", "root:123456@/contactDB")
	checkErr(err)

	router := mux.NewRouter()
	router.HandleFunc("/contacts", endpoints.GetAllContacts).Methods("GET")
	router.HandleFunc("/contacts/{id}", endpoints.GetContactProfile).Methods("GET")
	router.HandleFunc("/contacts", endpoints.CreateNewContact).Methods("POST")
	router.HandleFunc("/contacts/{id}", endpoints.EditContact).Methods("PUT")
	router.HandleFunc("/contacts/{id}", endpoints.DeleteContact).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":12345", router))

	defer db.Close()


}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}



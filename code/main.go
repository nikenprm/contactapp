package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"net/http"


	"github.com/gorilla/mux"
	"encoding/json"
	"log"
)

type Contact struct {
	Id int `json:"id"`
	Name string `json:"name,"`
	PhoneNum string `json:"phoneNum"`
	Address string `json:"address"`
}

func GetAllContacts(w http.ResponseWriter, req *http.Request) {
	db, err := sql.Open("mysql", "root:123456@/contactDB")
	checkErr(err)

	rows, err := db.Query("SELECT * FROM user")
	checkErr(err)

	for rows.Next() {
		var id int
		var name string
		var phoneNum string
		var address string
		err = rows.Scan(&id, &name, &phoneNum, &address)
		checkErr(err)
		contact := &Contact{id, name, phoneNum, address})
		json.NewEncoder(w).Encode(contact)
	}

	db.Close()

}


func main() {

	router := mux.NewRouter()
	router.HandleFunc("/contacts", GetAllContacts).Methods("GET")
	log.Fatal(http.ListenAndServe(":12345", router))


}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
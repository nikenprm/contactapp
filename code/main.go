package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"log"
	"net/http"

	"./endpoints"

)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/contacts", endpoints.GetAllContacts).Methods("GET")
	router.HandleFunc("/contacts/{id}", endpoints.GetContactProfile).Methods("GET")
	router.HandleFunc("/contacts", endpoints.CreateNewContact).Methods("POST")
	router.HandleFunc("/contacts/{id}", endpoints.EditContact).Methods("PUT")
	router.HandleFunc("/contacts/{id}", endpoints.DeleteContact).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":12345", router))

}





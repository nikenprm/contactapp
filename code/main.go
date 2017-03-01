package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"log"
	"net/http"

	"./endpoints"
	"./contact"
	"./config"

	"fmt"
	"os"
	"flag"
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args)>0 {
		switch args[0] {
		case "serve":
			executeServer()
		default:
			fmt.Println("Invalid command")
			os.Exit(1)
		}
	} else {
		fmt.Println("Invalid command, please retry with the following format <exec> serve")
		os.Exit(1)
	}

}

func executeServer() {
	if err := config.LoadConfigFile("./config.json"); err != nil {
		fmt.Printf("Error: %s loading configuration file: %s\n", "./config.json", err)
		os.Exit(1)
	}

	contact.ConnectDB()
	router := mux.NewRouter()
	router.HandleFunc("/contacts", endpoints.GetAllContacts).Methods("GET")
	router.HandleFunc("/contacts/{id}", endpoints.GetContactProfile).Methods("GET")
	router.HandleFunc("/contacts", endpoints.CreateNewContact).Methods("POST")
	router.HandleFunc("/contacts/{id}", endpoints.EditContact).Methods("PUT")
	router.HandleFunc("/contacts/{id}", endpoints.DeleteContact).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":12345", router))

	fmt.Println("HTTP is listening on port 12345")

	defer contact.CloseDB()
}





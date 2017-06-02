package main

import (
	"github.com/buaazp/fasthttprouter"
	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"

	"github.com/contactapp/config"
	"github.com/contactapp/endpoints"
	"github.com/contactapp/repository"

	"flag"
	"fmt"
	"os"
)

var migrationPath = flag.String("m", "./migrations/postgres", "path to migration directory")

func main() {
	flag.Parse()

	if err := config.LoadConfigFile("./config.json"); err != nil {
		fmt.Printf("Error: %s loading configuration file: %s\n", "./config.json", err)
		os.Exit(1)
	}

	args := flag.Args()

	if len(args) > 0 {
		switch args[0] {
		case "migration":
			executeMigration(args)
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
	repository.ConnectDB()
	router := fasthttprouter.New()

	router.GET("/contacts", endpoints.GetAllContacts)
	router.GET("/contacts/:id", endpoints.GetContactProfile)
	router.GET("/download", endpoints.DownloadContactProfile)
	router.POST("/contacts", endpoints.CreateNewContact)
	router.PUT("/contacts/:id", endpoints.EditContact)
	router.DELETE("/contacts/:id", endpoints.DeleteContact)
	fmt.Println(fasthttp.ListenAndServe(":12345", router.Handler))

	defer repository.CloseDB()
}

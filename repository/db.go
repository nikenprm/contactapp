package repository

import (
	"database/sql"
	"fmt"

	"github.com/contactapp/config"
)

var db *sql.DB

func ConnectDB() {
	var err error
	url, err := config.Config.DB.ConnectionString()
	//db, err = sql.Open("mysql", "root:123456@/contactDB")
	db, err = sql.Open(config.Config.DB.Type, url)
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	checkErr(err)
	PrepareStatements(err)

}

func CloseDB() {
	fmt.Println("db is closed")

	db.Close()
}

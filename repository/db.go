package repository

import (
	"database/sql"
	"fmt"

	"github.com/contactapp/config"

	_ "github.com/mattes/migrate/driver/mysql"
	_ "github.com/mattes/migrate/driver/postgres"
	"github.com/mattes/migrate/migrate"
)

var db *sql.DB

func ConnectDB() {
	//var err error
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

func Migrate(migrationPath string, num int) (errors []error) {
	url, err := config.Config.DB.MigrationString()
	if err != nil {
		return append(errors, err)
	}
	if num == 0 {
		errors, _ = migrate.UpSync(url, migrationPath)
	} else {
		errors, _ = migrate.MigrateSync(url, migrationPath, num)
	}
	return
}

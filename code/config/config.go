package config


type DB struct {

	Type string `json:"type"`
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Name string `json:"name"`
	Pass string `json:"pass"`

}

var db *DB


func SetMySQLDB() {

	db.Type = "mysql"
	db.Host = "localhost"
	db.Port = "12345"
	db.User = "root"
	db.Name = "contactapp"
	db.Pass = "123456"

}

func (db *DB) setPostGresDB() {

	db.Type = "postgres"
	db.Host = "localhost"
	db.Port = "12345"
	db.User = "root"
	db.Name = "contactapp"
	db.Pass = "123456"

}
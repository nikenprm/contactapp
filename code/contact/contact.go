package contact

import (
	"database/sql"
	"fmt"
	"strconv"
	_ "github.com/lib/pq"
	"../config"

)

type Contact struct {
	Id string `json:"id, omitempty"`
	Name string `json:"name, omitempty"`
	PhoneNum string `json:"phoneNum, omitempty"`
	Address string `json:"address, omitempty"`
}

var db *sql.DB

var (
	getContactFromIDSQL = map[string]string{
		"mysql":    "SELECT * FROM user WHERE id = ?",
		"postgres": "SELECT * FROM contactinfo WHERE id= $1",
	}

	createContactSQL = map[string]string{
		"mysql": "INSERT INTO user(name, phoneNum, address) VALUES(?,?,?)",
		"postgres": "INSERT INTO contactinfo(name, phoneNum, address) VALUES($1,$2,$3)",
	}

	updateContactSQL = map[string]string{
		"mysql": "UPDATE user SET name=?, phoneNum=?, address=? WHERE id=?",
		"postgres": "UPDATE contactinfo SET name=$1, phoneNum=$2, address=$3 WHERE id=$4",
	}

	deleteContactSQL = map[string]string{
		"mysql": "DELETE FROM user WHERE id=?",
		"postgres": "DELETE FROM contactinfo WHERE id=$1",
	}

	getMaxIDSQL = map[string]string{
		"mysql": "SELECT MAX(id) FROM user",
		"postgres" : "SELECT MAX(id) FROM contactinfo",

	}
)

var (
	//database config.DB
	getContactFromIDStmt     *sql.Stmt
	createContactStmt 	 *sql.Stmt
	updateContactStmt        *sql.Stmt
	deleteContactStmt	 *sql.Stmt
)


func PrepareStatements(err error) {
	getContactFromIDStmt, err = db.Prepare(getContactFromIDSQL[config.Config.DB.Type])
	if err != nil {
		return
	}

	createContactStmt, err = db.Prepare(createContactSQL[config.Config.DB.Type])
	if err != nil {
		return
	}

	updateContactStmt, err = db.Prepare(updateContactSQL[config.Config.DB.Type])
	if err != nil {
		return
	}

	deleteContactStmt, err = db.Prepare(deleteContactSQL[config.Config.DB.Type])
	if err != nil {
		return
	}
}

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

func GetContactByID (userId int) Contact {

	var id string
	var name string
	var phoneNum string
	var address string

	maxID := getMaxID()
	fmt.Println(userId, maxID)

	if userId <= maxID {
		rows, err := getContactFromIDStmt.Query(userId)
		checkErr(err)

		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&id, &name, &phoneNum, &address)
			checkErr(err)
		}
	}
	contact := Contact{id, name, phoneNum, address}
	return contact
}

func GetAllContacts () []Contact {

	var contacts []Contact

	maxID := getMaxID()

	for i := 1;i <= maxID; i++ {
		contact := GetContactByID(i)
		contacts = append(contacts, contact)

	}

	return contacts

}

func CreateContact (name, phoneNum, address string) {

	_, error := createContactStmt.Exec(name, phoneNum, address)
	checkErr(error)
}

func UpdateContact (id, name, phoneNum, address string) {

	userID, err := strconv.Atoi(id)
	checkErr(err)

	contact := GetContactByID(userID)

	if name == "" {
		name = contact.Name
	}

	if phoneNum == "" {
		phoneNum = contact.PhoneNum
	}

	if address == "" {
		address = contact.Address
	}

	res, err := updateContactStmt.Exec(name,phoneNum,address,id)
	checkErr(err)

	if res != nil {
		fmt.Println("Success")
	}

}

func DeleteContact (p string) bool{

	id, _ := strconv.Atoi(p)
	maxID := getMaxID()

	if id<=maxID {
		_, error := deleteContactStmt.Exec(id)
		checkErr(error)
		return true
	} else {
		return false

	}

}

func checkErr(err error) {
	if err != nil {
		fmt.Println("error!!!!",err)
	}
}

func getMaxID() int {
	var maxID int

	err := db.QueryRow(getMaxIDSQL[config.Config.DB.Type]).Scan(&maxID)
	checkErr(err)

	return maxID
}
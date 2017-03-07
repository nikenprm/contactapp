package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	_ "github.com/lib/pq"
	"../config"

)

type Contact struct {
	Id string `json:"id"`
	Name string `json:"name"`
	PhoneNum string `json:"phoneNum"`
	Address string `json:"address"`
}

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


func GetContactByID (userId int) Contact {
	var id string
	var name string
	var phoneNum string
	var address string

	if isLessThanMaxID(userId) {
		err := getContactFromIDStmt.QueryRow(userId).Scan(&id, &name, &phoneNum, &address)
		checkErr(err)
	}
	contact := Contact{id, name, phoneNum, address}
	return contact
}

func GetAllContacts () []Contact {
	var contacts []Contact

	maxID := getMaxID()

	for i := 1;i <= maxID; i++ {
		contact := GetContactByID(i)

		//if we do not use this then the row that has been deleted will also appear even though
		//they have empty values
		if contact.Id!= "" {
			contacts = append(contacts, contact)
		}
	}
	return contacts
}

func CreateContact (name, phoneNum, address string) {
	_, error := createContactStmt.Exec(name, phoneNum, address)
	checkErr(error)
}

func UpdateContact (id, name, phoneNum, address string){
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

	_, err = updateContactStmt.Exec(name,phoneNum,address,id)
	checkErr(err)
}

func DeleteContact (p string) error{
	var err error
	id, _ := strconv.Atoi(p)

	contact:=GetContactByID(id)

	fmt.Println("Deleting contact ", contact.Name)

	if isLessThanMaxID(id){
		_, err = deleteContactStmt.Exec(id)
		checkErr(err)
		return nil
	} else {
		return err
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

func isLessThanMaxID(userID int) bool{
	maxID := getMaxID()

	return userID<=maxID
}
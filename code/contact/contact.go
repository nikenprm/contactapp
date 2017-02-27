package contact

import (
	"database/sql"
	"fmt"
	"strconv"
)

type Contact struct {
	Id string `json:"id"`
	Name string `json:"name,"`
	PhoneNum string `json:"phoneNum"`
	Address string `json:"address"`
}

var db *sql.DB

func connectDB() {
	var err error
	db, err = sql.Open("mysql", "root:123456@/contactDB")
	checkErr(err)

}

func GetContactByID (userId int) Contact {

	var id string
	var name string
	var phoneNum string
	var address string
	connectDB()

	maxID := getMaxID()
	fmt.Println(userId, maxID)

	if userId <= maxID {

		rows, err := db.Query("SELECT * FROM user WHERE id = ?", userId)
		checkErr(err)

		defer rows.Close()


		for rows.Next() {

			fmt.Println("im here")
			err = rows.Scan(&id, &name, &phoneNum, &address)
			checkErr(err)

		}

	}
	contact := Contact{id, name, phoneNum, address}
	fmt.Println("Model: ",contact.Id, contact.Name)
	return contact

}

func GetAllContacts () []Contact {

	var contacts []Contact

	connectDB()
	maxID := getMaxID()

	for i := 1;i <= maxID; i++ {
		contact := GetContactByID(i)
		contacts = append(contacts, contact)

	}

	return contacts

}

func CreateContact (name, phoneNum, address string) {
	connectDB()

	query, err := db.Prepare("INSERT INTO user(name, phoneNum, address) VALUES(?,?,?)")
	checkErr(err)

	_, error := query.Exec(name, phoneNum, address)
	checkErr(error)
}

func Update (id, name, phoneNum, address string) {
	connectDB()

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

	query, err := db.Prepare("UPDATE user SET name=?, phoneNum=?, address=? WHERE id=?")
	checkErr(err)

	res, err := query.Exec(name,phoneNum,address,id)
	checkErr(err)

	if res != nil {
		fmt.Println("Success")
	}

}

func Delete (p string) bool{
	connectDB()

	id, _ := strconv.Atoi(p)
	maxID := getMaxID()

	fmt.Println("Max ID: ", maxID)
	fmt.Println("ID: ", id)

	if id<=maxID {
		query, err := db.Prepare("DELETE FROM user WHERE id=?")
		checkErr(err)

		_, error := query.Exec(id)
		checkErr(error)
		return true
	} else {
		return false

	}


}



func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func getMaxID() int {
	var maxID int

	err := db.QueryRow("SELECT MAX(id) FROM user").Scan(&maxID)
	checkErr(err)

	return maxID

}
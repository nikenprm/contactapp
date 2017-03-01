package endpoints

import (

	"github.com/gorilla/mux"

	"net/http"
	"strconv"
	"encoding/json"
	"../contact"
	"fmt"
	"regexp"
)

func GetContactProfile(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	p := params["id"]

	id, _ := strconv.Atoi(p)

	fmt.Println(id)

	contact := contact.GetContactByID(id)
	fmt.Println("Endpoint: ",contact.Id, contact.Name)
	if contact.Id==""{
		sendErrorMessage(w, "There is no user with that ID")

	} else {


		json.NewEncoder(w).Encode(contact)
	}

}

func GetAllContacts(w http.ResponseWriter, req *http.Request) {
	var contacts []contact.Contact
	contacts = contact.GetAllContacts()

	for _,c := range contacts {
		json.NewEncoder(w).Encode(c)
	}
}

func CreateNewContact(w http.ResponseWriter, req *http.Request) {
	name := req.PostFormValue("name")
	phoneNum := req.PostFormValue("phoneNum")
	address:= req.PostFormValue("address")

	//decoder := json.NewDecoder(req.Body)

	var IsLetter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

	if IsLetter(name) {
		contact.CreateContact(name,phoneNum,address)

	} else {

		fmt.Println("name:",name)
		sendErrorMessage(w, "Cannot contain number")
	}
}

func EditContact(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id := params["id"]
	name := req.PostFormValue("name")
	phoneNum := req.PostFormValue("phoneNum")
	address:= req.PostFormValue("address")

	contact.UpdateContact(id,name,phoneNum,address)
}

func DeleteContact(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	p := params["id"]

	isValid := contact.DeleteContact(p)

	if isValid != true {

		sendErrorMessage(w, "There is no user with that ID")

	}

}


func sendErrorMessage(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusBadRequest)
}

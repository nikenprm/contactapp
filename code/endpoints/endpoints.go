package endpoints

import (

	"github.com/gorilla/mux"

	"net/http"
	"strconv"
	"encoding/json"
	"../contact"
	"fmt"
	"regexp"
	"strings"
)

type ContactCreationStruct struct {

	Id string `json:"id"`
	Name string `json:"name"`
	PhoneNum string `json:"phoneNum"`
	Address string `json:"address"`

}

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
	var contactStruct ContactCreationStruct

	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&contactStruct); err != nil {
		sendErrorMessage(w, "Error decoding the input")
		return
	}

	//used to check whether a string contains number or not, because value of name cannot contain number
	var IsLetter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

	//the function above will assume that whitespace is not letter
	//so we have to first strip all whitespace to correctly check the variable
	tempname := strings.Join(strings.Fields(contactStruct.Name),"")

	if IsLetter(tempname) {
		contact.CreateContact(contactStruct.Name,contactStruct.PhoneNum,contactStruct.Address)

	} else {

		fmt.Println("name:",contactStruct.Name)
		sendErrorMessage(w, "Name cannot contain number")
	}
}

func EditContact(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id := params["id"]

	var contactStruct ContactCreationStruct
	contactStruct.Id = id


	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&contactStruct); err != nil {
		sendErrorMessage(w, "Error decoding the input")
		return
	}

	var IsLetter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString
	tempname := strings.Join(strings.Fields(contactStruct.Name),"")

	if IsLetter(tempname) {
		contact.UpdateContact(contactStruct.Id, contactStruct.Name, contactStruct.PhoneNum, contactStruct.Address)
	} else {

		sendErrorMessage(w, "Name cannot contain number")
	}
}

func DeleteContact(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	p := params["id"]

	error := contact.DeleteContact(p)

	if error!=nil {

		sendErrorMessage(w, "There is no user with that ID")

	}

}


func sendErrorMessage(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusBadRequest)
}

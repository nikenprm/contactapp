package endpoints

import (
	"strconv"

	"github.com/valyala/fasthttp"

	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"bytes"

	"github.com/contactapp/repository"
)

type ContactCreationStruct struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	PhoneNum string `json:"phoneNum"`
	Address  string `json:"address"`
}

func GetContactProfile(ctx *fasthttp.RequestCtx) {
	param, _ := ctx.UserValue("id").(string)

	id, _ := strconv.Atoi(param)
	fmt.Println(id)

	contact := repository.GetContactByID(id)
	if contact.Id == "" {
		sendErrorMessage(ctx, "There is no user with that ID")
	} else {
		json.NewEncoder(ctx).Encode(contact)
	}
}

func GetAllContacts(ctx *fasthttp.RequestCtx) {
	var contacts []repository.Contact
	contacts = repository.GetAllContacts()

	for _, c := range contacts {
		json.NewEncoder(ctx).Encode(c)
	}
}

func CreateNewContact(ctx *fasthttp.RequestCtx) {
	var contactStruct ContactCreationStruct

	body := bytes.NewReader(ctx.PostBody())
	decoder := json.NewDecoder(body)

	if err := decoder.Decode(&contactStruct); err != nil {
		sendErrorMessage(ctx, "Error decoding the input")
		return
	}

	//used to check whether a string contains number or not, because value of name cannot contain number
	var IsLetter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

	//the function above will assume that whitespace is not letter
	//so we have to first strip all whitespace to correctly check the variable
	tempname := strings.Join(strings.Fields(contactStruct.Name), "")

	fmt.Println("Creating contact ", contactStruct.Name)

	if IsLetter(tempname) {
		repository.CreateContact(contactStruct.Name, contactStruct.PhoneNum, contactStruct.Address)
	} else {
		fmt.Println("name:", contactStruct.Name)
		sendErrorMessage(ctx, "Name cannot contain number")
	}
}

func EditContact(ctx *fasthttp.RequestCtx) {
	id, _ := ctx.UserValue("id").(string)

	var contactStruct ContactCreationStruct
	contactStruct.Id = id

	body := bytes.NewReader(ctx.PostBody())
	decoder := json.NewDecoder(body)

	if err := decoder.Decode(&contactStruct); err != nil {
		sendErrorMessage(ctx, "Error decoding the input")
		return
	}

	var IsLetter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString
	tempname := strings.Join(strings.Fields(contactStruct.Name), "")

	fmt.Println("Updating contact ", contactStruct.Name)

	if IsLetter(tempname) {
		repository.UpdateContact(contactStruct.Id, contactStruct.Name, contactStruct.PhoneNum, contactStruct.Address)
	} else {
		sendErrorMessage(ctx, "Name cannot contain number")
	}
}

func DeleteContact(ctx *fasthttp.RequestCtx) {
	id, _ := ctx.UserValue("id").(string)
	error := repository.DeleteContact(id)

	if error != nil {
		sendErrorMessage(ctx, "Name cannot contain number")
	}
}

func sendErrorMessage(ctx *fasthttp.RequestCtx, message string) {
	ctx.SetBody([]byte(message))
	ctx.SetStatusCode(fasthttp.StatusBadRequest)
}

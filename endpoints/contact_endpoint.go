package endpoints

import (
	"strconv"

	"github.com/valyala/fasthttp"

	"encoding/json"
	"fmt"
	"regexp"

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
	var IsLetter = regexp.MustCompile(`^[a-zA-Z]+[\s[a-zA-Z]+]*$`).MatchString

	fmt.Println("Creating contact ", contactStruct.Name)

	if IsLetter(contactStruct.Name) {
		repository.CreateContact(contactStruct.Name, contactStruct.PhoneNum, contactStruct.Address)
	} else {
		fmt.Println("invalid name:", contactStruct.Name)
		sendErrorMessage(ctx, "Name cannot contain number")
	}
}

func EditContact(ctx *fasthttp.RequestCtx) {
	id, _ := ctx.UserValue("id").(string)

	var contactStruct ContactCreationStruct
	contactStruct.Id = id

	body := bytes.NewReader(ctx.PostBody())
	decoder := json.NewDecoder(body)

	//contactToUpdate := repository.GetContactByID(id)

	if err := decoder.Decode(&contactStruct); err != nil {
		sendErrorMessage(ctx, "Error decoding the input")
		return
	}

	var IsLetter = regexp.MustCompile(`^[a-zA-Z]+[\s[a-zA-Z]+]*$`).MatchString
	fmt.Println("Updating contact ", contactStruct.Id)

	if IsLetter(contactStruct.Name) {
		repository.UpdateContact(contactStruct.Id, contactStruct.Name, contactStruct.PhoneNum, contactStruct.Address)
	} else {
		sendErrorMessage(ctx, "Name cannot contain number")
	}
}

func DeleteContact(ctx *fasthttp.RequestCtx) {
	id, _ := ctx.UserValue("id").(string)
	error := repository.DeleteContact(id)

	if error != nil {
		sendErrorMessage(ctx, "There is no user with that ID")
	}
}

func sendErrorMessage(ctx *fasthttp.RequestCtx, message string) {
	ctx.SetBody([]byte(message))
	ctx.SetStatusCode(fasthttp.StatusBadRequest)
}

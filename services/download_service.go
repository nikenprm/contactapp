package services

import (
	"encoding/csv"
	"log"

	"bytes"

	"github.com/contactapp/repository"
)

func ExportToCSV(listOfContacts []repository.Contact) ([]byte, error) {
	bytes := &bytes.Buffer{}
	w := csv.NewWriter(bytes)

	title := []string{"Name", "Phone Number", "Address"}
	w.Write(title)
	for _, contact := range listOfContacts {
		var record []string
		record = append(record, contact.Name)
		record = append(record, contact.PhoneNum)
		record = append(record, contact.Address)
		w.Write(record)
	}
	w.Flush()
	err := w.Error()
	if err != nil {
		log.Fatal(err)
	}
	return bytes.Bytes(), err
}

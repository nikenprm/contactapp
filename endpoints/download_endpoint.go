package endpoints

import (
	"github.com/contactapp/repository"
	"github.com/contactapp/services"
	"github.com/valyala/fasthttp"
)

func DownloadContactProfile(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "text/csv")
	ctx.Response.Header.Set("Content-Disposition", "attachment;filename=contact.csv")
	var contacts []repository.Contact
	contacts = repository.GetAllContacts()

	b, err := services.ExportToCSV(contacts)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	ctx.Response.Header.Set("Content-Disposition", "attachment; filename=contact.csv")
	ctx.Response.Header.Set("Content-Type", "text/csv")
	ctx.Write(b)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

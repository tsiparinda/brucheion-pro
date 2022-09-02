package api

import (
	"brucheion/models"

	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/logging"
)

type UserHandler struct {
	Repository models.Repository
	identity.User
	logging.Logger
}

func (h PassageHandler) GetUser() actionresults.ActionResult {
	// should be: {"status":"success","message":"","data":{"name":"albatros"}}
	// get userid
	userid := h.User.GetID()
	user := h.User.GetDisplayName()
	h.Logger.Debugf("Userid:", userid)
	resp := models.JSONResponse{
		Status:  "success",
		Message: "",
		Data:    models.User{Name: user},
	}

	return actionresults.NewJsonAction(resp)

}

// type ProductReference struct {
// 	models.Product
// 	CategoryID int
// }

// func (h RestHandler) PostProduct(p ProductReference) actionresults.ActionResult {
// 	if p.ID == 0 {
// 		return actionresults.NewJsonAction(h.processData(p))
// 	} else {
// 		return &StatusCodeResult{http.StatusBadRequest}
// 	}
// }

// func (h RestHandler) PutProduct(p ProductReference) actionresults.ActionResult {
// 	if p.ID > 0 {
// 		return actionresults.NewJsonAction(h.processData(p))
// 	} else {
// 		return &StatusCodeResult{http.StatusBadRequest}
// 	}
// }

// func (h RestHandler) processData(p ProductReference) models.Product {
// 	product := p.Product
// 	product.Category = &models.Category{
// 		ID: p.CategoryID,
// 	}
// 	h.Repository.SaveProduct(&product)
// 	return h.Repository.GetProduct(product.ID)
// }

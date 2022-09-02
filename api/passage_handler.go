package api

import (
	"brucheion/gocite"
	"brucheion/models"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	//"github.com/ThomasK81/gocite"
	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/logging"
)

type PassageHandler struct {
	Repository models.Repository
	identity.User
	logging.Logger
}

// func (h PassageHandler) GetPassages(urn string) actionresults.ActionResult {
// 	return actionresults.NewJsonAction(urn)
// }
func (h PassageHandler) GetPasageTest(urn string) actionresults.ActionResult {
	urn = "urn:cts:sktlit:skt0001.nyaya002.M3D:5.1.1"
	// get userid
	userid := 1 //h.User.GetID()
	user := h.User.GetDisplayName()
	h.Logger.Debugf("Userid:", urn)

	// check urn
	if !gocite.IsCTSURN(urn) {
		return actionresults.NewErrorAction(errors.New("Bad urn request"))
	}

	textRefs := h.Repository.SelectUserBuckets(userid)
	h.Logger.Debugf("textRefs:", textRefs)
	bucketName := strings.Join(strings.Split(urn, ":")[0:4], ":") + ":"
	h.Logger.Debugf("bucketName:", bucketName)

	d, err := h.Repository.SelectUserBucketKeyValue(userid, bucketName, urn)
	if err != nil {
		//http.Error(w, "Internal server error", 500)
		return actionresults.NewErrorAction(fmt.Errorf("Internal server error1", err))
	}
	c, err := h.Repository.SelectUserBucketKeyValue(userid, bucketName, bucketName)
	if err != nil {
		//http.Error(w, "Internal server error", 500)
		return actionresults.NewErrorAction(fmt.Errorf("Internal server error2", err))
	}

	catalog := models.BoltCatalog{}
	passage := gocite.Passage{}
	json.Unmarshal([]byte(d.JSON), &passage)
	json.Unmarshal([]byte(c.JSON), &catalog)

	text := passage.Text.TXT
	h.Logger.Debugf("passage.Text.TXT:", text)
	passages := strings.Split(text, "\r\n")
	work, err := h.retriveUserBucketWork(userid, bucketName)
	if err != nil {
		//http.Error(w, "Internal server error", 500)
		return actionresults.NewErrorAction(fmt.Errorf("Internal server error3", err))
	}
	h.Logger.Debugf("work:", work)
	var imageRefs []string
	for _, tmp := range passage.ImageLinks {
		imageRefs = append(imageRefs, tmp.Object)
	}

	p := models.Passage{
		ID:                 passage.PassageID,
		Transcriber:        user,
		TranscriptionLines: passages,
		PreviousPassage:    passage.Prev.PassageID,
		NextPassage:        passage.Next.PassageID,
		FirstPassage:       work.First.PassageID,
		LastPassage:        work.Last.PassageID,
		ImageRefs:          imageRefs,
		TextRefs:           textRefs,
		Catalog:            catalog,
	}
	//var p models.Passage
	//return actionresults.NewJsonAction(p)

	return actionresults.NewTemplateAction("testapi_layout.html",
		models.Passage{
			Transcriber: p.Transcriber,
			TextRefs:    []string{"1", "2"},
			ID:          p.ID,
		})
}

func (h PassageHandler) GetPassage(urn string) actionresults.ActionResult {
	urn = "urn:cts:sktlit:skt0001.nyaya002.M3D:5.1.1"
	h.Logger.Debugf("urn:", urn)
	// get userid
	userid := 1 //h.User.GetID()
	user := h.User.GetDisplayName()
	h.Logger.Debugf("Userid:", userid)

	// check urn
	if !gocite.IsCTSURN(urn) {
		return actionresults.NewErrorAction(errors.New("Bad urn request"))
	}

	textRefs := h.Repository.SelectUserBuckets(userid)
	// urn:cts:sktlit:skt0001.nyaya002.M3D:3.1.1 -> urn:cts:sktlit:skt0001.nyaya002.M3D:
	bucketName := strings.Join(strings.Split(urn, ":")[0:4], ":") + ":"

	d, err := h.Repository.SelectUserBucketKeyValue(userid, bucketName, urn)
	if err != nil {
		//http.Error(w, "Internal server error", 500)
		return actionresults.NewErrorAction(fmt.Errorf("Internal server error1", err))
	}
	c, err := h.Repository.SelectUserBucketKeyValue(userid, bucketName, bucketName)
	if err != nil {
		//http.Error(w, "Internal server error", 500)
		return actionresults.NewErrorAction(fmt.Errorf("Internal server error2", err))
	}

	catalog := models.BoltCatalog{}
	passage := gocite.Passage{}
	json.Unmarshal([]byte(d.JSON), &passage)
	json.Unmarshal([]byte(c.JSON), &catalog)

	text := passage.Text.TXT
	h.Logger.Debugf("passage.Text.TXT:", text)
	passages := strings.Split(text, "\r\n")
	work, err := h.retriveUserBucketWork(userid, bucketName)
	if err != nil {
		//http.Error(w, "Internal server error", 500)
		return actionresults.NewErrorAction(fmt.Errorf("Internal server error3", err))
	}
	//h.Logger.Debugf("work:", work.First.PassageID, work.Last.PassageID)
	var imageRefs []string
	for _, tmp := range passage.ImageLinks {
		imageRefs = append(imageRefs, tmp.Object)
	}

	p := models.Passage{
		ID:                 passage.PassageID,
		Transcriber:        user,
		TranscriptionLines: passages,
		PreviousPassage:    passage.Prev.PassageID,
		NextPassage:        passage.Next.PassageID,
		FirstPassage:       work.First.PassageID,
		LastPassage:        work.Last.PassageID,
		ImageRefs:          imageRefs,
		TextRefs:           textRefs,
		Catalog:            catalog,
	}
	h.Logger.Debugf("passage.PassageID:", passage.PassageID)

	resp := models.JSONResponse{
		Status:  "success",
		Message: "",
		Data:    p,
	}

	return actionresults.NewJsonAction(resp)
}

// SelectUserBucketWork retrieves an entire work from the users database as an (ordered) gocite.Work object
func (h PassageHandler) retriveUserBucketWork(userid int, urn string) (result gocite.Work, err error) {
	//h.Logger.Debugf("retriveUserBucketWork: input params:", userid, urn)
	dict := h.Repository.SelectUserBucketDict(userid, urn)
	result.WorkID = urn
	h.Logger.Debugf("retriveUserBucketWork: select user bucketdict: (len, urn) ", len(dict), result.WorkID)

	for _, pair := range dict {
		var passage gocite.Passage
		err := json.Unmarshal([]byte(pair.Value), &passage) //unmarshal the buffer and save the gocite.Passage
		if err != nil {
			return result, fmt.Errorf("retriveUserBucketWork: Error unmarshalling Passage: %s", err)
		}

		if passage.PassageID != "" {
			result.Passages = append(result.Passages, passage)
		}
	}
	return gocite.SortPassages(result)
}

// func (h RestHandler) GetProducts() actionresults.ActionResult {
// 	return actionresults.NewJsonAction(h.Repository.GetProducts())
// }

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

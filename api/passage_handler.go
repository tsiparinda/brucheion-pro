package api

import (
	"brucheion/gocite"
	"brucheion/models"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
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

// func (h PassageHandler) GetPassage(urn string) actionresults.ActionResult {
// 	h.Logger.Debugf("urn:", urn)
// 	return actionresults.NewJsonAction(urn)
// }

func (h PassageHandler) GetPassage(urn string) actionresults.ActionResult {

	h.Logger.Debugf("api.PassageHandler.GetPassage: urn=", urn)
	// get userid
	userid := h.User.GetID()
	user := h.User.GetDisplayName()
	h.Logger.Debugf("api.PassageHandler.GetPassage: Userid= ", userid)

	// receive an all of buckets in database
	textRefs := h.Repository.SelectUserBuckets(userid)
	if len(textRefs) == 0 {
		h.Logger.Info("api.PassageHandler.GetPassage: No user's buckets!")
	}

	// set the bucket default as first accessible user's bucket
	sort.Strings(textRefs)
	if urn == "undefined" {
		//urn = "urn:cts:sktlit:skt0001.nyaya002.M3D:5.1.1"
		urn = textRefs[0]
	}

	// check urn
	if !gocite.IsCTSURN(urn) {
		return actionresults.NewErrorAction(errors.New("api.PassageHandler.GetPassage: Bad urn request"))
	}

	// cut the end of URN for receive header
	// urn:cts:sktlit:skt0001.nyaya002.M3D:3.1.1 -> urn:cts:sktlit:skt0001.nyaya002.M3D:
	bucketName := strings.Join(strings.Split(urn, ":")[0:4], ":") + ":"

	// receive all of passages from user's bucket in sorted state
	work, err := h.retriveUserBucketWork(userid, bucketName)
	if err != nil {
		//http.Error(w, "Internal server error", 500)
		return actionresults.NewErrorAction(fmt.Errorf("api.PassageHandler.GetPassage: Internal server error3 %v", err))
	}

	//h.Logger.Infof("api.PassageHandler.GetPassage:", strings.LastIndex(urn, ":")+1, len(urn), urn, work.First.PassageID)
	// correct urn if it's short (from undefined)
	if strings.LastIndex(urn, ":")+1 == len(urn) {
		urn = work.First.PassageID
	}

	// receive a passage: key (urn....:x.y.z) -> value
	d, err := h.Repository.SelectUserBucketKeyValue(userid, bucketName, urn)
	if err != nil {
		//http.Error(w, "Internal server error", 500)
		return actionresults.NewErrorAction(fmt.Errorf("api.PassageHandler.GetPassage: Internal server error1 %v", err))
	}
	// receive a header
	c, err := h.Repository.SelectUserBucketKeyValue(userid, bucketName, bucketName)
	if err != nil {
		//http.Error(w, "Internal server error", 500)
		return actionresults.NewErrorAction(fmt.Errorf("api.PassageHandler.GetPassage: Internal server error2 %v", err))
	}

	catalog := models.BoltCatalog{}
	passage := gocite.Passage{}
	json.Unmarshal([]byte(d.JSON), &passage)
	json.Unmarshal([]byte(c.JSON), &catalog)

	// split passage lines to passages array
	text := passage.Text.TXT
	//h.Logger.Debugf("api.PassageHandler.GetPassage: passage.Text.TXT:", text)
	passages := strings.Split(text, "\r\n")

	//h.Logger.Debugf("work:", work.First.PassageID, work.Last.PassageID)
	var imageRefs []string
	for _, tmp := range passage.ImageLinks {
		imageRefs = append(imageRefs, tmp.Object)
	}

	p := models.Passage{
		ID:                 passage.PassageID,      // for current urn
		Transcriber:        user,                   // user name
		TranscriptionLines: passages,               // for current urn
		PreviousPassage:    passage.Prev.PassageID, // for current urn
		NextPassage:        passage.Next.PassageID, // for current urn
		FirstPassage:       work.First.PassageID,   // for current urn first number x.y.z
		LastPassage:        work.Last.PassageID,    // for current urn last number x.y.z
		ImageRefs:          imageRefs,              // for current urn
		TextRefs:           textRefs,               // array all of users urns
		Catalog:            catalog,                // header of current urn
	}
	//h.Logger.Debugf("api.PassageHandler.GetPassage: passage.PassageID:", passage.PassageID)

	// generate responce in json format
	resp := models.JSONResponse{
		Status:  "success",
		Message: "",
		Data:    p,
	}

	return actionresults.NewJsonAction(resp)
}

// SelectUserBucketWork retrieves an entire work from the users database as an (ordered) gocite.Work object
func (h PassageHandler) retriveUserBucketWork(userid int, urn string) (result gocite.Work, err error) {
	//h.Logger.Debugf("api.PassageHandler.retriveUserBucketWork: input params:", userid, urn)
	dict := h.Repository.SelectUserBucketDict(userid, urn)
	result.WorkID = urn
	h.Logger.Debugf("api.PassageHandler.retriveUserBucketWork: select user bucketdict: (len(dist), urn) ", len(dict), result.WorkID)

	for _, pair := range dict {
		var passage gocite.Passage
		err := json.Unmarshal([]byte(pair.Value), &passage) //unmarshal the buffer and save the gocite.Passage
		if err != nil {
			return result, fmt.Errorf("api.PassageHandler.retriveUserBucketWork: Error unmarshalling Passage: %s", err)
		}

		if passage.PassageID != "" {
			result.Passages = append(result.Passages, passage)
		}
	}
	return gocite.SortPassages(result)
}

package api

import (
	"github.com/vedicsociety/platform/http/actionresults"
)

type StatusCodeResult struct {
	code int
}

func (action *StatusCodeResult) Execute(ctx *actionresults.ActionContext) error {
	ctx.ResponseWriter.WriteHeader(action.code)
	return nil
}

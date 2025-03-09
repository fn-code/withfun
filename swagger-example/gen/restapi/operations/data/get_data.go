// Code generated by go-swagger; DO NOT EDIT.

package data

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetDataHandlerFunc turns a function with the right signature into a get data handler
type GetDataHandlerFunc func(GetDataParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetDataHandlerFunc) Handle(params GetDataParams) middleware.Responder {
	return fn(params)
}

// GetDataHandler interface for that can handle valid get data params
type GetDataHandler interface {
	Handle(GetDataParams) middleware.Responder
}

// NewGetData creates a new http.Handler for the get data operation
func NewGetData(ctx *middleware.Context, handler GetDataHandler) *GetData {
	return &GetData{Context: ctx, Handler: handler}
}

/* GetData swagger:route GET /v1/api/data data getData

Get Data show all user data, that already inserted

*/
type GetData struct {
	Context *middleware.Context
	Handler GetDataHandler
}

func (o *GetData) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetDataParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

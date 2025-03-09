package data

import (
	"github.com/fn-code/swagger-example/model"
	"net/http"
)


// swagger:route GET /v1/api/data data GetData
// Get Data show all user data, that already inserted
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// responses:
//  200: dataModelResponseWrapper
func GetData(w http.ResponseWriter, r *http.Request) {
	_ = model.DataModel{}
	w.Write([]byte("hola"))
}

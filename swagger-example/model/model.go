package model

// swagger:response dataModelResponseWrapper
type dataModelResponseWrapper struct {
	// All response for this user
	// in: body
	Body DataModel
}

// swagger:model dataModel
type DataModel struct {
	// id for userid
	//
	// required: true
	// example: 1
	ID int `json:"id"`
	// name is user name
	//
	// required: true
	// example: ludinnento
	Name string `json:"name"`
}

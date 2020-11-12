package controllers

// H is a struct for JSON Response
type H struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// NewH is the constructor of H
func NewH(message string, data interface{}) *H {
	H := new(H)
	H.Message = message
	H.Data = data
	return H
}

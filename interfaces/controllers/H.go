package controllers
import (
	"strconv"
	"github.com/hiroyaonoe/todoapp-server/usecase"
)

type H struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewH(message string, data interface{}) *H {
	H := new(H)
	H.Message = message
	H.Data = data
	return H
}

func NewHForRes(res *usecase.ResultStatus) *H {
	H := new(H)
	H.Message = strconv.Itoa(res.StatusCode)
	H.Data = res.StatusMessage
	return H
}

package controllers

import (
	"fmt"
)

type ErrorForJSON struct {
	Code int `json:"code"`
	Err string `json:"error"`
}

// ErrorToJSON はエラーが発生したときにステータスコードとメッセージをJSONにしてレスポンスを返す
func ErrorToJSON(c Context, statusCode int, err error) {
	fmt.Printf(err.Error())
	c.JSON(statusCode, &ErrorForJSON{
		Code:statusCode,
		Err:err.Error(),
	})
	return
}
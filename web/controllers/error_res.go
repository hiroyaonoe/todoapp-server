package controllers

// import (
// 	"fmt"
// )

type errorRes struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func newErrorRes(code int, err string) (res *errorRes) {
	res = &errorRes{
		Code: code,
		Err:  err,
	}
	return
}

// ErrorToJSON はエラーが発生したときにステータスコードとメッセージをJSONにしてレスポンスを返す
func errorToJSON(c Context, statusCode int, err error) {
	// fmt.Printf("[Error] %s ", err.Error())
	c.JSON(statusCode, newErrorRes(statusCode, err.Error()))
	return
}

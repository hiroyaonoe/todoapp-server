package controllers

type ErrorForJSON struct {
	Code int `json:"code"`
	Err string `json:"error"`
}

// errorToJSON はエラーが発生したときにステータスコードとメッセージをJSONにしてレスポンスを返す
func ErrorToJSON(c Context, statusCode int, err error) {
	c.JSON(statusCode, &ErrorForJSON{
		Code:statusCode,
		Err:err.Error(),
	})
	return
}
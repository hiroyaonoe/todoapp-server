package controllers

type errorForJSON struct {
	code int `json:"code"`
	err string `json:"error"`
}

// errorToJSON はエラーが発生したときにステータスコードとメッセージをJSONにしてレスポンスを返す
func errorToJSON(c Context, statusCode int, err error) {
	c.JSON(statusCode, errorForJSON{
		code:statusCode,
		err:err.Error(),
	})
	return
}
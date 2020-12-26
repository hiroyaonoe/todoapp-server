package controllers

// getUserIDFromCookie はcookieからuseridを取得する
func getUserIDFromCookie(c Context) (id string, err error) {
	id, err = c.Cookie("id")
	return
}

// getTaskIDFromParam はURIのParamかtaskidを取得する
func getTaskIDFromParam(c Context) (tid string) {
	tid = c.Param("id")
	return
}

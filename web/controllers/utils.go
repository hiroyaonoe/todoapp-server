package controllers

// getUserIDFromCookie はcookieからuseridを取得する
func getUserIDFromCookie(c Context) (id string, err error) {
	id, err = c.Cookie("id")
	return
}

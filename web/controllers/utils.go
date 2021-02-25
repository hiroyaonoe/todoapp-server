package controllers

import (
	"errors"
	"net/http"
)

// getUserIDFromCookie はcookieからuseridを取得する
func getUserIDFromCookie(c Context) (id string, err error) {
	id, err = c.Cookie("id")
	return
}

// getTaskIDFromParam はURIのParamからtaskidを取得する
func getTaskIDFromParam(c Context) (tid string, err error) {
	tid = c.Param("id")
	if tid == "" {
		return "", errors.New("URI parameter 'id' is nil")
	}
	return
}

// unexpectedErrorHandling は予期せぬエラーが発生したときのエラーハンドリングを行う
func unexpectedErrorHandling(c Context, _ error) {
	// panic(err.Error())
	errorToJSON(c, http.StatusInternalServerError, ErrInternalServerError)
	return
}

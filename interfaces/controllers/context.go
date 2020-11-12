package controllers

// Context is a interface for gin.Context
type Context interface {
	Param(key string) string
	JSON(code int, obj interface{})
	Cookie(name string) (string, error)
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
	ShouldBindJSON(obj interface{}) error
}

/*
controllers is Interface Adapters.

*/
package controllers

type Context interface {
    Param(key string) string
    JSON(code int, obj interface{})
}
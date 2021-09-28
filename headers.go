// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package neveauth

import "github.com/gin-gonic/gin"

const (
	HeaderKeyResource = "Neve-Resource"
	HeaderKeyAc = "Neve-Resource"
)

type headerOpt struct{}

var Gin headerOpt

func (o *headerOpt) WithResource(resource string) func(*gin.Context) {
	return func(context *gin.Context) {
		context.Header(HeaderKeyResource, resource)
	}
}

func (o *headerOpt) WithAction(action string) func(*gin.Context) {
	return func(context *gin.Context) {
		context.Header(HeaderKeyResource, action)
	}
}

// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package neveauth

import (
	"github.com/gin-gonic/gin"
	"github.com/xfali/neve-auth/constants"
)

type headerOpt struct{}

var Gin headerOpt

func (o *headerOpt) WithResource(resource string) func(*gin.Context) {
	return func(context *gin.Context) {
		context.Header(constants.HeaderKeyResource, resource)
	}
}

func (o *headerOpt) WithAction(action string) func(*gin.Context) {
	return func(context *gin.Context) {
		context.Header(constants.HeaderKeyAction, action)
	}
}

func (o *headerOpt) RequirePermission(resource string, action string) func(*gin.Context) {
	return func(context *gin.Context) {
		context.Header(constants.HeaderKeyResource, resource)
		context.Header(constants.HeaderKeyAction, action)
	}
}


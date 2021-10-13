// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package neveauth

import (
	"github.com/gin-gonic/gin"
	"github.com/xfali/neve-auth/attribute"
)

type headerOpt struct{}

var Gin headerOpt

func (o *headerOpt) RequirePermission(resource string, action string) func(*gin.Context) {
	return func(context *gin.Context) {
		context.Request = attribute.WithRequestInfo(context.Request, &attribute.RequestInfo{
			Resource: resource,
			Action:   action,
		})
	}
}

// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package filter

import (
	"github.com/gin-gonic/gin"
	"github.com/xfali/neve-auth/attribute"
	"github.com/xfali/neve-core/appcontext"
)

const(
	PermissionHandlerName = "neve_permission"
)

type PermissionHandler interface {
	RequirePermission(resource string, action string) func(*gin.Context)
}

type defaultPermissionHandler struct {
	filter *GinAuthorizeFilter
}

func NewPermissionHandler(filter *GinAuthorizeFilter) *defaultPermissionHandler {
	return &defaultPermissionHandler{
		filter: filter,
	}
}

func (o *defaultPermissionHandler) RequirePermission(resource string, action string) func(*gin.Context) {
	return func(context *gin.Context) {
		context.Request = attribute.WithRequestInfo(context.Request, &attribute.RequestInfo{
			Resource: resource,
			Action:   action,
		})
		o.filter.Authorize(context)
	}
}

type PermissionHandlerHolder struct {
	ph PermissionHandler
}

func (h *PermissionHandlerHolder) RegisterFunction(registry appcontext.InjectFunctionRegistry) error {
	return registry.RegisterInjectFunction(func(ph PermissionHandler) {
		h.ph = ph
	}, PermissionHandlerName)
}

func (h *PermissionHandlerHolder) RequirePermission(resource string, action string) func(*gin.Context) {
	return h.ph.RequirePermission(resource, action)
}

/*
 * Copyright (C) 2019-2024, Xiongfa Li.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package filter

import (
	"github.com/gin-gonic/gin"
	"github.com/xfali/neve-auth/attribute"
	"github.com/xfali/neve-core/appcontext"
)

const (
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

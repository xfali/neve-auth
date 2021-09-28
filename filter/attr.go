// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package filter

import (
	"github.com/gin-gonic/gin"
	"github.com/xfali/neve-auth/authorizer"
)

type AttributeParser interface {
	ParseAttribute(ctx *gin.Context) (authorizer.Attribute, error)
}


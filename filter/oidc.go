// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package filter

import (
	"github.com/gin-gonic/gin"
	"github.com/xfali/neve-auth/authorizer"
)

type oidcAttributeParser struct{}

func (p *oidcAttributeParser) ParseAttribute(ctx *gin.Context) (authorizer.Attribute, bool, error) {

}

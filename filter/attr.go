// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package filter

import (
	"github.com/xfali/neve-auth/attribute"
	"net/http"
)

type AttributeParser interface {
	ParseAttribute(req *http.Request) (attribute.Attribute, bool, error)
}


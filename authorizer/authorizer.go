// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package authorizer

import (
	"context"
	"github.com/xfali/neve-auth/attribute"
)

type Authorizer interface {
	Authorize(ctx context.Context, attr attribute.Attribute) Result
}

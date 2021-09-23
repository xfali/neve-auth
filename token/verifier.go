// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package token

import "context"

type Verifier interface {
	Verify(ctx context.Context, token string) (interface{}, error)
}

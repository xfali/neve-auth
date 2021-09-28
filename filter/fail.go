// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package filter

import "net/http"

type FailHandler interface {
	OnFailed(err error, resp http.ResponseWriter, req *http.Request)
}

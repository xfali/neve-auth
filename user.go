// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package neveauth

type UserInfo struct {
	UserID   string
	Username string

	Groups []string
	Extra  interface{}
}

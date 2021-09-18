// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package auth

import "time"

type Token struct {
	ID      string    `json:"id"`
	Refresh string    `json:"refresh"`
	Expire  time.Time `json:"expire"`
}

// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package config

import (
	"fmt"
	"github.com/xfali/fig"
)

func RedirectUrl(conf fig.Properties) (string, error) {
	//port := conf.Get("neve.server.port", "")
	//if port == "" {
	//	return "", fmt.Errorf("server port is empty")
	//}
	url := conf.Get("neve.auth.router.callback", "")
	if url == "" {
		return "", fmt.Errorf("callback url is empty")
	}
	//return fmt.Sprintf("https://"), nil

	addr := conf.Get("neve.auth.dex.externalAddr", "")
	if addr == "" {
		return "", fmt.Errorf("external address is empty")
	}
	return fmt.Sprintf("%s/%s", addr, url), nil
}

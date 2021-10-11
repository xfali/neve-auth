// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package constants

import (
	"fmt"
	"github.com/xfali/fig"
)

func RedirectUrl(conf fig.Properties) (string, error) {
	//port := conf.Get("neve.server.port", "")
	//if port == "" {
	//	return "", fmt.Errorf("server port is empty")
	//}
	url := conf.Get(RouterCallbackKey, "")
	if url == "" {
		return "", fmt.Errorf("callback url is empty")
	}
	//return fmt.Sprintf("https://"), nil
	if url[0] == '/' {
		url = url[1:]
	}

	addr := conf.Get(ExternalAddrKey, "")
	if addr == "" {
		return "", fmt.Errorf("external address is empty")
	}
	if addr[len(addr)-1] == '/' {
		addr = addr[:len(addr)-1]
	}


	return fmt.Sprintf("%s/%s", addr, url), nil
}

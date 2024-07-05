/*
 * Copyright (C) 2019-2024, Xiongfa Li.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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

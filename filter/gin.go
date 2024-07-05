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

package filter

import (
	"github.com/gin-gonic/gin"
)

const (
	ctxIngoreKey = "_neve_auth_ignore_"
)

func setIgnoreAuth(ctx *gin.Context) {
	ctx.Set(ctxIngoreKey, true)
}

func isIgnoreAuth(ctx *gin.Context) bool {
	v, ok := ctx.Get(ctxIngoreKey)
	if ok {
		return v.(bool)
	}
	return false
}

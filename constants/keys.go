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

const (
	IssuerKey       = "neve.auth.authenticator.spec.issuer"
	IssuerCAKey     = "neve.auth.authenticator.spec.issuerCA"
	ClientIdKey     = "neve.auth.authenticator.spec.client.id"
	ClientSecretKey = "neve.auth.authenticator.spec.client.secret"
	ExternalAddrKey = "neve.auth.authenticator.spec.externalAddr"

	RouterRedirectKey = "neve.auth.authenticator.spec.router.redirect"
	RouterUserInfoKey = "neve.auth.authenticator.spec.router.userinfo"
	RouterCallbackKey = "neve.auth.authenticator.spec.router.callback"

	IncludesKey = "neve.auth.includes"
	ExcludesKey = "neve.auth.excludes"
)

const (
	HeaderKeyResource = "Neve-Resource"
	HeaderKeyAction   = "Neve-Action"
)

const (
	DigestRealmKey = "neve.auth.digest.realm"
)

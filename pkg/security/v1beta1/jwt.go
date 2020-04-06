// Copyright © 2020 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1beta1

// JSON Web Token (JWT) token format for authentication as defined by
// [RFC 7519](https://tools.ietf.org/html/rfc7519). See [OAuth 2.0](https://tools.ietf.org/html/rfc6749) and
// [OIDC 1.0](http://openid.net/connect) for how this is used in the whole
// authentication flow.
//
// Examples:
//
// Spec for a JWT that is issued by `https://example.com`, with the audience claims must be either
// `bookstore_android.apps.example.com` or `bookstore_web.apps.example.com`.
// The token should be presented at the `Authorization` header (default). The Json web key set (JWKS)
// will be discovered followwing OpenID Connect protocol.
//
// ```yaml
// issuer: https://example.com
// audiences:
// - bookstore_android.apps.example.com
//   bookstore_web.apps.example.com
// ```
//
// This example specifies token in non-default location (`x-goog-iap-jwt-assertion` header). It also
// defines the URI to fetch JWKS explicitly.
//
// ```yaml
// issuer: https://example.com
// jwksUri: https://example.com/.secret/jwks.json
// jwtHeaders:
// - "x-goog-iap-jwt-assertion"
// ```
type JWTRule struct {
	// Identifies the issuer that issued the JWT. See
	// [issuer](https://tools.ietf.org/html/rfc7519#section-4.1.1)
	// A JWT with different `iss` claim will be rejected.
	//
	// Example: https://foobar.auth0.com
	// Example: 1234567-compute@developer.gserviceaccount.com
	Issuer string `json:"issuer,omitempty"`
	// The list of JWT
	// [audiences](https://tools.ietf.org/html/rfc7519#section-4.1.3).
	// that are allowed to access. A JWT containing any of these
	// audiences will be accepted.
	//
	// The service name will be accepted if audiences is empty.
	//
	// Example:
	//
	// ```yaml
	// audiences:
	// - bookstore_android.apps.example.com
	//   bookstore_web.apps.example.com
	// ```
	Audiences []string `json:"audiences,omitempty"`
	// URL of the provider's public key set to validate signature of the
	// JWT. See [OpenID Discovery](https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderMetadata).
	//
	// Optional if the key set document can either (a) be retrieved from
	// [OpenID
	// Discovery](https://openid.net/specs/openid-connect-discovery-1_0.html) of
	// the issuer or (b) inferred from the email domain of the issuer (e.g. a
	// Google service account).
	//
	// Example: `https://www.googleapis.com/oauth2/v1/certs`
	//
	// Note: Only one of jwks_uri and jwks should be used. jwks_uri will be ignored if it does.
	JwksURI string `json:"jwksUri,omitempty"`
	// JSON Web Key Set of public keys to validate signature of the JWT.
	// See https://auth0.com/docs/jwks.
	//
	// Note: Only one of jwks_uri and jwks should be used. jwks_uri will be ignored if it does.
	Jwks string `json:"jwks,omitempty"`
	// List of header locations from which JWT is expected. For example, below is the location spec
	// if JWT is expected to be found in `x-jwt-assertion` header, and have "Bearer " prefix:
	// ```
	//   fromHeaders:
	//   - name: x-jwt-assertion
	//     prefix: "Bearer "
	// ```
	FromHeaders []*JWTHeader `json:"fromHeaders,omitempty"`
	// List of query parameters from which JWT is expected. For example, if JWT is provided via query
	// parameter `my_token` (e.g /path?my_token=<JWT>), the config is:
	// ```
	//   fromParams:
	//   - "my_token"
	// ```
	FromParams []string `json:"fromParams,omitempty"`
	// This field specifies the header name to output a successfully verified JWT payload to the
	// backend. The forwarded data is `base64_encoded(jwt_payload_in_JSON)`. If it is not specified,
	// the payload will not be emitted.
	OutputPayloadToHeader string `json:"outputPayloadToHeader,omitempty"`
	// If set to true, the original token will be kept for the ustream request. Default is false.
	ForwardOriginalToken bool `json:"forwardOriginalToken,omitempty"`
}

// This message specifies a header location to extract JWT token.
type JWTHeader struct {
	// The HTTP header name.
	Name string `json:"name,omitempty"`
	// The prefix that should be stripped before decoding the token.
	// For example, for "Authorization: Bearer <token>", prefix="Bearer " with a space at the end.
	// If the header doesn't have this exact prefix, it is considerred invalid.
	Prefix string `json:"prefix,omitempty"`
}

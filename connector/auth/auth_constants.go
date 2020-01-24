// Copyright (c) 2020 InfraCloud Technologies
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package auth

var (
	// ToChannelFromBotLoginURL : Login URL
	//
	//DEPRECATED: DO NOT USE
	ToChannelFromBotLoginURL = []string{
		"https://login.microsoftonline.com/botframework.com/oauth2/v2.0/token",
	}

	// ToBotFromChannelOpenIDMetadataURL : OpenID metadata document for tokens coming from MSA
	ToBotFromChannelOpenIDMetadataURL = []string{
		"https://login.botframework.com/v1/.well-known/openidconfiguration",
	}

	// ToBotFromEnterpriseChannelOpenIDMetadataURLFormat : OpenID metadata document for tokens coming from MSA
	ToBotFromEnterpriseChannelOpenIDMetadataURLFormat = []string{
		"https://{channelService}.enterprisechannel.botframework.com",
		"/v1/.well-known/openidconfiguration",
	}

	// ToBotFromEmulatorOpenIDMetadataURL : OpenID metadata document for tokens coming from MSA
	ToBotFromEmulatorOpenIDMetadataURL = []string{
		"https://login.microsoftonline.com/common/v2.0/.well-known/openid-configuration",
	}

	// AllowedSigningAlgorithms : Tokens come from channels to the bot. The code
	//that uses this also supports tokens coming from the emulator.
	AllowedSigningAlgorithms = []string{"RS256", "RS384", "RS512"}
)

const (
	// ToChannelFromBotLoginURLPrefix : Login URL prefix
	ToChannelFromBotLoginURLPrefix = "https://login.microsoftonline.com/"

	// ToChannelFromBotTokenEndpointPathTOCHANNELFROMBOTTOKENENDPOINTPATH : Login URL token endpoint path
	ToChannelFromBotTokenEndpointPathTOCHANNELFROMBOTTOKENENDPOINTPATH = "/oauth2/v2.0/token"

	// DefaultChannelAuthTenant : Default tenant from which to obtain a token for bot to channel communication
	DefaultChannelAuthTenant = "botframework.com"

	// ToChannelFromBotOauthScope : OAuth scope to request
	ToChannelFromBotOauthScope = "https://api.botframework.com/.default"

	// ToBotFromChannelTokenIssuer : Token issuer
	ToBotFromChannelTokenIssuer = "https://api.botframework.com"

	// BotOpenIDMetadataKey : Application Setting Key for the OpenIdMetadataURL value.
	BotOpenIDMetadataKey = "BotOpenIdMetadata"

	// ChannelService : Application Setting Key for the ChannelService value.
	ChannelService = "ChannelService"

	// OauthURLKey Application Setting Key for the OAuthURL value.
	OauthURLKey = "OAuthApiEndpoint"

	// EmulateOauthCardsKey : Application Settings Key for whether to emulate OAuthCards when using the emulator.
	EmulateOauthCardsKey = "EmulateOAuthCards"

	// AuthorizedParty "azp" Claim.
	//Authorized party - the party to which the ID Token was issued.
	//This claim follows the general format set forth in the OpenID Spec.
	//    http://openid.net/specs/openid-connect-core-10.html#IDToken
	AuthorizedParty = "azp"

	/*AudienceClaim From RFC 7519.
	      https://tools.ietf.org/html/rfc7519#section-4.1.3
	  The "aud" (audience) claim identifies the recipients that the JWT is
	  intended for.  Each principal intended to process the JWT MUST
	  identify itself with a value in the audience claim.If the principal
	  processing the claim does not identify itself with a value in the
	  "aud" claim when this claim is present, then the JWT MUST be
	  rejected.In the general case, the "aud" value is an array of case-
	  sensitive strings, each containing a StringOrURI value.In the
	  special case when the JWT has one audience, the "aud" value MAY be a
	  single case-sensitive string containing a StringOrURI value.The
	  interpretation of audience values is generally application specific.
	  Use of this claim is OPTIONAL.
	*/
	AudienceClaim = "aud"

	/*IssuerClaim  From RFC 7519.
	      https://tools.ietf.org/html/rfc7519#section-4.1.1
	  The "iss" (issuer) claim identifies the principal that issued the
	  JWT.  The processing of this claim is generally application specific.
	  The "iss" value is a case-sensitive string containing a StringOrURI
	  value.  Use of this claim is OPTIONAL.
	*/
	IssuerClaim = "iss"

	/*KeyIDHeader From RFC 7515
	      https://tools.ietf.org/html/rfc7515#section-4.1.4
	  The "kid" (key ID) Header Parameter is a hint indicating which key
	  was used to secure the JWS. This parameter allows originators to
	  explicitly signal a change of key to recipients. The structure of
	  the "kid" value is unspecified. Its value MUST be a case-sensitive
	  string. Use of this Header Parameter is OPTIONAL.
	  When used with a JWK, the "kid" value is used to match a JWK "kid"
	  parameter value.
	*/
	KeyIDHeader = "kid"

	// VersionClaim Token version claim name. As used in Microsoft AAD tokens.
	VersionClaim = "ver"

	// AppIDClaim App ID claim name. As used in Microsoft AAD 1.0 tokens.
	AppIDClaim = "appid"

	// ServiceURLClaim Service URL claim name. As used in Microsoft Bot Framework v3.1 auth.
	ServiceURLClaim = "serviceurl"
)

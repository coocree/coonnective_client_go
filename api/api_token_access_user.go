package api

import (
	"github.com/lestrrat-go/jwx/jwt"
	"time"
)

type TokenAccessUserModel struct {
	apiUserAuthorize TokenAuthorizeUserModel
	clientId         string
	tokenType        string
	expiresIn        string
	accessToken      string
	refreshToken     string
	isValid          bool
	user             *UserModel
}

func TokenAccessUser(apiUserAuthorize TokenAuthorizeUserModel) *TokenAccessUserModel {
	return &TokenAccessUserModel{
		apiUserAuthorize: apiUserAuthorize,
		user:             User(nil),
	}
}

func (a *TokenAccessUserModel) Token() bool {
	apiConnection, _ := Connection(a.apiUserAuthorize.nonceToken, a.apiUserAuthorize.apiUri)

	const params = `
        mutation AclToken($input: TokenInput!) {
            AclToken(input: $input) {
                result {
                    accessToken
                    expiresIn
                    refreshToken
                    tokenType
                }
                success
                error {
                    code
                    createdAt
                    messages
                    module
                    path
                    variables
                }
                elapsedTime
            }
        }
    `
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"clientId":     a.user.ClientId(),
			"clientSecret": a.apiUserAuthorize.clientSecret,
			"code":         a.apiUserAuthorize.code,
			"codeVerifier": a.apiUserAuthorize.codeVerifier64,
			"grantType":    "authorization_code",
		},
	}

	apiResponse := apiConnection.Mutation(params, variables)
	if !apiResponse.IsValid() {
		apiResponse.ThrowException()
	}
	token := apiResponse.EndpointAuth("AclToken")

	if !token.IsValid() {
		token.ThrowException()
	}

	tokenResult := token.Result.(map[string]interface{})
	a.tokenType = tokenResult["tokenType"].(string)
	a.expiresIn = tokenResult["expiresIn"].(string)
	a.accessToken = tokenResult["accessToken"].(string)
	a.refreshToken = tokenResult["refreshToken"].(string)
	a.isValid = true
	a.user.AccessToken(a.accessToken)
	a.user.RefreshToken(a.refreshToken)

	jwToken, errJwt := jwt.Parse([]byte(a.accessToken))
	if errJwt != nil {
		apiError := ApiError([]string{"Token de acesso de usuário 'user' inválido.", errJwt.Error()},
			"TokenAccessUser", "Token", "ACCESS_TOKEN_FAIL", time.Now(), token)
		E(apiError.ToString(), apiError.Code, nil)
		panic(apiError.Code)
	}

	iss := jwToken.Issuer()
	if iss != "" {
		a.user.ServerUri(iss)
	} else {
		return false
	}

	return a.isValid
}

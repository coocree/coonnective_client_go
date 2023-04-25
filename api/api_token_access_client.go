package api

import (
	"fmt"
	"github.com/lestrrat-go/jwx/jwt"
	"time"
)

type TokenAccessClientModel struct {
	apiClientAuthorize TokenAuthorizeClientModel
	clientId           string
	tokenType          string
	expiresIn          string
	accessToken        string
	refreshToken       string
	apiUri             string
	user               *UserModel
}

func TokenAccessClient(apiClientAuthorize TokenAuthorizeClientModel) TokenAccessClientModel {
	return TokenAccessClientModel{
		apiClientAuthorize: apiClientAuthorize,
		clientId:           apiClientAuthorize.clientId,
		user:               User(nil),
	}
}

func (a *TokenAccessClientModel) Token() string {
	apiConnection, _ := Connection(a.apiClientAuthorize.nonceToken, a.apiClientAuthorize.apiUri)

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
			"clientId":     a.apiClientAuthorize.clientId,
			"clientSecret": a.apiClientAuthorize.clientSecret,
			"code":         a.apiClientAuthorize.code,
			"codeVerifier": a.apiClientAuthorize.codeVerifier64,
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

	a.user.AccessToken(a.accessToken)
	a.user.RefreshToken(a.refreshToken)

	return a.accessToken
}

func (a *TokenAccessClientModel) IsValid() (bool, error) {
	token := a.user.GetAccessToken()

	if token == "" {
		apiError := ApiError([]string{"Nenhum token de 'autorização' foi definido na inicialização da aplicação."},
			"TokenAuthorizeUser", "TokenAuthorizeUser", "AUTHORIZE_TOKEN_NOT_FOUND", time.Now(), nil)
		E(apiError.ToString(), apiError.Code, nil)
		panic(apiError.Code)
	}

	jwToken, errJwt := jwt.Parse([]byte(token))
	if errJwt != nil {
		apiError := ApiError([]string{"Token de autorização de usuário 'user' inválido.", errJwt.Error()},
			"TokenAuthorizeUser", "TokenAuthorizeUser", "AUTHORIZE_TOKEN_FAIL", time.Now(), token)
		E(apiError.ToString(), apiError.Code, nil)
		panic(apiError.Code)
	}

	isValid := false

	iss := jwToken.Issuer()
	sub := jwToken.Subject()
	fmt.Println("iss: ", iss)
	if iss != "" && sub == "access_token" {
		a.apiUri = iss
		isValid = true
		a.user.ServerUri(iss)
	}

	return isValid, nil
}

package api

import (
	"fmt"
	"github.com/lestrrat-go/jwx/jwt"
	"time"
)

type TokenAuthorizeUserModel struct {
	apiUri         string
	token          string
	clientSecret   string
	code           string
	codeChallenge  string
	codeVerifier   string
	codeVerifier64 string
	nonceToken     string
	state          string
	apiSecurity    *SecurityModel
	jwToken        jwt.Token
}

func TokenAuthorizeUser(token string) TokenAuthorizeUserModel {
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

	apiSecurity := Security()
	codeVerifier := apiSecurity.RandomBytes(32)
	codeVerifier64 := apiSecurity.Base64URLEncode(codeVerifier)
	codeChallenge := apiSecurity.Base64URLEncode(apiSecurity.EncodeSha256(codeVerifier64))
	state := apiSecurity.RandomBytes(16)

	return TokenAuthorizeUserModel{
		apiUri:         "",
		token:          token,
		clientSecret:   "",
		code:           "",
		codeChallenge:  codeChallenge,
		codeVerifier:   codeVerifier,
		codeVerifier64: codeVerifier64,
		nonceToken:     "",
		state:          state,
		jwToken:        jwToken,
	}
}

func (t *TokenAuthorizeUserModel) IsValid() bool {
	isValid := false
	if t.token != "" {
		iss := t.jwToken.Issuer()
		sub := t.jwToken.Subject()
		if iss != "" && sub == "access_token" {
			t.apiUri = iss
			isValid = true
		}
	}
	return isValid
}

func (t *TokenAuthorizeUserModel) Authorize() bool {
	if t.token == "" {
		apiError := ApiError([]string{"Nenhum token de 'autorização' foi definido na inicialização da aplicação."},
			"TokenAuthorizeUser", "Authorize", "AUTHORIZE_TOKEN_NOT_FOUND", time.Now(), t.token)

		E(apiError.ToString(), apiError.Code, nil)
		panic(apiError.Code)
	}

	apiUser := User(nil)
	clientId := apiUser.ClientId()

	iss := t.jwToken.Issuer()
	if iss != "" {
		t.apiUri = iss
		t.clientSecret = t.apiSecurity.EncodeSha256(clientId + t.codeVerifier64)
	}

	fmt.Println("clientId: ", clientId, iss)

	apiConnection, _ := Connection(t.token, t.apiUri)

	params := `
mutation AclAuthorize($input: AuthorizeInput!) {
	AclAuthorize(input: $input) {
		result {
			code
			nonceToken
			state
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

	apiStorage := Storage()
	if clientId != "" && apiStorage != nil {
		apiStorage.Add("clientId", clientId)
	}

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"clientId":      clientId,
			"codeChallenge": t.codeChallenge,
			"responseType":  "code",
			"scope":         "mut:authorize",
			"state":         t.state,
		},
	}

	apiResponse := apiConnection.Mutation(params, variables)

	var response AclAuthorizeResponse
	apiResponse.Endpoint(&response)

	authorizeResult := response.Data.AclAuthorize.Result

	if response.Data.AclAuthorize.Success {
		if authorizeResult.State == t.state {
			t.code = authorizeResult.Code
			t.nonceToken = authorizeResult.NonceToken
			return true
		}
	} else {
		//TODO: Tratar erro de autorização
		//response.Data.AclAuthorize.ThrowException()
	}

	return false
}
